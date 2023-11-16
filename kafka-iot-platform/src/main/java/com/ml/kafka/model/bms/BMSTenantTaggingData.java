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

    public BMSTenantTaggingData() {
    }

    public BMSTenantTaggingData(String bmsId, String tenantId, String itemType, String projectId, int status,
            double value, Long timestamp) {
        this.bmsId = bmsId;
        this.tenantId = tenantId;
        this.itemType = itemType;
        this.projectId = projectId;
        this.status = status;
        this.value = value;
        this.timestamp = timestamp;
    }

}