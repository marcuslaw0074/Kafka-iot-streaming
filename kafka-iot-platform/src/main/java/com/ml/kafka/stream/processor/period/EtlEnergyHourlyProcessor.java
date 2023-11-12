package com.ml.kafka.stream.processor.period;

import org.apache.kafka.streams.processor.api.Processor;
import org.apache.kafka.streams.processor.api.ProcessorContext;
import org.apache.kafka.streams.processor.api.Record;
import org.apache.kafka.streams.state.KeyValueStore;

import com.ml.kafka.model.bms.BMSDataType;
import com.ml.kafka.model.bms.BMSDeltaData;
import com.ml.kafka.serialization.timestamp.UnixTimestampConvertor;

final public class EtlEnergyHourlyProcessor implements Processor<BMSDataType, BMSDeltaData, BMSDataType, BMSDeltaData> {
    private KeyValueStore<BMSDataType, BMSDeltaData> hourlyKvStore;

    private ProcessorContext<BMSDataType, BMSDeltaData> context;

    private String hourlyStateStoreName;


    public EtlEnergyHourlyProcessor(String hourlyStateStoreName) {
        this.hourlyStateStoreName = hourlyStateStoreName;
    }

    @Override
    public void init(final ProcessorContext<BMSDataType, BMSDeltaData> context) {
        this.context = context;
        hourlyKvStore = this.context.getStateStore(this.hourlyStateStoreName);
        System.out.println(String.format("Initialized proccessor with state store: %s", this.hourlyStateStoreName));
    }

    private void hourlyProcess(final Record<BMSDataType, BMSDeltaData> record) {
            final BMSDeltaData oldValue = hourlyKvStore.get(record.key());
            final BMSDeltaData currentValue = record.value();
            System.out.println("hourly");
            System.out.println(oldValue);
            System.out.println(currentValue);

            if (oldValue == null) {
                hourlyKvStore.put(record.key(), record.value());
                System.out.println("NULL");
            } else {
                UnixTimestampConvertor oldTime = new UnixTimestampConvertor(oldValue.timestamp, +8);
                UnixTimestampConvertor currentTime = new UnixTimestampConvertor(currentValue.timestamp, +8);

                if (currentTime.isNextHour(oldTime) && (currentValue.value >= oldValue.value)) {
                    this.context.forward(new Record<>(record.key(),
                            new BMSDeltaData(currentValue.id, currentValue.value - oldValue.value,
                                    currentValue.timestamp,
                                    currentValue.timestamp - oldValue.timestamp,
                                    1, "hourly"),
                            currentValue.timestamp));
                    hourlyKvStore.put(record.key(), currentValue);
                } else {
                    System.out.println("hourly-pass");
                }
            }
        }

    @Override
    public void process(final Record<BMSDataType, BMSDeltaData> record) {
        System.out.println("Start");

        this.hourlyProcess(record);
    }

    @Override
    public void close() {
        hourlyKvStore.close();
    }
}
