package com.ml.kafka.model.bms;

import com.ml.kafka.model.bms.json.JSONSerdeCompatible;

final public class BMSRawData extends BMSData implements JSONSerdeCompatible {
    public String id;
    public double value;
    public String status;
    public Long timestamp;
}
