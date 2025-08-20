#!/bin/bash
set -e

# Configurations -- UPDATE_ME
MONGO_HOST="localhost"
MONGO_PORT="27017"
MONGO_ADMIN_USER="mongo" 
MONGO_ADMIN_PASS='mongo'
APP_PWD='mongo'
PURGER_PWD='mongo'

if [ -z "$APP_PWD" ] || [ -z "$PURGER_PWD" ]; then
    echo "Error: All passwords must be set"
    exit 1
fi

DB_NAME="shop"
# DB Changes
if ! command -v mongosh &> /dev/null; then
  wget https://downloads.mongodb.com/compass/mongodb-mongosh_2.5.6_amd64.deb
  sudo apt install ./mongodb-mongosh_2.5.6_amd64.deb -y
  rm -f ./mongodb-mongosh_2.5.6_amd64.deb
  mongosh --version
fi

mongosh --host "$MONGO_HOST" --port "$MONGO_PORT" \
    -u "$MONGO_ADMIN_USER" -p "$MONGO_ADMIN_PASS" --authenticationDatabase "admin" <<EOF
use ${DB_NAME}
if (!db.getRole("appRole")) {
    db.createRole({
      role: "appRole",
      privileges: [
        {
          resource: { db: "${DB_NAME}", collection: "" },
          actions: [
            "analyze",
            "bypassDocumentValidation",
            "changeOwnCustomData",
            "changeOwnPassword",
            "changeStream",
            "collStats",
            "dbStats",
            "find",
            "insert",
            "killCursors",
            "listCollections",
            "listIndexes",
            "listSearchIndexes",
            "update",
          ]
        }
      ],
      roles: []
    }, { w: "majority" })
}

if (!db.getRole("purgerRole")) {
    db.createRole({
      role: "purgerRole",
      privileges: [
        {
          resource: { db: "${DB_NAME}", collection: "" },
          actions: [
            "changeOwnCustomData",
            "changeOwnPassword",
            "remove",
            "find",
            "listCollections"
          ]
        }
      ],
      roles: []
    }, { w: "majority" })
}

if (!db.getUser("apiUser")) {
    db.createUser({
      user: "apiUser",
      pwd: "${APP_PWD}",
      roles: [ { role: "appRole", db: "${DB_NAME}" } ],
      mechanisms: ["SCRAM-SHA-256" ]
    },{ w: "majority" })
}

if (!db.getUser("purger")) {
    db.createUser({
      user: "purger",
      pwd: "${PURGER_PWD}",
      roles: [ { role: "purgerRole", db: "${DB_NAME}" } ],
      mechanisms: [ "SCRAM-SHA-256" ]
    }, { w: "majority" })
}
exit()
EOF

echo 'Created user with roles: {apiUser: appRole, purger: purgerRole}'
      