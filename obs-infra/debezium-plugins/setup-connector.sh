#!/bin/bash

curl -LO https://repo1.maven.org/maven2/io/debezium/debezium-mongodb/connector/1.5.0.Final/debezium-mongodb-connector-1.5.0.Final.jar

PLUGIN_PATH=/kafka/connect/plugins

mkdir -p $PLUGIN_PATH
cp debezium-mongodb-connector-1.5.0.Final.jar $PLUGIN_PATH

echo "done!"