package com.ml.kafka.etl.group;

import java.util.Properties;
import java.util.concurrent.CountDownLatch;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.common.utils.Bytes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.KeyValue;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.Consumed;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.KTable;
import org.apache.kafka.streams.kstream.Materialized;
import org.apache.kafka.streams.state.KeyValueIterator;
import org.apache.kafka.streams.state.KeyValueStore;
import org.apache.kafka.streams.state.StoreBuilder;
import org.apache.kafka.streams.state.Stores;
import org.apache.kafka.streams.processor.api.Record;

import com.ml.kafka.model.bms.BMSDataType;
import com.ml.kafka.model.bms.BMSDeltaData;
import com.ml.kafka.model.bms.BMSMapData;
import com.ml.kafka.model.bms.BMSTenantTagging;
import com.ml.kafka.model.bms.json.JSONDeserializer;
import com.ml.kafka.model.bms.json.JSONSerde;
import com.ml.kafka.model.bms.json.JSONSerializer;
import com.ml.kafka.stream.processor.CustomGroupProcessor;

public class customProcessor {
    static final String stateStoreName = "etl-processor-group-store";
    static final String contextStateStoreName = "etl-processor-context-store";
    static final String tableStoreName = "kafka-table-store";

    public static void main(final String[] args) {
        final String application_id = "kafka-etl-energy-hourly-app";
        final String stream_name = "kafka-stream-group";
        final String sink_name = "kafka-sink-group";
        final String table_name = "kafka-table-group";
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

        final KTable<BMSDataType, BMSTenantTagging> tenantKTable = builder.table(table_name,
                Materialized.<BMSDataType, BMSTenantTagging, KeyValueStore<Bytes, byte[]>>as(tableStoreName));

        tenantKTable
                .filter((key, value) -> {
                    System.out.println(key);
                    System.out.println(value);
                    return true;
                });

        stream
                .map((key, value) -> {
                    System.out.println(value);
                    return KeyValue.pair(key, value);
                })
                .leftJoin(tenantKTable, (key, value) -> {
                    return key;
                    // return new BMSTenantTaggingData(key.id, value.tenantId, value.itemType,
                    // value.projectId, key.status, key.value, key.timestamp);
                })
                .process(() -> new CustomGroupProcessor<BMSDataType, BMSDeltaData, BMSDataType, BMSMapData>(
                        stateStoreName) {
                    @Override
                    public Record<BMSDataType, BMSMapData> generatRecord(
                            KeyValueIterator<BMSDataType, BMSDeltaData> iter) {
                        return null;
                    }
                },
                        stateStoreName, contextStateStoreName)
                .map((key, value) -> {
                    System.out.println("_______");
                    System.out.println(value);
                    System.out.println("_______");
                    return KeyValue.pair(key, value);
                })
                .to(sink_name);

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
