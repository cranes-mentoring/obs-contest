curl -X POST -H "Content-Type: application/json" --data '{
  "name": "mongo-connector",
  "config": {
    "connector.class": "io.debezium.connector.mongodb.MongoDbConnector",
    "mongodb.connection.string": "mongodb://example_user:example_password@mongo:27017",
    "mongodb.user": "example_user",
    "mongodb.password": "example_password",
    "mongodb.name": "project_one",
    "mongodb.servers.name": "mongo",
    "database.include.list": "purchase",
    "collection.include.list": "project_one.purchase",
    "tasks.max": "1",
    "topic.prefix": "mongo",
    "schema.history.internal.kafka.bootstrap.servers": "kafka:9092",
    "schema.history.internal.kafka.topic": "schema-changes.mongo"
  }
}' http://localhost:8083/connectors
