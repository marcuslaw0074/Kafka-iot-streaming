package com.ml.kafka.stream.consumer;

import java.time.Duration;
import java.util.Arrays;
import java.util.Properties;

import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;

public class ConsumerEtlCountByKey {
    public static void main(String[] args) {
        final String sink_name = "kafka-sink-count-output";

        Properties props = new Properties();
        props.setProperty("bootstrap.servers", "localhost:9092");
        props.setProperty("group.id", "test");
        props.setProperty("enable.auto.commit", "true");
        props.setProperty("auto.commit.interval.ms", "1000");
        props.setProperty("key.deserializer", "org.apache.kafka.common.serialization.StringDeserializer");
        props.setProperty("value.deserializer", "com.ml.kafka.serialization.EtlCountByKeyDeserializer");
        
        KafkaConsumer<String, Long> consumer = new KafkaConsumer<>(props);
        consumer.subscribe(Arrays.asList(sink_name));
        while (true) {
            ConsumerRecords<String, Long> records = consumer.poll(Duration.ofMillis(100));
            for (ConsumerRecord<String, Long> record : records)
                System.out.printf("offset = %d, key = %s, value = %d%n", record.offset(), record.key(), record.value());
        }
    }
}
