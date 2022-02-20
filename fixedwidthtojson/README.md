# FixedWidthToJSONTransformation
The FixedWidthToJSONTransformation object expects a fixed width file sent via content-type `application/text`, transforms it into a JSON object and emits the event to the provided sink.


## Deploying with Koby

### Prerequisites
* Ensure that you have installed [Koby](https://github.com/triggermesh/koby) on the target cluster.

### Configuring the FixedWidthToJSONTransformation CRD with Koby
The FixedWidthToJSONTransformation CRD can be configured with [Koby](https://github.com/triggermesh/koby) by applying the provided manifest in `/config/100-registration.yaml`
```cmd
kubectl apply -f /config/100-registration.yaml
```

### Deploying an instance of the FixedWidthToJSONTransformation
The FixedWidthToJSONTransformation can now be deployed by applying the provided manifest in `/config/200-deployment.yaml`. Note that the FixedWidthToJSONTransformation must be deployed with a configured sink. For example purposes
the example deployment comes pre-configured with an event-display.
```cmd
kubectl apply -f /config/200-deployment.yaml
```
*note* The `dw_spell` field is required and must be a valid Dataweave spell. The `output_content_type` and `incoming_content_type` fields are optional, defaulted to `application/json`, and can be used to specify the content-type of the cloudevent response and the expected content-type of incoming cloudevents respectively.

### Interacting with the FixedWidthToJSONTransformation

send this
```
curl --location --request POST 'http://fixedwidthtojsontransformations-fwtojson.default.tmkongdemo.triggermesh.io' \
--header 'Ce-Specversion: 1.0' \
--header 'Ce-Type: io.triggermesh.sample.event' \
--header 'Ce-Id: 123123' \
--header 'Ce-Source: ser' \
--header 'Content-Type: text/plain' \
--data-raw 'NAME                STATE     TELEPHONE

John Smith          WA        418-Y11-4111

Mary Hartford       CA        319-Z19-4341

Evan Nolan          IL        219-532-c301
'
```
Expecting a response like this inside the logs of the event-display:
```
☁️  cloudevents.Event
Context Attributes,
  specversion: 1.0
  type: io.triggermesh.sample.event
  source: ser
  id: 123123
  time: 2022-02-03T19:25:00.094490019Z
  datacontenttype: application/json
Data,
{
  "fields": [
    {
      "value": "NAME",
      "spaceLeft": 0,
      "lineNumber": 0
    },
    {
      "value": "STATE",
      "spaceLeft": 16,
      "lineNumber": 0
    },
    {
      "value": " TELEPHONE",
      "spaceLeft": 4,
      "lineNumber": 0
    },
    {
      "value": "John Smith",
      "spaceLeft": 0,
      "lineNumber": 2
    },
    {
      "value": "WA",
      "spaceLeft": 10,
      "lineNumber": 2
    },
    {
      "value": "418-Y11-4111",
      "spaceLeft": 8,
      "lineNumber": 2
    },
    {
      "value": "Mary Hartford",
      "spaceLeft": 0,
      "lineNumber": 4
    },
    {
      "value": " CA",
      "spaceLeft": 6,
      "lineNumber": 4
    },
    {
      "value": "319-Z19-4341",
      "spaceLeft": 8,
      "lineNumber": 4
    },
    {
      "value": "Evan Nolan",
      "spaceLeft": 0,
      "lineNumber": 6
    },
    {
      "value": "IL",
      "spaceLeft": 10,
      "lineNumber": 6
    },
    {
      "value": "219-532-c301",
      "spaceLeft": 8,
      "lineNumber": 6
    }
  ]
}
```
