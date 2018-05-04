# youtubeat

Welcome to youtubeat.

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/bienkma/youtubeat`

## Getting Started with youtubeat

### Requirements

* [Golang](https://golang.org/dl/) 1.10

### Build

To build the binary for youtubeat run the command below. This will generate a binary
in the same directory with the name youtubeat.

```
make
```

### Requirement 
Install kafka and elasticsearch or logstash on server
```
$ mkdir demo && cd demo
$ wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-5.6.3.tar.gz
$ wget http://mirror.downloadvn.com/apache/kafka/1.1.0/kafka_2.12-1.1.0.tgz
$ tar -xvf kafka_2.12-1.1.0.tgz
$ tar -xvf elasticsearch-5.6.3.tar.gz
$ demo/kafka_2.12-1.1.0/bin/zookeeper-server-start.sh -daemon demo/kafka_2.12-1.1.0/config/zookeeper.properties
$ demo/kafka_2.12-1.1.0/bin/kafka-server-start.sh --daemon demo/kafka_2.12-1.1.0/config/server.properties
$ elasticsearch-5.6.3/bin/elasticsearch -d
```

### Run

To run youtubeat with debugging output enabled, run:

```
./youtubeat -c youtubeat.yml -e -d "*"
kafka_2.12-1.1.0/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic youtube
> topmusic
> any_key_word
```
### Check index
To query all docs on elasticsearch
```
http://localhost:9200/_cat/indices?v
http://localhost:9200/keywordlog-6.0.0-beta1-2018.05.04/_search?pretty=true&q=*:*
```

