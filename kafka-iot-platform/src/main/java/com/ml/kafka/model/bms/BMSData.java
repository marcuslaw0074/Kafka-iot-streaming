package com.ml.kafka.model.bms;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;

public class BMSData {
    @Override
    public String toString() {
        Gson gson = new GsonBuilder().setPrettyPrinting().create();
        return gson.toJson(this);
    }
}
