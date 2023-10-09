package com.ml.kafka.model.bms;

import com.ml.kafka.model.bms.json.JSONSerdeCompatible;

final public class BMSRealtimeData extends BMSData implements JSONSerdeCompatible {
    public String id;
    public String block;
    public String buildingName;
    public String equipmentName;
    public String functionType;
    public String prefername;
    public String status;
    public Long value;
    public Long timestamp;
}
