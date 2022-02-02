# JSONToXMLTransformation
the JSONToXMLTransformation kind exposes a service that allows the user to transform a JSON payload
into an XML payload.

***Notes on the conversion***
An object is wrapped in <object></object>
An array is wrapped in <array></array>
A boolean is wrapped in <boolean></boolean> , with either "true" or "false" as chardata
A number is wrapped in <number></number>
A string is wrapped in <string></string>
A null becomes <null></null> , with no chardata

## Deploying with Koby

### Prerequisites
* Ensure that you have installed [Koby](https://github.com/triggermesh/koby) on the target cluster.

### Configuring the JSONToXMLTransformation CRD with Koby
The JSONToXMLTransformation CRD can be configured with [Koby](https://github.com/triggermesh/koby) by applying the provided manifest in `/config/100-registration.yaml`
```cmd
kubectl apply -f /config/100-registration.yaml
```

### Deploying an instance of the JSONToXMLTransformation
The JSONToXMLTransformation must be deployed alongside a sink. For example purposes the example deployment comes pre-configured with an event-display. It can be deployed by applying the provided manifest in `/config/200-deployment.yaml`.
```cmd
kubectl apply -f /config/200-deployment.yaml
```

### Interacting with the JSONToXMLTransformation
The JSONToXMLTransformation object will accept any event with a valid JSON payload.
One can try the following example event:
```
curl -v "http://jsontoxmltransformations-hello-jtx.j2x.35.202.146.138.sslip.io" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.transform.me" \
       -H "Ce-Source: dev.knative.samples/helloworldsource" \
       -H "Content-Type: application/json" \
       -d '{"foo":"richard@triggermesh.com", "numVal":12}'
```

Expecting a response like this inside the logs of the event-display:
```
☁️  cloudevents.Event
Context Attributes,
  specversion: 1.0
  type: io.triggermesh.transform.me
  source: dev.knative.samples/helloworldsource
  id: 536808d3-88be-4077-9d7a-a3f162705f79
  time: 2022-02-02T17:51:12.717774863Z
  datacontenttype: application/xml
Data,
  <object><string name="foo">richard@triggermesh.com</string><number name="numVal">12</number></object>
```
