# FROM ubuntu

# RUN apt-get update -y
# RUN apt-get -y install python3-pip
# RUN apt-get install -y python3


# FROM golang:python3

# COPY ./ /root/kafka-iot-connect/


# WORKDIR /root/kafka-iot-connect


FROM golang:1.19

WORKDIR /root/kafka-iot-connect

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build .

CMD ["./kafka-iot-connect"]