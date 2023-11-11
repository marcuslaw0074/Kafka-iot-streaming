package com.ml.kafka.etl;

import java.time.Duration;
import java.util.Locale;
import java.util.Properties;
import java.util.concurrent.CountDownLatch;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.serialization.Serde;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.common.utils.Bytes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.KeyValue;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.Consumed;
import org.apache.kafka.streams.kstream.Grouped;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.Materialized;
import org.apache.kafka.streams.kstream.Named;
import org.apache.kafka.streams.kstream.SlidingWindows;
import org.apache.kafka.streams.processor.PunctuationType;
import org.apache.kafka.streams.processor.api.Processor;
import org.apache.kafka.streams.processor.api.ProcessorContext;
import org.apache.kafka.streams.processor.api.Record;
import org.apache.kafka.streams.state.KeyValueIterator;
import org.apache.kafka.streams.state.KeyValueStore;
import org.apache.kafka.streams.state.StoreBuilder;
import org.apache.kafka.streams.state.Stores;
import org.apache.kafka.streams.state.WindowStore;

import com.ml.kafka.model.TimeseriesData;
import com.ml.kafka.model.bms.*;
import com.ml.kafka.model.bms.json.JSONDeserializer;
import com.ml.kafka.model.bms.json.JSONSerde;
import com.ml.kafka.model.bms.json.JSONSerializer;


/*

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-streams-etl-raw-data-store


kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-streams-etl-realtime-data-store


kafka-console-producer \
  --topic kafka-streams-etl-raw-data-store \
  --bootstrap-server localhost:9092 \
  --property parse.key=true \
  --property key.separator=":::"


{"_t":"bms.type","id":"user2","itemType":"elec"}:::{"_t":"bms.etl","id":"user1","value":31.431,"timestamp":15934567,"status":"1"}

{"_t":"bms.type","id":"user2","itemType":"elec"}:::{"_t":"bms.etl","id":"user1","value":41.431,"timestamp":15934567,"status":"1"}


mvn exec:java -Dexec.mainClass=com.ml.kafka.etl.EtlEnergyStore


 */



@SuppressWarnings({ "WeakerAccess", "unused" })
public class EtlEnergyStore {

    static final String stateStoreName = "etl-energy-store";

    static public class EtlEnergyProcessor implements Processor<BMSDataType, BMSEtlData, BMSDataType, BMSEtlData> {
        private KeyValueStore<BMSDataType, BMSEtlData> kvStore;

        private ProcessorContext<BMSDataType, BMSEtlData> context;

        @Override
        public void init(final ProcessorContext<BMSDataType, BMSEtlData> context) {
            this.context = context;
            kvStore = this.context.getStateStore(stateStoreName);
        }

        @Override
        public void process(final Record<BMSDataType, BMSEtlData> record) {
            final BMSEtlData oldValue = kvStore.get(record.key());
            final BMSEtlData currentValue = record.value();
            if (currentValue.value > oldValue.value) {
                this.context.forward(new Record<>(record.key(),
                        new BMSEtlData(currentValue.id, currentValue.value - oldValue.value, currentValue.timestamp, 1),
                        0));
                kvStore.put(record.key(), currentValue);
            }

        }

        @Override
        public void close() {
            // close any resources managed by this processor
            // Note: Do not close any StateStores as these are managed by the library
        }
    }

    public static void main(final String[] args) {
        final String application_id = "kafka-stream-etl-energy-store";
        final String stream_name = "kafka-streams-etl-raw-data-store";
        final String sink_name = "kafka-streams-etl-realtime-data-store";
        final String server = "localhost:9092";

        final Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, application_id);
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, server);
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.COMMIT_INTERVAL_MS_CONFIG, 1000);
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");

        final StreamsBuilder builder = new StreamsBuilder();


        JSONDeserializer<BMSDataType> BMSDataTypeJsonDeserializer = new JSONDeserializer<>(BMSDataType.class);
        JSONSerializer<BMSDataType> BMSDataTypeJsonSerializer = new JSONSerializer<>();

        Serde<BMSDataType> BMSDataTypeSerde = Serdes.serdeFrom(BMSDataTypeJsonSerializer, BMSDataTypeJsonDeserializer);

        JSONDeserializer<BMSEtlData> BMSEtlDataJsonDeserializer = new JSONDeserializer<>(BMSEtlData.class);
        JSONSerializer<BMSEtlData> BMSEtlDataJsonSerializer = new JSONSerializer<>();

        Serde<BMSEtlData> BMSEtlDataSerde = Serdes.serdeFrom(BMSEtlDataJsonSerializer, BMSEtlDataJsonDeserializer);

        StoreBuilder<KeyValueStore<BMSDataType, BMSEtlData>> keyValueStoreBuilder = Stores.keyValueStoreBuilder(
                Stores.persistentKeyValueStore(stateStoreName),
                BMSDataTypeSerde,
                BMSEtlDataSerde);
        builder.addStateStore(keyValueStoreBuilder);

        final KStream<BMSDataType, BMSEtlData> stream = builder.stream(stream_name, Consumed.with(new JSONSerde<>(), new JSONSerde<>()));

        stream.map((key, value) -> {
            System.out.println(key.toString());
            System.out.println(value.toString());
                    return KeyValue.pair(key, value);
                }).process(() -> new EtlEnergyProcessor(), Named.as(sink_name)).to(sink_name);

        final KafkaStreams streams = new KafkaStreams(builder.build(), props);
        final CountDownLatch latch = new CountDownLatch(1);

        // attach shutdown handler to catch control-c
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
