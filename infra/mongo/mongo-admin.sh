#!/bin/bash
set -e

# Configurations -- UPDATE_ME
MONGO_HOST="localhost"
MONGO_PORT="27017"
MONGO_ADMIN_USER="" 
MONGO_ADMIN_PASS=''
APP_PWD=''
PURGER_PWD=''
MIGRATOR_PWD=''

if [ -z "$APP_PWD" ] || [ -z "$PURGER_PWD" ] || [ -z "$MIGRATOR_PWD" ]; then
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
    })
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
    })
}

if (!db.getRole("migratorRole")) {
    db.createRole({
      role: "migratorRole",
      privileges: [
        {
          resource: { db: "${DB_NAME}", collection: "" },
          actions: [
            "bypassDocumentValidation",
            "changeOwnCustomData",
            "changeOwnPassword",
            "collMod",
            "convertToCapped",
            "createCollection",
            "createIndex",
            "createSearchIndexes",
            "dropCollection",
            "dropIndex",
            "dropSearchIndex",
            "listCollections",
            "listIndexes",
            "listSearchIndexes",
            "reIndex",
            "updateSearchIndex",
            "renameCollectionSameDB",
            "validate",
            "storageDetails",
            "update"
          ]
        }
      ],
      roles: []
    })
}

if (!db.getUser("apiUser")) {
    db.createUser({
      user: "apiUser",
      pwd: "${APP_PWD}",
      roles: [ { role: "appRole", db: "${DB_NAME}" } ],
      mechanisms: ["SCRAM-SHA-256"]
    })
}

if (!db.getUser("purger")) {
    db.createUser({
      user: "purger",
      pwd: "${PURGER_PWD}",
      roles: [ { role: "purgerRole", db: "${DB_NAME}" } ],
      mechanisms: [ "SCRAM-SHA-256"]
    })
}

if (!db.getUser("migrator")) {
    db.createUser({
      user: "migrator",
      pwd: "${MIGRATOR_PWD}",
      roles: [ { role: "migratorRole", db: "${DB_NAME}" } ],
      mechanisms: [ "SCRAM-SHA-256"]
    })
}
exit()
EOF

echo 'Created user with roles: {apiUser: appRole, purger: purgerRole, migrator: migratorRole}'
      