package com.ml.kafka.etl;

import java.util.Properties;
import java.util.concurrent.CountDownLatch;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.common.utils.Bytes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.Materialized;
import org.apache.kafka.streams.kstream.Produced;
import org.apache.kafka.streams.state.KeyValueStore;
import org.apache.kafka.streams.KeyValue;

/*
 * 

docker exec -it kafka-broker bash   


kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-stream-count-input

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-sink-count-output



kafka-console-producer \
  --topic kafka-stream-count-input \
  --bootstrap-server localhost:9092 \
  --property parse.key=true \
  --property key.separator=":"

kafka-console-consumer --bootstrap-server localhost:9092 \
--topic kafka-sink-count-output \
--from-beginning \
--formatter kafka.tools.DefaultMessageFormatter \
--property print.timestamp=true \
--property print.key=true \
--property print.value=true

 */
public class EtlWindow {
    
    public static void main(String[] args) throws Exception {
        final String stream_name = "kafka-stream-count-input";
        final String sink_name = "kafka-sink-count-output";
        final String application_id = "kafka-stream-etl-client";
        final String server = "localhost:9092";
        final String store = "counts-store";
        final String shutdown = "kafka-stream-shutdown-hook";

        Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, application_id);
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, server);
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, Serdes.String().getClass());
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, Serdes.String().getClass());

        final StreamsBuilder builder = new StreamsBuilder();

        KStream<String, String> source = builder.stream(stream_name);
        source.groupByKey()
            .count(Materialized.<String, Long, KeyValueStore<Bytes, byte[]>>as(store))
            .toStream()
            .map((key, value) -> {
                System.out.println(key);
                System.out.println(value);
                return KeyValue.pair(key, Long.toString(value));
            })
            .to(sink_name, Produced.with(Serdes.String(), Serdes.String()));

        final KafkaStreams streams = new KafkaStreams(builder.build(), props);
        final CountDownLatch latch = new CountDownLatch(1);

        Runtime.getRuntime().addShutdownHook(new Thread(shutdown) {
            @Override
            public void run() {
                streams.close();
                latch.countDown();
            }
        });

        try {
            streams.start();
            latch.await();
        } catch (final Throwable e) {
            e.printStackTrace();
            System.exit(1);
        }
        System.exit(0);
    }
}
