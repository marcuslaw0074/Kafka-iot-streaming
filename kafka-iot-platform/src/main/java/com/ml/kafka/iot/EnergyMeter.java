package com.ml.kafka.iot;

import java.time.Duration;
import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;
import java.util.Properties;
import java.util.concurrent.CountDownLatch;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.serialization.Serde;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.KeyValue;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.Consumed;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.Produced;

import com.ml.kafka.model.*;
import com.ml.kafka.model.TimeseriesData.TimeseriesDataPoint;
import com.ml.kafka.model.bms.*;
import com.ml.kafka.model.bms.json.JSONDeserializer;
import com.ml.kafka.model.bms.json.JSONSerde;
import com.ml.kafka.model.bms.json.JSONSerializer;

@SuppressWarnings({ "WeakerAccess", "unused" })
public class EnergyMeter {

    public static void main(final String[] args) {
        final String application_id = "test-kafka-stream-etl-energy";
        final String stream_name = "test-kafka-streams-energy-raw-data";
        final String sink_name = "test-kafka-streams-energy-realtime-data";
        final String server = "localhost:9092";

        final Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, application_id);
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, server);
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, JSONSerde.class);
        props.put(StreamsConfig.COMMIT_INTERVAL_MS_CONFIG, 1000);
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");

        final StreamsBuilder builder = new StreamsBuilder();

        JSONDeserializer<TimeseriesData> timeseriesJsonDeserializer = new JSONDeserializer<>(TimeseriesData.class);
        JSONSerializer<TimeseriesData> timeseriesJsonSerializer = new JSONSerializer<>();

        Serde<TimeseriesData> timeseriesSerde = Serdes.serdeFrom(timeseriesJsonSerializer, timeseriesJsonDeserializer);

        final KStream<BMSDataType, TimeseriesData> timeData = builder.stream(stream_name,
                Consumed.with(new JSONSerde<>(), timeseriesSerde));

        timeData
                .map((key, value) -> {
                    System.out.println(key);
                    System.out.println(value.toString());
                    return KeyValue.pair(key, value);
                })
                .flatMap((key, value) -> {
                    System.out.println(key);
                    System.out.println(value.toString());
                    List<KeyValue<BMSDataType, TimeseriesDataPoint>> result = new LinkedList<>();
                    for (int i = 0; i < value.data.length; i++) {
                        result.add(KeyValue.pair(key, value.data[i]));
                        System.out.println(key);
                        System.out.println(value.data[i].toString());
                    }
                    return result;
                }).map((key, value) -> {
                    System.out.println("HAHA");
                    System.out.println(key);
                    System.out.println(value.toString());
                    return KeyValue.pair(key, value);
                }).to(sink_name);


        final KafkaStreams streams = new KafkaStreams(builder.build(), props);
        final CountDownLatch latch = new CountDownLatch(1);

        Runtime.getRuntime().addShutdownHook(new Thread("test-streams-energy-shutdown-hook") {
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
