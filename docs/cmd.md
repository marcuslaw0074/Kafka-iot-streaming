# Command: switch to java 8 

/usr/libexec/java_home -V       

export JAVA_HOME=`/usr/libexec/java_home -v 1.8.0_292`

# Run stream app

mvn clean package
mvn exec:java -Dexec.mainClass=com.ml.kafka.etl.EtlJson


# create topics for inputs

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-pipe-input

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-linesplit-input

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-wordcount-input

# create topics for outputs

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-pipe-output

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-linesplit-output

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-wordcount-output \
--config cleanup.policy=compact

# inspect topics 

kafka-topics --bootstrap-server localhost:9092 --describe

# docker exec 

docker exec -it kafka-broker bash
cd /bin/


# Publish msg 

kafka-console-producer --broker-list localhost:9092 --topic streams-pipe-input

# subscrible msg

kafka-console-consumer --bootstrap-server localhost:9092 \
--topic streams-pipe-output \
--from-beginning \
--formatter kafka.tools.DefaultMessageFormatter \
--property print.key=true \
--property print.value=true \
--property key.deserializer=org.apache.kafka.common.serialization.StringDeserializer \
--property value.deserializer=org.apache.kafka.common.serialization.StringDeserializer













kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-pageview-input

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-userprofile-input

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-pageviewstats-typed-output


kafka-console-producer \
  --topic streams-pageview-input \
  --bootstrap-server localhost:9092 \
  --property parse.key=true \
  --property key.separator=":"

SEND: user1:{"_t":"pv","user":"user1","page":"gjrosa","timestamp":1234567}


kafka-console-producer \
  --topic streams-userprofile-input \
  --bootstrap-server localhost:9092 \
  --property parse.key=true \
  --property key.separator=":"


SEND: user1:{"_t":"up","region":"region3","timestamp":1234567}


kafka-console-consumer --bootstrap-server localhost:9092 \
--topic streams-pageviewstats-typed-output \
--from-beginning \
--formatter kafka.tools.DefaultMessageFormatter \
--property print.key=true \
--property print.value=true




kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-pageview-input


kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-userprofile-input

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-pageviewstats-typed-output



{"u1":{"_t": "pv", "user": "u1", "page": "p1", "timestamp": 162747187913}}
{"u1":{"_t": "up", "region": "r1", "timestamp": 162747187913}}








docker exec -it kafka-broker bash



kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-bms-raw-data

kafka-topics --create \
--bootstrap-server localhost:9092 \
--replication-factor 1 \
--partitions 1 \
--topic streams-bms-realtime-data


kafka-console-producer \
  --topic streams-bms-raw-data \
  --bootstrap-server localhost:9092 \
  --property parse.key=true \
  --property key.separator=":::"


SEND: {"_t":"bms.type","id":"user1","itemType":"ukfs30"}:::{"_t":"bms.raw","id":"user1","value":31.431,"timestamp":1234567,"status":"1"}

{"_t":"bms.type","id":"user2","itemType":"ukfs30"}:::{"_t":"bms.raw","id":"user1","value":31.431,"timestamp":12934567,"status":"1"}

{"_t":"bms.type","id":"user2","itemType":"ukfs30"}:::{"_t":"bms.raw","id":"user1","value":31.431,"timestamp":15934567,"status":"1"}

{"_t":"bms.type","id":"user2","itemType":"elec"}:::{"_t":"bms.raw","id":"user1","value":31.431,"timestamp":15934567,"status":"1"}


kafka-console-consumer --bootstrap-server localhost:9092 \
--topic streams-bms-realtime-data \
--from-beginning \
--formatter kafka.tools.DefaultMessageFormatter \
--property print.key=true \
--property print.value=true