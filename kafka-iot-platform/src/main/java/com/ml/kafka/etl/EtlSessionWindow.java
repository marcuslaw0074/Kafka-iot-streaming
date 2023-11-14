package com.ml.kafka.etl;

import java.io.IOException;
import java.time.Duration;
import java.util.Arrays;
import java.util.Map;
import java.util.Properties;
import java.util.concurrent.CountDownLatch;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.errors.SerializationException;
import org.apache.kafka.common.serialization.Deserializer;
import org.apache.kafka.common.serialization.Serde;
import org.apache.kafka.common.serialization.Serializer;
import org.apache.kafka.common.utils.Bytes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.KeyValue;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.Aggregator;
import org.apache.kafka.streams.kstream.Consumed;
import org.apache.kafka.streams.kstream.Grouped;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.Materialized;
import org.apache.kafka.streams.kstream.Merger;
import org.apache.kafka.streams.kstream.Named;
import org.apache.kafka.streams.kstream.SessionWindows;
import org.apache.kafka.streams.kstream.SlidingWindows;
import org.apache.kafka.streams.kstream.Suppressed;
import org.apache.kafka.streams.kstream.TimeWindows;
import org.apache.kafka.streams.state.SessionStore;
import org.apache.kafka.streams.state.WindowStore;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonSubTypes;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.ml.kafka.etl.EtlSlidingWindow.BMSAggregationData;
import com.ml.kafka.model.bms.BMSDataType;
import com.ml.kafka.model.bms.BMSEtlData;
import com.ml.kafka.model.bms.BMSRawData;
import com.ml.kafka.model.bms.json.JSONSerde;

import static org.apache.kafka.streams.kstream.Suppressed.BufferConfig.unbounded;

/*
 * 

docker exec -it kafka-broker bash   


mvn exec:java -Dexec.mainClass=com.ml.kafka.etl.EtlSlidingWindow

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-sink-test

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-stream-test

kafka-console-producer \
  --topic kafka-stream-test \
  --bootstrap-server localhost:9092 \
  --property parse.key=true \
  --property key.separator=":::"

{"_t":"bms.type","id":"user2","itemType":"elec"}:::{"_t":"bms.raw","id":"user1","value":31.431,"timestamp":15934567,"status":"1"}

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
public class EtlSessionWindow {

    public static void main(final String[] args) {

        final String application_id = "kafka-stream-etl-session-window";
        final String stream_name = "kafka-stream-test";
        final String sink_name = "kafka-sink-test";
        final String server = "localhost:9092";

        final Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, application_id);
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, server);
        // props.put(StreamsConfig.DEFAULT_TIMESTAMP_EXTRACTOR_CLASS_CONFIG,
        // AggregationTimestampExtractor.class);
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.COMMIT_INTERVAL_MS_CONFIG, 1000);

        // setting offset reset to earliest so that we can re-run the demo code with the
        // same pre-loaded data
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");

        final StreamsBuilder builder = new StreamsBuilder();

        final KStream<BMSDataType, BMSRawData> bmsData = builder.stream(stream_name,
                Consumed.with(new JSONSerde<>(), new JSONSerde<>()));

        Duration timeDifference = Duration.ofSeconds(3);
        Duration gracePeriod = Duration.ofSeconds(0);
        

        bmsData
                .map((key, value) -> {
                    // System.out.println(value);
                    return KeyValue.pair(key, value);
                })
                .filter((key, value) -> "elec".equals(key.itemType))
                .groupByKey(Grouped.with(new JSONSerde<>(), new JSONSerde<>()))
                .windowedBy(SessionWindows.ofInactivityGapAndGrace(timeDifference, gracePeriod))
                .aggregate(
                        () -> 0L,
                        new Aggregator<BMSDataType,BMSRawData,Long>() {
                            @Override
                            public Long apply(final BMSDataType key, final BMSRawData value, final Long aggregate) {
                                System.out.println(aggregate);
                                return aggregate;
                            }
                        },
                        new Merger<BMSDataType,Long>() {
                            @Override
                            public Long apply(final BMSDataType aggKey, final Long aggOne, final Long aggTwo) {
                                return aggOne;
                            }
                        }, 
                        Named.as("MyVeryCustomAggregator"),
                        Materialized.<BMSDataType, Long, SessionStore<Bytes, byte[]>>as(
                                "session-windowed-aggregated-stream-store-bms-data")
                        )
                // .count()
                .filter((key, value) -> {
                    System.out.println("filter");
                    return "elec".equals(key.key().itemType);
                })
                // .suppress(Suppressed.untilWindowCloses(unbounded()))
                .toStream()
                .map((key, value) -> {
                    System.out.println("End");
                    return KeyValue.pair(key.key(), value);
                });

        // .count()
        // .suppress(Suppressed.untilWindowCloses(unbounded()))
        // .toStream()
        // .foreach((key, value) -> {
        // System.out.println(value);
        // System.out.println("End");
        // });

        // bmsProcess.to(sink_name);

        final KafkaStreams streams = new KafkaStreams(builder.build(), props);
        // final CountDownLatch latch = new CountDownLatch(1);

        // attach shutdown handler to catch control-c
        streams.cleanUp();

        // Now run the processing topology via `start()` to begin processing its input
        // data.
        streams.start();
        System.out.println("S2f");

        // Add shutdown hook to respond to SIGTERM and gracefully close the Streams
        // application.
        Runtime.getRuntime().addShutdownHook(new Thread(streams::close));

        System.out.println("S2q");
    }

}
