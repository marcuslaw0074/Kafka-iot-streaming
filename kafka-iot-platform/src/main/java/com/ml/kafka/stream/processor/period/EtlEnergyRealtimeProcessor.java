package com.ml.kafka.stream.processor.period;

import org.apache.kafka.streams.processor.api.Processor;
import org.apache.kafka.streams.processor.api.ProcessorContext;
import org.apache.kafka.streams.processor.api.Record;
import org.apache.kafka.streams.state.KeyValueStore;

import com.ml.kafka.model.bms.BMSDataType;
import com.ml.kafka.model.bms.BMSDeltaData;

final public class EtlEnergyRealtimeProcessor implements Processor<BMSDataType, BMSDeltaData, BMSDataType, BMSDeltaData> {
    private KeyValueStore<BMSDataType, BMSDeltaData> realtimeKvStore;

    private ProcessorContext<BMSDataType, BMSDeltaData> context;

    private String realtimeStateStoreName;


    public EtlEnergyRealtimeProcessor(String realtimeStateStoreName) {
        this.realtimeStateStoreName = realtimeStateStoreName;
    }

    @Override
    public void init(final ProcessorContext<BMSDataType, BMSDeltaData> context) {
        this.context = context;
        realtimeKvStore = this.context.getStateStore(this.realtimeStateStoreName);
        System.out.println(String.format("Initialized proccessor with state store: %s", this.realtimeStateStoreName));
    }

    private void realtimeProcess(final Record<BMSDataType, BMSDeltaData> record) {
        final BMSDeltaData oldValue = realtimeKvStore.get(record.key());
        final BMSDeltaData currentValue = record.value();
        System.out.println("realtime"); 
        System.out.println(oldValue);
        System.out.println(currentValue);

        if (oldValue == null) {
            realtimeKvStore.put(record.key(), record.value());
        } else {
            if (currentValue.value >= oldValue.value) {
                this.context.forward(new Record<>(record.key(),
                        new BMSDeltaData(currentValue.id, currentValue.value - oldValue.value,
                                currentValue.timestamp,
                                currentValue.timestamp - oldValue.timestamp,
                                1, "realtime"),
                        currentValue.timestamp));
                realtimeKvStore.put(record.key(), currentValue);
            } else {
                System.out.println("realtime-pass");
            }
        }

    }

    @Override
    public void process(final Record<BMSDataType, BMSDeltaData> record) {
        System.out.println("Start");

        this.realtimeProcess(record);
    }

    @Override
    public void close() {
        realtimeKvStore.close();
    }
}
