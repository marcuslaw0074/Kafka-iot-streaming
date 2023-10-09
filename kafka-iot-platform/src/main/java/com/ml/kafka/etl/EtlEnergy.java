package com.ml.kafka.etl;

import java.time.Duration;
import java.util.Properties;
import java.util.concurrent.CountDownLatch;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.utils.Bytes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.KeyValue;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.Consumed;
import org.apache.kafka.streams.kstream.Grouped;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.Materialized;
import org.apache.kafka.streams.kstream.SlidingWindows;
import org.apache.kafka.streams.state.WindowStore;

import com.ml.kafka.model.bms.*;
import com.ml.kafka.model.bms.json.JSONSerde;

/*

docker exec -it kafka-broker bash   


mvn exec:java -Dexec.mainClass=com.ml.kafka.etl.EtlEnergy

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-streams-energy-raw-data

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-streams-energy-realtime-data

kafka-console-producer \
  --topic kafka-streams-energy-raw-data \
  --bootstrap-server localhost:9092 \
  --property parse.key=true \
  --property key.separator=":::"

kafka-console-consumer --bootstrap-server localhost:9092 \
--topic kafka-streams-energy-realtime-data \
--from-beginning \
--formatter kafka.tools.DefaultMessageFormatter \
--property print.timestamp=true \
--property print.key=true \
--property print.value=true




 */

@SuppressWarnings({ "WeakerAccess", "unused" })
public class EtlEnergy {
    public static void main(final String[] args) {
        final String application_id = "kafka-stream-etl-energy";
        final String stream_name = "kafka-streams-energy-raw-data";
        final String sink_name = "kafka-streams-energy-realtime-data";
        final String server = "localhost:9092";

        final Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, application_id);
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, server);
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.COMMIT_INTERVAL_MS_CONFIG, 1000);
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");

        final StreamsBuilder builder = new StreamsBuilder();

        final KStream<BMSDataType, BMSRawData> bmsData = builder.stream(stream_name,
                Consumed.with(new JSONSerde<>(), new JSONSerde<>()));

        Duration timeDifference = Duration.ofSeconds(4);
        Duration gracePeriod = Duration.ofSeconds(0);

        final KStream<BMSDataType, BMSEtlData> bmsProcess = bmsData
                .filter((key, value) -> "elec".equals(key.itemType))
                .groupByKey(Grouped.with(new JSONSerde<>(), new JSONSerde<>()))
                .windowedBy(SlidingWindows.ofTimeDifferenceAndGrace(timeDifference, gracePeriod))
                .aggregate(
                        () -> new BMSAggregationData(0, 0),
                        (aggKey, newValue, aggValue) -> {
                            System.out.println("Start");
                            System.out.println(newValue.toString());
                            System.out.println(aggValue.toString());
                            aggValue.setType("diff");
                            aggValue.setId(aggKey.id);
                            aggValue.setTimestamp(newValue.timestamp);
                            aggValue.addValues(newValue.value);
                            aggValue.agg();
                            System.out.println(aggValue.toString());
                            System.out.println("End");
                            return aggValue;
                        },
                        Materialized.<BMSDataType, BMSAggregationData, WindowStore<Bytes, byte[]>>as(
                                "sliding-windowed-aggregated-stream-store-energy-data"))
                .toStream()
                .filter((key, value) -> value.status == 1)
                .map((key, value) -> {
                    return KeyValue.pair(key.key(),
                            new BMSEtlData(value.id, value.value, key.window().end(), value.status));
                });

        bmsProcess.to(sink_name);

        final KafkaStreams streams = new KafkaStreams(builder.build(), props);
        final CountDownLatch latch = new CountDownLatch(1);

        Runtime.getRuntime().addShutdownHook(new Thread("streams-energy-shutdown-hook") {
            @Override
            public void run() {
                streams.close();
                latch.countDown();
            }
        });

        try {
            streams.cleanUp();
            streams.start();
            latch.await();
        } catch (final Throwable e) {
            e.printStackTrace();
            System.exit(1);
        }
        System.exit(0);
    }
}
