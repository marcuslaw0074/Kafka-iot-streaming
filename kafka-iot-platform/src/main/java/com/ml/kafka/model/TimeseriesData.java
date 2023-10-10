package com.ml.kafka.model;

import com.ml.kafka.model.bms.BMSData;

public class TimeseriesData extends BMSData {
    public TimeseriesDataPoint[] data;

    static public class TimeseriesDataPoint extends BMSData{
        public String id;
        public String value;
        public String time;
    }
}
