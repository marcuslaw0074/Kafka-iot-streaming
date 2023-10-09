package com.ml.kafka.model.bms;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.ml.kafka.model.bms.json.JSONSerdeCompatible;

@JsonIgnoreProperties(ignoreUnknown = true)
final public class BMSEtlData extends BMSData implements JSONSerdeCompatible {
    public double value;
    public String id;
    public Long timestamp;
    public int status;

    public BMSEtlData() {
    }

    public BMSEtlData(String id, double value, Long timestamp, int status) {
        this.id = id;
        this.value = value;
        this.timestamp = timestamp;
        this.status = status;
    }
}