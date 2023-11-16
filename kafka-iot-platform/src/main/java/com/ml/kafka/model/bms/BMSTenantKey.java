package com.ml.kafka.model.bms;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.ml.kafka.model.bms.json.JSONSerdeCompatible;

@JsonIgnoreProperties(ignoreUnknown = true)
public class BMSTenantKey extends BMSData implements JSONSerdeCompatible {
    public String tenantId;
    public String projectId;

    public BMSTenantKey() {
    }

    public BMSTenantKey(String tenantId, String projectId) {
        this.tenantId = tenantId;
        this.projectId = projectId;
    }

}
