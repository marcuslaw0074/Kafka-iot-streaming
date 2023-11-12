package com.ml.kafka.serialization.timestamp;

import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.Date;

final public class UnixTimestampConvertor {

    private Long unix;
    private Date date;
    private Calendar c;
    private int timezone;

    public UnixTimestampConvertor(Long unix, int timezone) {
        this.timezone = timezone;
        this.validUnix(unix);
        this.convertUnixToDate();
    }

    private void validUnix(Long unix) {
        int numDigits = unix.toString().length();
        switch (numDigits) {
            case 10:
                this.unix = unix * 1000;
                break;
            case 13:
                this.unix = unix;
                break;
            case 16:
                this.unix = unix / 1000;
                break;
            default:
                break;
        }
    }

    private Date convertUnixToDate() {
        if (this.unix != null) {
            this.date = new java.util.Date((long) this.unix);
            this.c = Calendar.getInstance();
            this.c.setTime(this.date);
        }
        return this.date;
    }

    public int getDayOfMonth() {
        if (this.date != null) {
            return this.c.get(Calendar.DAY_OF_MONTH);
        } else {
            return -1;
        }

    }

    public String getDateString() {
        if (this.date != null) {
            return (new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSS'Z'")).format(this.date);
        } else {
            return "";
        }

    }

    public boolean isNextDay(UnixTimestampConvertor obj) {
        this.c.add(Calendar.HOUR_OF_DAY, +this.timezone);
        obj.c.add(Calendar.HOUR_OF_DAY, +this.timezone);

        if (this.c.get(Calendar.DAY_OF_MONTH) > obj.c.get(Calendar.DAY_OF_MONTH)) {
            return true;
        } else {
            return false;
        }
    }

    public boolean isNextHour(UnixTimestampConvertor obj) {
        if (this.c.get(Calendar.HOUR_OF_DAY) > obj.c.get(Calendar.HOUR_OF_DAY)) {
            return true;
        } else {
            return false;
        }
    }

    public boolean isMonday() {
        if (this.c.get(Calendar.DAY_OF_WEEK) == 1) {
            return true;
        } else {
            return false;
        }
    }

}
