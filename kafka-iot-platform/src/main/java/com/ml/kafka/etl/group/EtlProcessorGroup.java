package com.ml.kafka.etl.group;

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
import com.ml.kafka.stream.processor.EtlGroupProcessor;

/*
 * 

docker exec -it kafka-broker bash   


mvn exec:java -Dexec.mainClass=com.ml.kafka.etl.EtlSlidingWindow

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-stream-group

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-sink-group

kafka-console-producer \
  --topic kafka-stream-group \
  --bootstrap-server localhost:9092 \
  --property parse.key=true \
  --property key.separator=":::"

{"_t":"bms.type","id":"user1","itemType":"elec"}:::{"_t":"bms.delta","id":"user1","value":31.431,"timestamp":15934567,"status":"1"}
{"_t":"bms.type","id":"user2","itemType":"elec"}:::{"_t":"bms.delta","id":"user2","value":12.43134,"timestamp":15934567,"status":"1"}
{"_t":"bms.type","id":"user3","itemType":"elec"}:::{"_t":"bms.delta","id":"user3","value":321.431,"timestamp":15934567,"status":"1"}
{"_t":"bms.type","id":"user4","itemType":"elec"}:::{"_t":"bms.delta","id":"user4","value":-120.4931,"timestamp":15934567,"status":"1"}

kafka-console-consumer --bootstrap-server localhost:9092 \
--topic kafka-streams-bms-realtime-data \
--from-beginning \
--formatter kafka.tools.DefaultMessageFormatter \
--property print.timestamp=true \
--property print.key=true \
--property print.value=true


 * 
 * 
 */


@SuppressWarnings({ "WeakerAccess", "unused" })
public class EtlProcessorGroup {
    static final String stateStoreName = "etl-processor-group-store";
    static final String contextStateStoreName = "etl-processor-context-store";

    public static void main(final String[] args) {
        final String application_id = "kafka-etl-energy-hourly-app";
        final String stream_name = "kafka-stream-group";
        final String sink_name = "kafka-sink-group";
        final String server = "localhost:9092";

        final Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, application_id);
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, server);
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.COMMIT_INTERVAL_MS_CONFIG, 1000);
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");

        final StreamsBuilder builder = new StreamsBuilder();

        StoreBuilder<KeyValueStore<BMSDataType, BMSDeltaData>> etlEnergyhourlyStoreBuilder = Stores
                .keyValueStoreBuilder(
                        Stores.persistentKeyValueStore(stateStoreName),
                        Serdes.serdeFrom(new JSONSerializer<>(), new JSONDeserializer<>(BMSDataType.class)),
                        Serdes.serdeFrom(new JSONSerializer<>(), new JSONDeserializer<>(BMSDeltaData.class)));

        StoreBuilder<KeyValueStore<String, Long>> contextStoreBuilder = Stores
                .keyValueStoreBuilder(
                        Stores.persistentKeyValueStore(contextStateStoreName),
                        Serdes.String(),
                        Serdes.Long());


        builder.addStateStore(etlEnergyhourlyStoreBuilder);
        builder.addStateStore(contextStoreBuilder);

        final KStream<BMSDataType, BMSDeltaData> stream = builder.stream(stream_name,
                Consumed.with(new JSONSerde<>(), new JSONSerde<>()));

        stream
                .map((key, value) -> {
                    System.out.println(value);
                    return KeyValue.pair(key, value);
                })
                .process(() -> new EtlGroupProcessor(stateStoreName, contextStateStoreName),
                        stateStoreName, contextStateStoreName)
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
