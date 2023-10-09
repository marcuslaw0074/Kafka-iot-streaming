package com.ml.kafka.model.bms;

import com.ml.kafka.model.bms.json.JSONSerdeCompatible;

final public class BMSDataType extends BMSData implements JSONSerdeCompatible {
    public String id;
    public String itemType;
}
