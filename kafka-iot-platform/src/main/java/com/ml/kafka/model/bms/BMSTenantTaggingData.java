package com.ml.kafka.model.bms;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.ml.kafka.model.bms.json.JSONSerdeCompatible;

@JsonIgnoreProperties(ignoreUnknown = true)
public class BMSTenantTaggingData extends BMSData implements JSONSerdeCompatible {
    public String bmsId;
    public String tenantId;
    public String itemType;
    public String projectId;
    public int status;
    public double value;
    public Long timestamp;
}