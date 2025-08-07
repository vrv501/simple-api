
// connect to mongodb using admin creds then switch to pets server 
// mongosh -u purger -p myPassword  --authenticationDatabase shop shop

use shop
db.createRole({
  role: "dbPowerUser",
  privileges: [
    {
      resource: { db: 'shop', collection: '' },
      actions: [
        'analyze',
        'bypassDocumentValidation',
        'changeOwnCustomData',
        'changeOwnPassword',
        'changeStream',
        'collMod',
        'collStats',
        'compact',
        'configureQueryAnalyzer',
        'convertToCapped',
        'createCollection',
        'createIndex',
        'createSearchIndexes',
        'dbHash',
        'dbStats',
        'dropIndex',
        'dropSearchIndex',
        'enableProfiler',
        'find',
        'insert',
        'killCursors',
        'listCollections',
        'listIndexes',
        'listSearchIndexes',
        'updateSearchIndex',
        'planCacheIndexFilter',
        'planCacheRead',
        'planCacheWrite',
        'reIndex',
        'renameCollectionSameDB',
        'storageDetails',
        'update',
        'validate',
        'viewRole',
        'viewUser'
      ]
    }
  ],
  roles: []
})

db.createRole({
  role: "dataPurger",
  privileges: [
    {
      resource: { db: 'shop', collection: '' },
      actions: [
        'cleanupStructuredEncryptionData',
        'compactStructuredEncryptionData',
        'changeOwnCustomData',
        'changeOwnPassword',
        'remove',
      ]
    }
  ],
  roles: []
})

db.createUser({
  user: "api-server",
  pwd: "myPassword",
  roles: [ { role: "dbPowerUser", db: "shop" } ],
  mechanisms: ["SCRAM-SHA-256"]
})

db.createUser({
  user: "purger",
  pwd: "myPassword",
  roles: [ { role: "dataPurger", db: "shop" } ],
  mechanisms: [ "SCRAM-SHA-256"]
})

db.createCollection("animal-categories", {
   validator: {
      $jsonSchema: {
         bsonType: "object",
         required: [ "name", "created_on" ],
         properties: {
            name: {
               bsonType: "string",
               minLength: 1,
               maxLength: 150,
               description: "animal category"
            },
            created_on: {
               bsonType: "date",
               description: "creation date time(UTC) of document"
            },
            updated_on: {
               bsonType: "date",
               description: "last updated date time(UTC) of document"
            },
         }
      }
   }
})

db.getCollection("animal-categories").createIndex({})
