package com.ml.kafka.etl;

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
import org.apache.kafka.streams.processor.api.Processor;
import org.apache.kafka.streams.processor.api.ProcessorContext;
import org.apache.kafka.streams.processor.api.Record;
import org.apache.kafka.streams.state.KeyValueStore;
import org.apache.kafka.streams.state.StoreBuilder;
import org.apache.kafka.streams.state.Stores;

import com.ml.kafka.model.bms.*;
import com.ml.kafka.model.bms.json.JSONDeserializer;
import com.ml.kafka.model.bms.json.JSONSerde;
import com.ml.kafka.model.bms.json.JSONSerializer;
import com.ml.kafka.serialization.timestamp.UnixTimestampConvertor;
import com.ml.kafka.stream.processor.period.EtlEnergyRealtimeProcessor;

/*

docker exec -it kafka-broker bash   

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-streams-etl-raw-period-store


kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-streams-etl-realtime-period-store


kafka-console-producer \
  --topic kafka-streams-etl-raw-period-store \
  --bootstrap-server localhost:9092 \
  --property parse.key=true \
  --property key.separator=":::"


{"_t":"bms.type","id":"user3","itemType":"elec"}:::{"_t":"bms.delta","id":"user1","value":31.431,"timestamp":1699760626130,"status":"1"}
{"_t":"bms.type","id":"user3","itemType":"elec"}:::{"_t":"bms.delta","id":"user1","value":41.431,"timestamp":1699760726130,"status":"1"}
{"_t":"bms.type","id":"user3","itemType":"elec"}:::{"_t":"bms.delta","id":"user1","value":141.431,"timestamp":1699760826130,"status":"1"}
{"_t":"bms.type","id":"user3","itemType":"elec"}:::{"_t":"bms.delta","id":"user1","value":176.4431,"timestamp":1699760926130,"status":"1"}
{"_t":"bms.type","id":"user3","itemType":"elec"}:::{"_t":"bms.delta","id":"user1","value":476.4431,"timestamp":1699860578000,"status":"1"}

{"_t":"bms.type","id":"user2","itemType":"elec"}:::{"_t":"bms.delta","id":"user1","value":41.431,"timestamp":1699761626130,"status":"1"}



{"_t":"bms.type","id":"user3","itemType":"elec"}:::{"_t":"bms.delta","id":"user1","value":1.4431,"timestamp":1699777778000,"status":"1"}
{"_t":"bms.type","id":"user3","itemType":"elec"}:::{"_t":"bms.delta","id":"user1","value":2.54245,"timestamp":1699781378000,"status":"1"}

mvn exec:java -Dexec.mainClass=com.ml.kafka.etl.EtlEnergyPeriodStore


 */

@SuppressWarnings({ "WeakerAccess", "unused" })
public class EtlEnergyPeriodStore {

    static final String realtimeStateStoreName = "etl-energy-realtime-store";
    static final String hourlyStateStoreName = "etl-energy-hourly-store";
    static final String dailyStateStoreName = "etl-energy-daily-store";
    static final String weeklyStateStoreName = "etl-energy-weekly-store";
    static final String monthlyStateStoreName = "etl-energy-monthly-store";

    static public class EtlEnergyProcessor implements Processor<BMSDataType, BMSDeltaData, BMSDataType, BMSDeltaData> {
        private KeyValueStore<BMSDataType, BMSDeltaData> realtimeKvStore;
        private KeyValueStore<BMSDataType, BMSDeltaData> hourlyKvStore;
        private KeyValueStore<BMSDataType, BMSDeltaData> dailyKvStore;

        private ProcessorContext<BMSDataType, BMSDeltaData> context;

        @Override
        public void init(final ProcessorContext<BMSDataType, BMSDeltaData> context) {
            this.context = context;
            realtimeKvStore = this.context.getStateStore(realtimeStateStoreName);
            hourlyKvStore = this.context.getStateStore(hourlyStateStoreName);
            dailyKvStore = this.context.getStateStore(dailyStateStoreName);
            System.out.println("Initialized");
        }

        private void realtimeProcess(final Record<BMSDataType, BMSDeltaData> record) {
            final BMSDeltaData oldValue = realtimeKvStore.get(record.key());
            final BMSDeltaData currentValue = record.value();
            System.out.println("realtime");
            System.out.println(oldValue);
            System.out.println(currentValue);

            if (oldValue == null) {
                realtimeKvStore.put(record.key(), record.value());
            } else {

                if (currentValue.value >= oldValue.value) {
                    this.context.forward(new Record<>(record.key(),
                            new BMSDeltaData(currentValue.id, currentValue.value - oldValue.value,
                                    currentValue.timestamp,
                                    currentValue.timestamp - oldValue.timestamp,
                                    1, "realtime"),
                            currentValue.timestamp));
                    realtimeKvStore.put(record.key(), currentValue);
                } else {
                    System.out.println("realtime-pass");
                }
            }
        
        }

        private void hourlyProcess(final Record<BMSDataType, BMSDeltaData> record) {
            final BMSDeltaData oldValue = hourlyKvStore.get(record.key());
            final BMSDeltaData currentValue = record.value();
            System.out.println("hourly");
            System.out.println(oldValue);
            System.out.println(currentValue);

            if (oldValue == null) {
                hourlyKvStore.put(record.key(), record.value());
                System.out.println("NULL");
            } else {
                UnixTimestampConvertor oldTime = new UnixTimestampConvertor(oldValue.timestamp, +8);
                UnixTimestampConvertor currentTime = new UnixTimestampConvertor(currentValue.timestamp, +8);

                if (currentTime.isNextHour(oldTime) && (currentValue.value >= oldValue.value)) {
                    this.context.forward(new Record<>(record.key(),
                            new BMSDeltaData(currentValue.id, currentValue.value - oldValue.value,
                                    currentValue.timestamp,
                                    currentValue.timestamp - oldValue.timestamp,
                                    1, "hourly"),
                            currentValue.timestamp));
                    hourlyKvStore.put(record.key(), currentValue);
                } else {
                    System.out.println("hourly-pass");
                }
            }
        }

        private void dailyProcess(final Record<BMSDataType, BMSDeltaData> record) {
            final BMSDeltaData oldValue = dailyKvStore.get(record.key());
            final BMSDeltaData currentValue = record.value();
            System.out.println("daily");
            System.out.println(oldValue);
            System.out.println(currentValue);

            if (oldValue == null) {
                dailyKvStore.put(record.key(), record.value());
            } else {
                UnixTimestampConvertor oldTime = new UnixTimestampConvertor(oldValue.timestamp, +8);
                UnixTimestampConvertor currentTime = new UnixTimestampConvertor(currentValue.timestamp, +8);

                if (currentTime.isNextDay(oldTime) && (currentValue.value >= oldValue.value)) {
                    this.context.forward(new Record<>(record.key(),
                            new BMSDeltaData(currentValue.id, currentValue.value - oldValue.value,
                                    currentValue.timestamp,
                                    currentValue.timestamp - oldValue.timestamp,
                                    1, "daily"),
                            currentValue.timestamp));
                    dailyKvStore.put(record.key(), currentValue);
                } else {
                    System.out.println("daily-pass");
                }
            }
        }

        @Override
        public void process(final Record<BMSDataType, BMSDeltaData> record) {
            System.out.println("Start");

            this.realtimeProcess(record);
            this.hourlyProcess(record);
            this.dailyProcess(record);
        }

        @Override
        public void close() {
            realtimeKvStore.close();
            hourlyKvStore.close();
            dailyKvStore.close();
        }
    }

    public static void main(final String[] args) {
        final String application_id = "kafka-stream-etl-energy-period-store";
        final String stream_name = "kafka-streams-etl-raw-period-store";
        final String sink_name = "kafka-streams-etl-realtime-period-store";
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

        StoreBuilder<KeyValueStore<BMSDataType, BMSDeltaData>> etlEnergyhourlyStoreBuilder = Stores
                .keyValueStoreBuilder(
                        Stores.persistentKeyValueStore(hourlyStateStoreName),
                        Serdes.serdeFrom(new JSONSerializer<>(), new JSONDeserializer<>(BMSDataType.class)),
                        Serdes.serdeFrom(new JSONSerializer<>(), new JSONDeserializer<>(BMSDeltaData.class)));

        StoreBuilder<KeyValueStore<BMSDataType, BMSDeltaData>> etlEnergydailyStoreBuilder = Stores
                .keyValueStoreBuilder(
                        Stores.persistentKeyValueStore(dailyStateStoreName),
                        Serdes.serdeFrom(new JSONSerializer<>(), new JSONDeserializer<>(BMSDataType.class)),
                        Serdes.serdeFrom(new JSONSerializer<>(), new JSONDeserializer<>(BMSDeltaData.class)));

        builder.addStateStore(etlEnergyRealtimeStoreBuilder);
        builder.addStateStore(etlEnergyhourlyStoreBuilder);
        builder.addStateStore(etlEnergydailyStoreBuilder);

        final KStream<BMSDataType, BMSDeltaData> stream = builder.stream(stream_name,
                Consumed.with(new JSONSerde<>(), new JSONSerde<>()));

        stream
                .map((key, value) -> {
                    System.out.println(value);
                    return KeyValue.pair(key, value);
                })
                .process(() -> new EtlEnergyProcessor(),
                        realtimeStateStoreName, hourlyStateStoreName, dailyStateStoreName)
                // .process(() -> new EtlEnergyRealtimeProcessor(realtimeStateStoreName),
                //         realtimeStateStoreName, hourlyStateStoreName, dailyStateStoreName)
                .map((key, value) -> {
                    System.out.println("_______");
                    System.out.println(value);
                    System.out.println("_______");
                    return KeyValue.pair(key, value);
                }).to(sink_name);

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
