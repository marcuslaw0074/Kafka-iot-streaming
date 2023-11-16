package com.ml.kafka.stream.processor;

import java.time.Duration;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

import org.apache.kafka.streams.KeyValue;
import org.apache.kafka.streams.processor.Cancellable;
import org.apache.kafka.streams.processor.PunctuationType;
import org.apache.kafka.streams.processor.api.Processor;
import org.apache.kafka.streams.processor.api.ProcessorContext;
import org.apache.kafka.streams.processor.api.Record;
import org.apache.kafka.streams.state.KeyValueIterator;
import org.apache.kafka.streams.state.KeyValueStore;

import com.ml.kafka.model.bms.BMSDataType;
import com.ml.kafka.model.bms.BMSDeltaData;
import com.ml.kafka.model.bms.BMSEtlData;
import com.ml.kafka.model.bms.BMSMapData;
import com.ml.kafka.stream.punctuator.EtlGroupPunctuator;

final public class EtlGroupProcessor implements Processor<BMSDataType, BMSDeltaData, BMSDataType, BMSMapData> {
    private KeyValueStore<BMSDataType, BMSDeltaData> kvStore;
    // private KeyValueStore<String, Long> ckvStore;

    private ProcessorContext<BMSDataType, BMSMapData> context;

    private String stateStoreName;
    // private String contextStateStoreName;

    private Cancellable can = null;

    public EtlGroupProcessor(String stateStoreName, String contextStateStoreName) {
        // System.out.println(this);
        this.stateStoreName = stateStoreName;
        // this.contextStateStoreName = contextStateStoreName;
    }

    @Override
    public void init(final ProcessorContext<BMSDataType, BMSMapData> context) {
        this.context = context;
        kvStore = this.context.getStateStore(this.stateStoreName);
        // ckvStore = this.context.getStateStore(this.contextStateStoreName);
        // ckvStore.delete("context-lock");
        // System.out.println(ckvStore.get("context-lock"));
        System.out.println(String.format("Initialized proccessor with state store: %s", this.stateStoreName));
        // System.out.println(
        // String.format("Initialized proccessor with context state store: %s",
        // this.contextStateStoreName));
    }

    public void cancel() {
        this.can.cancel();
        this.can = null;
    }

    @Override
    public void process(final Record<BMSDataType, BMSDeltaData> record) {
        System.out.println("Start");
        final BMSDeltaData oldValue = kvStore.get(record.key());
        final BMSDeltaData currentValue = record.value();

        if (oldValue == null) {
            kvStore.put(record.key(), record.value());
        } else {
            // Long lock = this.ckvStore.get("context-lock");
            if (this.can == null) {
                // this.ckvStore.put("context-lock", 1L);
                this.can = this.context.schedule(Duration.ofSeconds(10), PunctuationType.WALL_CLOCK_TIME,
                        new EtlGroupPunctuator() {
                            @Override
                            public void punctuate(long timestamp) {
                                KeyValueIterator<BMSDataType, BMSDeltaData> kvIter = kvStore.all();
                                Record<BMSDataType, BMSMapData> newRecord;
                                // ckvStore.delete("context-lock");
                                HashMap<String, Double> map = new HashMap<>();
                                List<BMSEtlData> d = new ArrayList<>();
                                KeyValue<BMSDataType, BMSDeltaData> var;
                                while (kvIter.hasNext()) {
                                    var = kvIter.next();
                                    System.out.print(var.value);
                                    map.put(var.value.id, var.value.value);
                                    d.add(new BMSEtlData(var.value.id, var.value.value, timestamp, 1));
                                }
                                BMSMapData m = new BMSMapData(map, timestamp, 1, d);
                                kvIter.close();
                                newRecord = new Record<BMSDataType, BMSMapData>(new BMSDataType(), m, timestamp);
                                context.forward(newRecord);
                                cancel();
                                System.out.println("TEST");
                                context.commit();
                            }
                        });
            } else {

            }
            kvStore.put(record.key(), currentValue);
        }
    }

    @Override
    public void close() {
        kvStore.close();
    }
}
