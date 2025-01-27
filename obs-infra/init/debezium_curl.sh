curl -X POST -H "Content-Type: application/json" --data '{
  "name": "mongo-connector",
  "config": {
    "connector.class": "io.debezium.connector.mongodb.MongoDbConnector",
    "mongodb.connection.string": "mongodb://mongo-primary:27017,mongo-secondary:27017,mongo-secondary-2:27017/?replicaSet=rs0",
    "mongodb.user": "example_user",
    "mongodb.password": "example_password",
    "mongodb.name": "project_one",
    "database.include.list": "purchase",
    "collection.include.list": "project_one.purchase",
    "tasks.max": "1",
    "topic.prefix": "mongo",
    "schema.history.internal.kafka.bootstrap.servers": "kafka:9092",
    "schema.history.internal.kafka.topic": "schema-changes.mongo"
  }
}' http://localhost:8083/connectors
