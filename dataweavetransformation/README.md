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
After updating the `dw_spell`, `output_content_type`, and `incoming_content_type` spec fields. the DataweaveTransformation can now be deployed by applying the provided manifest in `/config/200-deployment.yaml`. Note that the DataweaveTransformation must be deployed with a configured sink. For example purposes
the example deployment comes pre-configured with an event-display.
```cmd
kubectl apply -f /config/200-deployment.yaml
```
*note* The `dw_spell` field is required and must be a valid Dataweave spell. The `output_content_type` and `incoming_content_type` fields are optional, defaulted to `application/json`, and can be used to specify the content-type of the cloudevent response and the expected content-type of incoming cloudevents respectively.

### Interacting with the DataweaveTransformation
The DataweaveTransformation object will accept any event and use the Dataweave spell provided in the spec to transform it.
If it was deployed with the example Dataweave spell, one can try the following example event:
```cmd
curl --location --request POST 'http://dataweavetransformations-hello-dw.dw.35.202.146.138.sslip.io' \
--header 'Ce-Specversion: 1.0' \
--header 'Ce-Type: io.triggermesh.mongodb.query.kv' \
--header 'Ce-Id: 123123' \
--header 'Ce-Source: ser' \
--header 'Content-Type: application/json' \
--data-raw '{
    "books": [
      {
        "-category": "cooking",
        "title":"Everyday Italian",
        "author": "Giada De Laurentiis",
        "year": "2005",
        "price": "30.00"
      },
      {
        "-category": "children",
        "title": "Harry Potter",
        "author": "J K. Rowling",
        "year": "2005",
        "price": "29.99"
      },
      {
        "-category": "web",
        "title":  "XQuery Kick Start",
        "author": [
          "James McGovern",
          "Per Bothner",
          "Kurt Cagle",
          "James Linn",
          "Vaidyanathan Nagarajan"
        ],
        "year": "2003",
        "price": "49.99"
      },
      {
        "-category": "web",
        "-cover": "paperback",
        "title": "Learning XML",
        "author": "Erik T. Ray",
        "year": "2003",
        "price": "39.95"
      }
    ]
}'
```
Expecting a response like this inside the logs of the event-display:
```
☁️  cloudevents.Event
Context Attributes,
  specversion: 1.0
  type: dev.knative.sources.ping
  source: /apis/v1/namespaces/dw/pingsources/cj
  id: 4c1a3770-cc4b-46b1-94b8-f50f9e918c50
  time: 2022-02-03T19:25:00.094490019Z
  datacontenttype: application/json
Data,
  {
    "items": [
      {
        "book": {
          "-CATEGORY": "cooking",
          "TITLE": "Everyday Italian",
          "AUTHOR": "Giada De Laurentiis",
          "YEAR": "2005",
          "PRICE": "30.00"
        }
      },
      {
        "book": {
          "-CATEGORY": "children",
          "TITLE": "Harry Potter",
          "AUTHOR": "J K. Rowling",
          "YEAR": "2005",
          "PRICE": "29.99"
        }
      },
      {
        "book": {
          "-CATEGORY": "web",
          "TITLE": "XQuery Kick Start",
          "AUTHOR": [
            "James McGovern",
            "Per Bothner",
            "Kurt Cagle",
            "James Linn",
            "Vaidyanathan Nagarajan"
          ],
          "YEAR": "2003",
          "PRICE": "49.99"
        }
      },
      {
        "book": {
          "-CATEGORY": "web",
          "-COVER": "paperback",
          "TITLE": "Learning XML",
          "AUTHOR": "Erik T. Ray",
          "YEAR": "2003",
          "PRICE": "39.95"
        }
      }
    ]
  }
```
