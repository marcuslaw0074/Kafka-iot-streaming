package com.ml.kafka.model.bms;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.ml.kafka.model.bms.json.JSONSerdeCompatible;


@JsonIgnoreProperties(ignoreUnknown = true)
final public class BMSDeltaData extends BMSData implements JSONSerdeCompatible {
    public double value;
    public String id;
    public String timeType;
    public Long timestamp;
    public Long timestampDelta;
    public int status;

    public BMSDeltaData() {
    }

    public BMSDeltaData(String id, double value, Long timestamp, Long timestampDelta, int status, String timeType) {
        this.id = id;
        this.value = value;
        this.timestamp = timestamp;
        this.timestampDelta = timestampDelta;
        this.status = status;
        this.timeType = timeType;
    }
}
