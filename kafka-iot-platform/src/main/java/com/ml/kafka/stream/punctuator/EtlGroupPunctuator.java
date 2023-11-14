package com.ml.kafka.stream.punctuator;

import org.apache.kafka.streams.processor.Cancellable;
import org.apache.kafka.streams.processor.Punctuator;

public class EtlGroupPunctuator implements Punctuator {
    Cancellable schedule;
    boolean first_run;

    public void punctuate(final long timestamp) {
        if (!first_run) {
            schedule.cancel();
            first_run = true;
        }
    }
}
