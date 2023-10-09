package com.ml.kafka.model.bms;

import java.util.Arrays;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.ml.kafka.model.bms.json.JSONSerdeCompatible;

@JsonIgnoreProperties(ignoreUnknown = true)
final public class BMSAggregationData extends BMSData implements JSONSerdeCompatible {
    public int count;
    public double value;
    public String type;
    public String id;
    public Long timestamp;
    public int status;

    // @JsonIgnore
    public double[] values = {};

    public BMSAggregationData() {
    }

    public BMSAggregationData(int count, double value) {
        this.count = count;
        this.value = value;
    }

    public void setTimestamp(Long timestamp) {
        this.timestamp = timestamp;
    }

    public void setId(String id) {
        this.id = id;
    }

    public void setType(String type) {
        this.type = type;
    }

    // @JsonIgnore
    public double[] getValues() {
        // System.out.println(Arrays.toString(this.values));
        return this.values;
    }

    // @JsonIgnore
    public void addValues(double value) {
        int i;
        int len = this.values.length;
        // System.out.println(len);
        double newarr[] = new double[len + 1];
        for (i = 0; i < len; i++) {
            newarr[i] = this.values[i];
        }
        newarr[len] = value;
        this.count++;
        this.values = newarr;
    }

    public BMSAggregationData agg() {
        double _value;
        int _status;
        switch (this.type) {
            case "mean":
                _value = Arrays.stream(this.values).sum() / this.count;
                _status = 1;
                break;
            case "sum":
                _value = Arrays.stream(this.values).sum();
                _status = 1;
                break;
            case "count":
                _value = this.count;
                _status = 1;
                break;
            case "max":
                try {
                    _value = Arrays.stream(this.values).max().getAsDouble();
                    _status = 1;
                } catch (Exception e) {
                    _value = 0;
                    _status = 0;
                }
                break;
            case "min":
                try {
                    _value = Arrays.stream(this.values).min().getAsDouble();
                    _status = 1;
                } catch (Exception e) {
                    _value = 0;
                    _status = 0;
                }
                break;
            case "first":
                _value = this.values[0];
                _status = 1;
                break;
            case "last":
                _value = this.values[this.values.length - 1];
                _status = 1;
                break;
            case "diff":
                if (this.values.length < 2) {
                    _value = 0;
                    _status = 0;
                } else {
                    _value = this.values[this.values.length - 1] - this.values[0];
                    _status = 1;
                }
                break;
            default:
                _value = 0;
                _status = 0;
                System.out.println("DEFAULT");
        }
        this.value = _value;
        this.status = _status;
        return this;
    }
}
