package com.ml.kafka.model.bms.json;

import com.fasterxml.jackson.annotation.JsonSubTypes;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import com.ml.kafka.model.bms.*;

@SuppressWarnings("DefaultAnnotationParam") // being explicit for the example
@JsonTypeInfo(use = JsonTypeInfo.Id.NAME, include = JsonTypeInfo.As.PROPERTY, property = "_t")
@JsonSubTypes({
        @JsonSubTypes.Type(value = BMSRawData.class, name = "bms.raw"),
        @JsonSubTypes.Type(value = BMSDataType.class, name = "bms.type"),
        @JsonSubTypes.Type(value = BMSRealtimeData.class, name = "bms.realtime"),
        @JsonSubTypes.Type(value = BMSAggregationData.class, name = "bms.aggregation"),
        @JsonSubTypes.Type(value = BMSEtlData.class, name = "bms.etl")
})
public interface JSONSerdeCompatible {

}