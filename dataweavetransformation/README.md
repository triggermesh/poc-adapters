# DataweaveTransformation
The DataweaveTransformation exposes a service that allows the user to transform a JSON payload
by using the [Dataweave Language](https://docs.mulesoft.com/mule-runtime/3.9/dataweave).

## Deploying with Koby

### Prerequisites
* Ensure that you have installed [Koby](https://github.com/triggermesh/koby) on the target cluster.

### Configuring the DataweaveTransformation CRD with Koby
The DataweaveTransformation CRD can be configured with [Koby](https://github.com/triggermesh/koby) by applying the provided manifest in `/config/100-registration.yaml`
```cmd
kubectl apply -f /config/100-registration.yaml
```

### Deploying an instance of the DataweaveTransformation
After updating the `dw_spell`, `output_content_type`, and `incoming_content_type` spec fields. the DataweaveTransformation can now be deployed by applying the provided manifest in `/config/200-deployment.yaml`.
```cmd
kubectl apply -f /config/200-deployment.yaml
```
*note* The `dw_spell` field is required and must be a valid Dataweave spell. The `output_content_type` and `incoming_content_type` fields are optional, defaulted to `application/json`, and can be used to specify the content-type of the cloudevent response and the expected content-type of incoming cloudevents respectively.

### Interacting with the DataweaveTransformation
The DataweaveTransformation object will accept any event and use the Dataweave spell provided in the spec to transform it.
If it was deployed with the example Dataweave spell, one can try the following example event:
```cmd
curl --location --request POST  'http://dataweavetransformations-hello-dw.default.34.133.226.173.sslip.io' \
--header 'Ce-Specversion: 1.0' \
--header 'Ce-Type: io.triggermesh.mongodb.query.kv' \
--header 'Ce-Id: 123123' \
--header 'Ce-Source: ser' \
--header 'Content-Type: application/json' \
--data-raw '[
  {
    "name": "User1",
    "age": 19
  },
  {
    "name": "User2",
    "age": 18
  },
  {
    "name": "User3",
    "age": 15
  },
  {
    "name": "User4",
    "age": 13
  },
  {
    "name": "User5",
    "age": 16
  }
]'
```
and expect a response of:
```
[
  {
    "name": "User1",
    "age": 19
  },
  {
    "name": "User2",
    "age": 18
  }
]
```
