package com.ml.kafka.etl.tenant;

import java.time.Duration;
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
import org.apache.kafka.streams.kstream.KTable;
import org.apache.kafka.streams.kstream.SessionWindows;

import com.ml.kafka.model.bms.BMSAggregationTenant;
import com.ml.kafka.model.bms.BMSDataType;
import com.ml.kafka.model.bms.BMSDeltaData;
import com.ml.kafka.model.bms.BMSTenantTagging;
import com.ml.kafka.model.bms.BMSTenantTaggingData;
import com.ml.kafka.model.bms.json.JSONDeserializer;
import com.ml.kafka.model.bms.json.JSONSerde;
import com.ml.kafka.model.bms.json.JSONSerializer;

/*

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-etl-tenant-raw

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-etl-tenant-realtime


kafka-console-producer \
  --topic kafka-etl-tenant-raw \
  --bootstrap-server localhost:9092 \
  --property parse.key=true \
  --property key.separator=":::"


{"_t":"bms.type","id":"user3","itemType":"elec"}:::{"_t":"bms.delta","id":"user1","value":31.431,"timestamp":1699760626130,"status":"1"}
{"_t":"bms.type","id":"user3","itemType":"elec"}:::{"_t":"bms.delta","id":"user2","value":14.431,"timestamp":1699760626130,"status":"1"}


 */

public class tenantEnergy {

    public static void main(final String[] args) {
        final String application_id = "kafka-etl-energy-tenant-app";
        final String stream_name = "kafka-etl-tenant-raw";
        final String table_name = "kafka-etl-tenant-tagging";
        // final String sink_name = "kafka-etl-tenant-realtime";
        final String server = "localhost:9092";

        final Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, application_id);
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, server);
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.COMMIT_INTERVAL_MS_CONFIG, 1000);
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");

        final StreamsBuilder builder = new StreamsBuilder();

        final KStream<BMSDataType, BMSDeltaData> stream = builder.stream(stream_name,
                Consumed.with(new JSONSerde<>(), new JSONSerde<>()));

        final KTable<BMSDataType, BMSTenantTagging> table = builder.stream(table_name, Consumed.with(
                Serdes.serdeFrom(new JSONSerializer<>(), new JSONDeserializer<>(BMSDataType.class)),
                Serdes.serdeFrom(new JSONSerializer<>(), new JSONDeserializer<>(BMSTenantTagging.class)))).toTable();

        stream
                .map((key, value) -> {
                    System.out.println(value);
                    return KeyValue.pair(key, value);
                })
                .leftJoin(table, (left, right) -> {
                    BMSTenantTaggingData data = new BMSTenantTaggingData();
                    data.bmsId = left.id;
                    data.itemType = left.timeType;
                    data.projectId = right.projectId;
                    data.value = left.value;
                    data.status = right.status;
                    data.tenantId = right.tenantId;
                    data.timestamp = left.timestamp;
                    return data;
                })
                .groupBy((key, value) -> {
                    BMSDataType newKey = new BMSDataType();
                    newKey.itemType = key.itemType;
                    newKey.id = "";
                    return newKey;
                })
                .windowedBy(
                        SessionWindows.ofInactivityGapAndGrace(
                                Duration.ofSeconds(5),
                                Duration.ofSeconds(30)))
                .aggregate(() -> new BMSAggregationTenant("", 0L, 1),
                        (aggKey, newValue, aggValue) -> {
                            aggValue.arr.add(newValue);
                            aggValue.map.put(newValue.bmsId, newValue.value);
                            System.out.println(aggValue);
                            return aggValue;
                        },
                        (aggKey, newValue, aggValue) -> {
                            for (BMSTenantTaggingData parts : newValue.arr) {
                                aggValue.arr.add(parts);
                            }
                            return aggValue;
                        });

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
