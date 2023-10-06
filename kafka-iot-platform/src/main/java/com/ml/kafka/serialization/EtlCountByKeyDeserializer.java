package com.ml.kafka.serialization;

import org.apache.kafka.common.errors.SerializationException;
import org.apache.kafka.common.serialization.Deserializer;

public class EtlCountByKeyDeserializer implements Deserializer<Long> {
    public Long deserialize(String topic, byte[] data) {
        if (data == null)
            return null;
        if (data.length != 8) {
            throw new SerializationException("Size of data received by LongDeserializer is not 8");
        }

        long value = 0;
        for (byte b : data) {
            value <<= 8;
            value |= b & 0xFF;
        }
        return value;
    }
}