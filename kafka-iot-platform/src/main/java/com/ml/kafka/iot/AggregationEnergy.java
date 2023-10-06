package com.ml.kafka.iot;

import com.fasterxml.jackson.annotation.JsonSubTypes;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.time.Duration;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.errors.SerializationException;
import org.apache.kafka.common.serialization.Deserializer;
import org.apache.kafka.common.serialization.Serde;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.common.serialization.Serializer;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.KeyValue;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.Consumed;
import org.apache.kafka.streams.kstream.WindowedSerdes;
import org.apache.kafka.streams.kstream.Produced;
import org.apache.kafka.streams.kstream.Windowed;
import org.apache.kafka.streams.kstream.Grouped;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.KTable;
import org.apache.kafka.streams.kstream.Materialized;
import org.apache.kafka.common.utils.Bytes;
import org.apache.kafka.streams.state.WindowStore;
import org.apache.kafka.streams.kstream.TimeWindows;

import com.fasterxml.jackson.databind.JsonNode;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.streams.processor.TimestampExtractor;

import java.io.IOException;
import java.util.Map;
import java.util.Properties;
import java.util.concurrent.CountDownLatch;
import java.util.Arrays;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;


@SuppressWarnings({"WeakerAccess", "unused"})
public class AggregationEnergy {


    public static class AggregationTimestampExtractor implements TimestampExtractor {

        @Override
        public long extract(final ConsumerRecord<Object, Object> record, final long partitionTime) {
            if (record.value() instanceof BMSRawData) {
                return ((BMSRawData) record.value()).timestamp;
            }
    
            if (record.value() instanceof BMSRealtimeData) {
                return ((BMSRealtimeData) record.value()).timestamp;
            }
    
            if (record.value() instanceof JsonNode) {
                return ((JsonNode) record.value()).get("timestamp").longValue();
            }
    
            throw new IllegalArgumentException("AggregationTimestampExtractor cannot recognize the record value " + record.value());
        }
    }
    

    public static class JSONSerde<T extends JSONSerdeCompatible> implements Serializer<T>, Deserializer<T>, Serde<T> {
        private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

        @Override
        public void configure(final Map<String, ?> configs, final boolean isKey) {}

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
        public void close() {}

        @Override
        public Serializer<T> serializer() {
            return this;
        }

        @Override
        public Deserializer<T> deserializer() {
            return this;
        }
    }

    @SuppressWarnings("DefaultAnnotationParam") // being explicit for the example
    @JsonTypeInfo(use = JsonTypeInfo.Id.NAME, include = JsonTypeInfo.As.PROPERTY, property = "_t")
    @JsonSubTypes({
                      @JsonSubTypes.Type(value = BMSRawData.class, name = "bms.raw"),
                      @JsonSubTypes.Type(value = BMSDataType.class, name = "bms.type"),
                      @JsonSubTypes.Type(value = BMSRealtimeData.class, name = "bms.realtime"),

                      @JsonSubTypes.Type(value = AggregationValue.class, name = "_agg.val")
                  })
    public interface JSONSerdeCompatible {

    }


    static public class BMSRawData implements JSONSerdeCompatible {
        public String id;
        public double value;
        public String status;
        public Long timestamp;
    }

    static public class BMSDataType implements JSONSerdeCompatible {
        public String id;
        public String itemType;
    }

    static public class BMSRealtimeData implements JSONSerdeCompatible {
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
    static public class AggregationValue implements JSONSerdeCompatible {
        public int _count;
        public double value;
        public String id;
        public Long timestamp;

        @JsonIgnore
        private double[] values = {};

        public AggregationValue() {
            
        }

        public AggregationValue (int _count, double value) {
            this._count = _count;
            this.value = value;
        }

        public void setTimestamp(Long timestamp) {
            this.timestamp = timestamp;
        }

        public void setId(String id) {
            this.id = id;
        }

        public double sum (double value) {
            this._count ++;
            this.value = this.value + value;
            this.addValues(value);
            return this.value;
        }

        public double mean (int _count, double value) {
            return this.value / this._count;
        }

        @JsonIgnore
        public double[] getValues() {
            System.out.println(Arrays.toString(this.values));
            return this.values;
        }

        @JsonIgnore
        public void addValues(double value) {
            int i;
            int len = this.values.length;
            System.out.println(len);
            double newarr[] = new double[len + 1];
            for (i = 0; i< len; i++) {
                newarr[i] = this.values[i];
            }
            newarr[len] = value;
            this.values = newarr;
        }
    }

    public static void main(final String[] args) {
        final Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, "streams-bms-data");
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
        props.put(StreamsConfig.DEFAULT_TIMESTAMP_EXTRACTOR_CLASS_CONFIG, AggregationTimestampExtractor.class);
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.DEFAULT_WINDOWED_KEY_SERDE_INNER_CLASS, JSONSerde.class);
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.DEFAULT_WINDOWED_VALUE_SERDE_INNER_CLASS, JSONSerde.class);
        props.put(StreamsConfig.CACHE_MAX_BYTES_BUFFERING_CONFIG, 0);
        props.put(StreamsConfig.COMMIT_INTERVAL_MS_CONFIG, 1000);



        // setting offset reset to earliest so that we can re-run the demo code with the same pre-loaded data
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");


        final StreamsBuilder builder = new StreamsBuilder();


        final KStream<BMSDataType, BMSRawData> bmsData = builder.stream("streams-bms-raw-data", Consumed.with(new JSONSerde<>(), new JSONSerde<>()));

        final KStream<BMSDataType, AggregationValue> bmsProcess =  bmsData
        // .map((key, value) -> new KeyValue<>(key.id, value))
        .filter((key, value) -> "elec".equals(key.itemType))
        .groupByKey(Grouped.with(new JSONSerde<>(), new JSONSerde<>()))
        .windowedBy(TimeWindows.ofSizeWithNoGrace(Duration.ofSeconds(1000)).advanceBy(Duration.ofSeconds(1000)))
        .aggregate(
        () -> new AggregationValue(0, 0.0),
        (aggKey, newValue, aggValue) -> {
            System.out.println(aggKey);
            System.out.println(aggValue.getValues());
            aggValue.sum(newValue.value);
            return aggValue;
        },
        Materialized.<BMSDataType, AggregationValue, WindowStore<Bytes, byte[]>>as("time-windowed-aggregated-stream-store-bms-data") /* state store name */
        // .withValueSerde(AggregationValue.class) /* serde for aggregate value */
        )
        .toStream()
        .map((key, value) -> {
            System.out.println(key);
            System.out.println(value.getValues());
            value.setTimestamp(key.window().start());
            value.setId(key.key().id);
            return KeyValue.pair(key.key(), value);
        });

        bmsProcess.to("streams-bms-realtime-data");

        // bmsData.to("streams-bms-realtime-data", Produced.with(new JSONSerde<>(), new JSONSerde<>()));


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
