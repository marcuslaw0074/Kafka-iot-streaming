package com.ml.kafka.stream.processor;

import java.time.Duration;

import org.apache.kafka.streams.processor.PunctuationType;
import org.apache.kafka.streams.processor.Punctuator;
import org.apache.kafka.streams.processor.api.Processor;
import org.apache.kafka.streams.processor.api.ProcessorContext;
import org.apache.kafka.streams.processor.api.Record;
import org.apache.kafka.streams.state.KeyValueStore;

import com.ml.kafka.model.bms.BMSDataType;
import com.ml.kafka.model.bms.BMSDeltaData;

final public class EtlGroupProcessor implements Processor<BMSDataType, BMSDeltaData, BMSDataType, BMSDeltaData> {
    private KeyValueStore<BMSDataType, BMSDeltaData> kvStore;

    private ProcessorContext<BMSDataType, BMSDeltaData> context;

    private String stateStoreName;

    public EtlGroupProcessor(String stateStoreName) {
        this.stateStoreName = stateStoreName;
    }

    @Override
    public void init(final ProcessorContext<BMSDataType, BMSDeltaData> context) {
        this.context = context;
        kvStore = this.context.getStateStore(this.stateStoreName);
        System.out.println(String.format("Initialized proccessor with state store: %s", this.stateStoreName));
    }

    @Override
    public void process(final Record<BMSDataType, BMSDeltaData> record) {
        System.out.println("Start");
        final BMSDeltaData oldValue = kvStore.get(record.key());
        final BMSDeltaData currentValue = record.value();
        System.out.println("group");
        System.out.println(oldValue);
        System.out.println(currentValue);

        if (oldValue == null) {
            kvStore.put(record.key(), record.value());
        } else {
            this.context.schedule(Duration.ofSeconds(10), PunctuationType.WALL_CLOCK_TIME, new Punctuator() {
                @Override
                public void punctuate(long timestamp) {
                    // TODO Auto-generated method stub
                }
            });
            kvStore.put(record.key(), currentValue);
        }
    }

    @Override
    public void close() {
        kvStore.close();
    }
}
