package com.ml.kafka.stream.processor;

import java.time.Duration;

import org.apache.kafka.streams.processor.Cancellable;
import org.apache.kafka.streams.processor.PunctuationType;
import org.apache.kafka.streams.processor.Punctuator;
import org.apache.kafka.streams.processor.api.Processor;
import org.apache.kafka.streams.processor.api.ProcessorContext;
import org.apache.kafka.streams.processor.api.Record;
import org.apache.kafka.streams.state.KeyValueIterator;
import org.apache.kafka.streams.state.KeyValueStore;

public abstract class CustomGroupProcessor<KIn, VIn, KOut, VOut>
        implements Processor<KIn, VIn, KOut, VOut> {

    private ProcessorContext<KOut, VOut> context;
    private String stateStoreName;
    private KeyValueStore<KIn, VIn> kvStore;
    private Cancellable can = null;

    public abstract Record<KOut, VOut> generatRecord(KeyValueIterator<KIn, VIn> iter);

    @Override
    public void init(ProcessorContext<KOut, VOut> context) {
        Processor.super.init(context);
        this.context = context;
        this.kvStore = this.context.getStateStore(this.stateStoreName);
        System.out.println(String.format("Initialized CustomGroupProcessor with state store: %s", this.stateStoreName));
    }

    @Override
    public void process(Record<KIn, VIn> record) {
        System.out.println("CustomGroupProcessor::process start");
        final VIn oldValue = this.kvStore.get(record.key());

        if (oldValue == null) {
            this.kvStore.put(record.key(), record.value());
        } else {
            this.addScheduler();
        }
    }

    public void cancel() {
        this.can.cancel();
        this.can = null;
    }

    public synchronized void addScheduler() {
        if (this.can == null) {
            this.can = this.context.schedule(Duration.ofSeconds(10), PunctuationType.WALL_CLOCK_TIME, new Punctuator() {
                @Override
                public void punctuate(long timestamp) {
                    KeyValueIterator<KIn, VIn> kvIter = kvStore.all();
                    Record<KOut, VOut> newRecord = generatRecord(kvIter);
                    kvIter.close();
                    context.forward(newRecord);
                    context.commit();
                    cancel();
                }
            });
        }
    }

    @Override
    public void close() {
        this.kvStore.close();
    }

    public CustomGroupProcessor(String stateStoreName) {
        this.stateStoreName = stateStoreName;
    }
}
