FROM maven:3.5.2-jdk-8-alpine AS MAVEN_TOOL_CHAIN


COPY pom.xml /root/kafka-iot-platform/
COPY src /root/kafka-iot-platform/src/
WORKDIR /root/kafka-iot-platform/
RUN mvn clean package

CMD [ "mvn", "exec:java", "-Dexec.mainClass=com.ml.kafka.etl.energy.EtlEnergyRealtime"]