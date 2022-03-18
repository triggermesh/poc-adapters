# TechulusPushTarget
the TechulusPushTarget kind exposes a service that allows the user to emit push notifications via [Techulus Push](https://push.techulus.com).

## Deploying with Koby

### Prerequisites
* Ensure that you have installed [Koby](https://github.com/triggermesh/koby) on the target cluster.

### Configuring the TechulusPushTarget CRD with Koby
The TechulusPushTarget CRD can be configured with [Koby](https://github.com/triggermesh/koby) by applying the provided manifest in `/config/100-registration.yaml`
```cmd
kubectl apply -f /config/100-registration.yaml
```

### Deploying an instance of the TechulusPushTarget
After updating the `api_key` field with a valid API key, a TechulusPushTarget can be deployed by applying the provided manifest in `/config/200-deployment.yaml`.
```cmd
kubectl apply -f /config/200-deployment.yaml
```

### Interacting with the TechulusPushTarget
The TechulusPushTarget object expects a cloudevent with a certin json payload. This payload can be defined as:
```go
type eventPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
```
An example curl request to be sent to the TechulusPushTarget service:
```
curl -v "localhost:8080" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.sample.event" \
       -H "Ce-Source: dev.knative.samples/helloworldsource" \
       -H "Content-Type: application/json" \
       -d '{"title":"Hello World","body":"Hello, world!"}'
```

Expecting a response like this
```
< HTTP/1.1 200 OK
< Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79
< Ce-Source: dev.knative.samples/helloworldsource
< Ce-Specversion: 1.0
< Ce-Time: 2022-03-01T14:13:58.683788Z
< Ce-Type: io.triggermesh.sample.event.response
< Content-Length: 82
< Content-Type: application/json
< Date: Tue, 01 Mar 2022 14:13:58 GMT
<
* Connection #0 to host localhost left intact
{"success":true,"responses":[{"success":true,"message":"Message send to device"}]}%
```
