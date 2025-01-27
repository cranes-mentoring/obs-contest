#!/bin/bash

mongosh <<EOF
rs.initiate({
  _id: "rs0",
  members: [
  { _id: 0, host: "mongo-secondary:27017" },
  { _id: 1, host: "mongo-secondary-2:27017" }
  ]
});
EOF


cfg = rs.conf()
cfg.members = [
  { _id: 0, host: "mongo-secondary:27017" },
  { _id: 1, host: "mongo-secondary-2:27017" }
]
rs.reconfig(cfg, { force: true })