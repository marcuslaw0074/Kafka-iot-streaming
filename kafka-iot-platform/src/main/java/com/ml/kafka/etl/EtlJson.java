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
import org.apache.kafka.streams.kstream.Consumed;
import org.apache.kafka.streams.kstream.Grouped;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.Materialized;
import org.apache.kafka.streams.kstream.SlidingWindows;
import org.apache.kafka.streams.kstream.TimeWindows;
import org.apache.kafka.streams.state.WindowStore;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonSubTypes;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.ml.kafka.etl.EtlJson.BMSAggregationData;

/*
 * 

docker exec -it kafka-broker bash   


kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-streams-bms-data

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic kafka-streams-bms-realtime-data

kafka-console-producer \
  --topic kafka-streams-bms-data \
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
public class EtlJson {

    @SuppressWarnings("DefaultAnnotationParam") // being explicit for the example
    @JsonTypeInfo(use = JsonTypeInfo.Id.NAME, include = JsonTypeInfo.As.PROPERTY, property = "_t")
    @JsonSubTypes({
            @JsonSubTypes.Type(value = BMSRawData.class, name = "bms.raw"),
            @JsonSubTypes.Type(value = BMSDataType.class, name = "bms.type"),
            @JsonSubTypes.Type(value = BMSRealtimeData.class, name = "bms.realtime"),
            @JsonSubTypes.Type(value = BMSAggregationData.class, name = "bms.aggregation"),
            @JsonSubTypes.Type(value = BMSEtlData.class, name = "bms.etl")
    })
    public interface JSONSerdeCompatible {

    }

    static public class BMSData {
        @Override
        public String toString() {
            Gson gson = new GsonBuilder().setPrettyPrinting().create();
            return gson.toJson(this);
        }
    }

    static public class BMSDataType extends BMSData implements JSONSerdeCompatible {
        public String id;
        public String itemType;
    }

    static public class BMSRawData extends BMSData implements JSONSerdeCompatible {
        public String id;
        public double value;
        public String status;
        public Long timestamp;
    }

    static public class BMSRealtimeData extends BMSData implements JSONSerdeCompatible {
        public String id;
        public String block;
        public String buildingName;
        public String equipmentName;
        public String functionType;
        public String prefername;
        public String status;
        public Long value;
        public Long timestamp;
    }

    @JsonIgnoreProperties(ignoreUnknown = true)
    static public class BMSEtlData extends BMSData implements JSONSerdeCompatible {
        public double value;
        public String id;
        public Long timestamp;

        public BMSEtlData() {
        }

        public BMSEtlData(String id, double value, Long timestamp) {
            this.id = id;
            this.value = value;
            this.timestamp = timestamp;
        }
    }

    @JsonIgnoreProperties(ignoreUnknown = true)
    static public class BMSAggregationData extends BMSData implements JSONSerdeCompatible {
        public int count;
        public double value;
        public String type;
        public String id;
        public Long timestamp;

        // @JsonIgnore
        public double[] values = {};

        public BMSAggregationData() {
        }

        public BMSAggregationData(int count, double value) {
            this.count = count;
            this.value = value;
        }

        public void setTimestamp(Long timestamp) {
            this.timestamp = timestamp;
        }

        public void setId(String id) {
            this.id = id;
        }

        public void setType(String type) {
            this.type = type;
        }

        // @JsonIgnore
        public double[] getValues() {
            // System.out.println(Arrays.toString(this.values));
            return this.values;
        }

        // @JsonIgnore
        public void addValues(double value) {
            int i;
            int len = this.values.length;
            // System.out.println(len);
            double newarr[] = new double[len + 1];
            for (i = 0; i < len; i++) {
                newarr[i] = this.values[i];
            }
            newarr[len] = value;
            this.count++;
            this.values = newarr;
        }

        public BMSAggregationData agg() {
            double _value;
            switch (this.type) {
                case "mean":
                    _value = Arrays.stream(this.values).sum() / this.count;
                    break;
                case "sum":
                    _value = Arrays.stream(this.values).sum();
                    break;
                case "count":
                    _value = this.count;
                    break;
                case "max":
                    try {
                        _value = Arrays.stream(this.values).max().getAsDouble();
                    } catch (Exception e) {
                        _value = 0;
                    }
                    break;
                case "min":
                    try {
                        _value = Arrays.stream(this.values).min().getAsDouble();
                    } catch (Exception e) {
                        _value = 0;
                    }
                    break;
                case "first":
                    _value = this.values[0];
                    break;
                case "last":
                    _value = this.values[this.values.length - 1];
                    break;
                default:
                _value = 0;
                System.out.println("DEFAULT");
            }
            this.value = _value;
            return this;
        }
    }

    public static class JSONSerde<T extends JSONSerdeCompatible> implements Serializer<T>, Deserializer<T>, Serde<T> {
        private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

        @Override
        public void configure(final Map<String, ?> configs, final boolean isKey) {
        }

        @SuppressWarnings("unchecked")
        @Override
        public T deserialize(final String topic, final byte[] data) {
            if (data == null) {
                return null;
            }

            try {
                return (T) OBJECT_MAPPER.readValue(data, JSONSerdeCompatible.class);
            } catch (final IOException e) {
                throw new SerializationException(e);
            }
        }

        @Override
        public byte[] serialize(final String topic, final T data) {
            if (data == null) {
                return null;
            }

            try {
                return OBJECT_MAPPER.writeValueAsBytes(data);
            } catch (final Exception e) {
                throw new SerializationException("Error serializing JSON message", e);
            }
        }

        @Override
        public void close() {
        }

        @Override
        public Serializer<T> serializer() {
            return this;
        }

        @Override
        public Deserializer<T> deserializer() {
            return this;
        }

    }

    public static void main(final String[] args) {

        final String application_id = "kafka-stream-etl-json";
        final String stream_name = "kafka-streams-bms-data";
        final String sink_name = "kafka-streams-bms-realtime-data";
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

        Duration timeDifference = Duration.ofSeconds(1000);
        Duration gracePeriod = Duration.ofSeconds(10);

        final KStream<BMSDataType, BMSEtlData> bmsProcess = bmsData
                .filter((key, value) -> "elec".equals(key.itemType))
                .groupByKey(Grouped.with(new JSONSerde<>(), new JSONSerde<>()))
                // .windowedBy(SlidingWindows.ofTimeDifferenceAndGrace(timeDifference, gracePeriod))
                .windowedBy(SlidingWindows.ofTimeDifferenceWithNoGrace(timeDifference))
                // .windowedBy(TimeWindows.ofSizeWithNoGrace(Duration.ofSeconds(5)).advanceBy(Duration.ofSeconds(5)))
                .aggregate(
                        () -> new BMSAggregationData(0, 0),
                        (aggKey, newValue, aggValue) -> {
                            // System.out.println(aggKey.toString());
                            System.out.println("Start");
                            System.out.println(newValue.toString());
                            System.out.println(aggValue.toString());
                            aggValue.setType("sum");
                            aggValue.addValues(newValue.value);
                            aggValue.agg();
                            System.out.println(aggValue.toString());
                            System.out.println("End");
                            return aggValue;
                        },
                        Materialized.<BMSDataType, BMSAggregationData, WindowStore<Bytes, byte[]>>as(
                                "sliding-windowed-aggregated-stream-store-bms-data"))
                .toStream()
                .map((key, value) -> {
                    // System.out.println(key.key().toString());
                    // System.out.println(value.toString());
                    return KeyValue.pair(key.key(), new BMSEtlData(value.id, value.value, value.timestamp));
                });
        
        bmsProcess.to(sink_name);

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
            streams.start();
            latch.await();
        } catch (final Throwable e) {
            e.printStackTrace();
            System.exit(1);
        }
        System.exit(0);
    }
}
