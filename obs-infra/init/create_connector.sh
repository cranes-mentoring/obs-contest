curl -X POST -H "Content-Type: application/json" --data '{
  "name": "mongo-connector",
  "config": {
    "connector.class": "io.debezium.connector.mongodb.MongoDbConnector",
    "mongodb.connection.string": "mongodb://mongo1:27017,mongo2:27017,mongo3:27017/?replicaSet=rs0",
    "mongodb.user": "root",
    "mongodb.password": "example",
    "database.include.list": "project_one",
    "collection.include.list": "project_one.purchases",
    "tasks.max": "1",
    "topic.prefix": "mongo",
    "schema.history.internal.kafka.bootstrap.servers": "kafka:9092",
    "schema.history.internal.kafka.topic": "schema-changes.mongo",
    "transforms": "route",
    "transforms.route.type" : "org.apache.kafka.connect.transforms.RegexRouter",
    "transforms.route.regex" : "([^.]+)\\.([^.]+)\\.([^.]+)",
    "transforms.route.replacement" : "$3"
  }
}' http://localhost:8083/connectors
