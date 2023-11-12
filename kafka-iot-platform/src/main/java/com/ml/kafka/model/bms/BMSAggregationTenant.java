package com.ml.kafka.model.bms;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.ml.kafka.model.bms.json.JSONSerdeCompatible;

@JsonIgnoreProperties(ignoreUnknown = true)
public class BMSAggregationTenant extends BMSData implements JSONSerdeCompatible {
    public String type;
    public List<BMSTenantTaggingData> arr;
    public HashMap<String, Double> map;
    public Long timestamp;
    public int status;

    public BMSAggregationTenant() {}


    public BMSAggregationTenant(String type, Long timestamp, int status) {
        this.type = type;
        this.timestamp = timestamp;
        this.status = status;
        this.arr = new ArrayList<BMSTenantTaggingData>();
        this.map = new HashMap<String, Double>();
    }


}
