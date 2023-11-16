package com.ml.kafka.model.bms;

import java.util.HashMap;
import java.util.List;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.ml.kafka.model.bms.json.JSONSerdeCompatible;


@JsonIgnoreProperties(ignoreUnknown = true)
public class BMSMapData extends BMSData implements JSONSerdeCompatible {
    public HashMap<String, Double> map;
    public List<BMSEtlData> data;
    public Long timestamp;
    public int status;

    public BMSMapData() {
    }

    public BMSMapData(HashMap<String, Double> map, Long timestamp, int status) {
        this.map = map;
        this.status = status;
        this.timestamp = timestamp;
    }

    public BMSMapData(HashMap<String, Double> map, Long timestamp, int status, List<BMSEtlData> data) {
        this.map = map;
        this.status = status;
        this.timestamp = timestamp;
        this.data = data;
    }

}
