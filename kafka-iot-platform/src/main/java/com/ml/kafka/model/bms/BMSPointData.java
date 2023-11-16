package com.ml.kafka.model.bms;

public class BMSPointData extends BMSData {
    public String id;
    public double value;
    public Long timestamp;

    public BMSPointData() {

    }

    public BMSPointData(String id, double value, Long timestamp) {
        this.id = id;
        this.value = value;
        this.timestamp = timestamp;
    }
}
