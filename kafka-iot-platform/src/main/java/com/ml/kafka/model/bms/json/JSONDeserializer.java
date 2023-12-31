package com.ml.kafka.model.bms.json;

import com.google.gson.Gson;
import org.apache.kafka.common.serialization.Deserializer;

import java.util.Map;

public class JSONDeserializer<T> implements Deserializer<T> {

    private Gson gson = new Gson();
    private Class<T> deserializedClass;

    public JSONDeserializer(Class<T> deserializedClass) {
        this.deserializedClass = deserializedClass;
    }

    public JSONDeserializer() {
    }

    @Override
    @SuppressWarnings("unchecked")
    public void configure(Map<String, ?> map, boolean b) {
        if(deserializedClass == null) {
            deserializedClass = (Class<T>) map.get("serializedClass");
        }
    }

    @Override
    public T deserialize(String s, byte[] bytes) {
         if(bytes == null){
             return null;
         }

         return gson.fromJson(new String(bytes),deserializedClass);

    }

    @Override
    public void close() {

    }
}
