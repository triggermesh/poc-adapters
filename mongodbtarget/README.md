
### io.triggermesh.mongodb.insert

Events of this type intend to post a single key:value pair to MongoDB

#### Example CE posting an event of type "io.triggermesh.mongodb.insert"


```cmd
curl -v  http://localhost:8080 \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.mongodb.insert" \
       -H "Ce-Source: ocimetrics/adapter" \
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

### io.triggermesh.mongodb.update

Events of this type intend to update a single pre-existing key:value pair

#### Example CE posting an event of type "io.triggermesh.mongodb.update"

```cmd
curl -v http://localhost:8080 \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.mongodb.update" \
       -H "Ce-Source: ocimetrics/adapter" \
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

### io.triggermesh.mongodb.query.kv

Events of this type intend to query a MongoDB

#### Example CE posting an event of type "io.triggermesh.mongodb.query.kv"

```cmd
curl -v http://localhost:8080 \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.mongodb.query.kv" \
       -H "Ce-Source: ocimetrics/adapter" \
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

### Export to run locally

export NAMESPACE=default

export K_LOGGING_CONFIG={}

export K_METRICS_CONFIG={}

export MONGODB_SERVER_URL="mongodb+srv://<user>:<password>@<database_url>/myFirstDatabase"
