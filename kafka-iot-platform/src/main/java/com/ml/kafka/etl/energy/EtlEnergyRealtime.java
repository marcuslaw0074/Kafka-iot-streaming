package com.ml.kafka.etl.energy;

import java.util.Properties;
import java.util.concurrent.CountDownLatch;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.KeyValue;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.Consumed;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.state.KeyValueStore;
import org.apache.kafka.streams.state.StoreBuilder;
import org.apache.kafka.streams.state.Stores;

import com.ml.kafka.model.bms.BMSDataType;
import com.ml.kafka.model.bms.BMSDeltaData;
import com.ml.kafka.model.bms.json.JSONDeserializer;
import com.ml.kafka.model.bms.json.JSONSerde;
import com.ml.kafka.model.bms.json.JSONSerializer;
import com.ml.kafka.stream.processor.period.EtlEnergyRealtimeProcessor;

/*

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-etl-energy-raw


kafka-console-producer \
  --topic kafka-etl-energy-raw \
  --bootstrap-server localhost:9092 \
  --property parse.key=true \
  --property key.separator=":::"

{"_t":"bms.type","id":"user3","itemType":"elec"}:::{"_t":"bms.delta","id":"user1","value":176.4431,"timestamp":1699786626130,"status":1}

 */

public class EtlEnergyRealtime {
    static final String realtimeStateStoreName = "etl-energy-realtime-store";

    public static void main(final String[] args) {
        final String application_id = "kafka-etl-energy-realtime-app";
        final String stream_name = "kafka-etl-energy-raw";
        final String sink_name = "kafka-etl-energy-realtime";
        final String server = "localhost:9092";

        final Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, application_id);
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, server);
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.COMMIT_INTERVAL_MS_CONFIG, 1000);
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");

        final StreamsBuilder builder = new StreamsBuilder();

        StoreBuilder<KeyValueStore<BMSDataType, BMSDeltaData>> etlEnergyRealtimeStoreBuilder = Stores
                .keyValueStoreBuilder(
                        Stores.persistentKeyValueStore(realtimeStateStoreName),
                        Serdes.serdeFrom(new JSONSerializer<>(), new JSONDeserializer<>(BMSDataType.class)),
                        Serdes.serdeFrom(new JSONSerializer<>(), new JSONDeserializer<>(BMSDeltaData.class)));

        builder.addStateStore(etlEnergyRealtimeStoreBuilder);

        final KStream<BMSDataType, BMSDeltaData> stream = builder.stream(stream_name,
                Consumed.with(new JSONSerde<>(), new JSONSerde<>()));

        stream
                .map((key, value) -> {
                    System.out.println(value);
                    return KeyValue.pair(key, value);
                })
                .process(() -> new EtlEnergyRealtimeProcessor(realtimeStateStoreName),
                        realtimeStateStoreName)
                .map((key, value) -> {
                    System.out.println("_______");
                    System.out.println(value);
                    System.out.println("_______");
                    return KeyValue.pair(key, value);
                }).to(sink_name);

        final KafkaStreams streams = new KafkaStreams(builder.build(), props);
        final CountDownLatch latch = new CountDownLatch(1);

        Runtime.getRuntime().addShutdownHook(new Thread("streams-bms-shutdown-hook") {
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
