# MongoDBTarget

The MongoDBTarget exposes several methods, via event types, that can be used to interact with the MongoDB database.
The MongoDBTarget can currently be deployed via [Koby](https://github.com/triggermesh/koby)

## Deploying with Koby

### Prerequisites
* Ensure that you have installed [Koby](https://github.com/triggermesh/koby) on the target cluster.
* A pre-existing MongoDB database with its associated connection string.

### Configuring the MongoDBTarget CRD with Koby
The MongoDBTarget CRD can be configured with Koby by applying the provided manifest in `/config/100-registration.yaml`
```cmd
kubectl apply -f /config/100-registration.yaml
```

### Deploying an instance of the MongoDBTarget
After providing a valid connection string for the MongoDB database under the `mongodb_server_url` spec field,
the MongoDBTarget can now be deployed by applying the provided manifest in `/config/200-deployment.yaml`.
```cmd
kubectl apply -f /config/200-deployment.yaml
```

# Interacting with the Event Target

## Arbitrary Event Types
The mongoDBTarget supports accepting arbitrary event types. These events will be inserted into the MongoDB database/collection specified in the `defaultDatabase` and `defaultCollection` spec fields.

So if one were to send the following event:

```cmd
curl -v  http://localhost:8080 \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.sample.event" \
       -H "Ce-Source: sample" \
       -H "Content-Type: application/json" \
       -d '{"example":"event"}'

```

We would expect this entry in the MongoDB database/collection:
```
{"_id":{"$oid":"62029b4e20372fe01225194d"},"example":"event"}
```

## Pre-defined Event Types
Use these event types and associated payloads to interact with the MongoDB Target.

### io.triggermesh.mongodb.insert

Events of this type intend to post a single key:value pair to MongoDB

#### Example CE posting an event of type "io.triggermesh.mongodb.insert"


```cmd
curl -v  http://localhost:8080 \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.mongodb.insert" \
       -H "Ce-Source: sample/source" \
       -H "Content-Type: application/json" \
       -d '{"database":"test","collection": "test","mapStrVal":{"test":"testdd1","test2":"test3"}}'
```


#### This type expects a JSON payload with the following properties:

| Name  |  Type |  Comment |
|---|---|---|
| **database** | string | The name of the database.  |
| **collection** | string | The value of the collection. |
| **itemName** | string | This value will be used to assing a Key.  |
| **mapStrVal** | map[string]string | This value will be used to assing a Value. |

**Note** the `database` and `collection` fields are not required. If not provided, the `defaultDatabase` and `defaultCollection` spec fields will be used.

### io.triggermesh.mongodb.update

Events of this type intend to update a single pre-existing key:value pair

#### Example CE posting an event of type "io.triggermesh.mongodb.update"

```cmd
curl -v http://localhost:8080 \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.mongodb.update" \
       -H "Ce-Source: sample/source" \
       -H "Content-Type: application/json" \
       -d '{"database":"test","collection": "test","searchKey":"test","searchValue":"testdd1","updateKey":"partstore","updateValue":"UP FOR GRABS"}'
```

#### This type expects a JSON payload with the following properties:

| Name  |  Type |  Comment |
|---|---|---|
| **database** | string | The name of the database.  |
| **collection** | string | The value of the collection. |
| **itemName** | string | This value will be used to assing a Key.  |
| **searchKey** | string | . |
| **searchValue** | string | .  |
| **updateKey** | string | .  |
| **updateValue** | string |. |

**Note** the `database` and `collection` fields are not required. If not provided, the `defaultDatabase` and `defaultCollection` spec fields will be used.

### io.triggermesh.mongodb.query.kv

Events of this type intend to query a MongoDB

#### Example CE posting an event of type "io.triggermesh.mongodb.query.kv"

```cmd
curl -v http://localhost:8080 \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.mongodb.query.kv" \
       -H "Ce-Source: sample/source" \
       -H "Content-Type: application/json" \
       -d '{"database":"test","Collection": "test","key":"partstore","value":"UP FOR GRABS"}'
```

#### This type expects a JSON payload with the following properties:

| Name  |  Type |  Comment |
|---|---|---|
| **database** | string | The name of the database.  |
| **collection** | string | The value of the collection. |
| **key** | string | the "Key" value to search  |
| **value** | string | the "Value" value to search |

**Note** the `database` and `collection` fields are not required. If not provided, the `defaultDatabase` and `defaultCollection` spec fields will be used.


# Local Development

To build and run this Target locally, run the following command(s):

```cmd
export MONGODB_SERVER_URL=mongodb+srv://<user>:<password>@<database_url>/myFirstDatabase
export MONGODB_DEFAULT_DATABASE=test
export MONGODB_DEFAULT_COLLECTION=testcol
export DISCARD_CE_CONTEXT=true
go run cmd/main.go
```
