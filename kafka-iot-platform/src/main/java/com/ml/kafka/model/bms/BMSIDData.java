package com.ml.kafka.model.bms;

public class BMSIDData extends BMSData {
    public String id;
    public String itemType;

    public BMSIDData() {

    }

    public BMSIDData(String id, String itemType) {
        this.id = id;
        this.itemType = itemType;
    }
}
