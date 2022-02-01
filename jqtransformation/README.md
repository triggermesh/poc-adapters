# JQTransformation
the JQTransformation exposes a service that allows the user to transform a JSON payload
by using a [JQ](https://devdocs.io/jq/) expression.

## Deploying with Koby

### Prerequisites
* Ensure that you have installed [Koby](https://github.com/triggermesh/koby) on the target cluster.

### Configuring the JQTransformation CRD with Koby
The JQTransformation CRD can be configured with [Koby](https://github.com/triggermesh/koby) by applying the provided manifest in `/config/100-registration.yaml`
```cmd
kubectl apply -f /config/100-registration.yaml
```

### Deploying an instance of the JQTransformation
After updating the `query` spec field with a valid JQ expression, the JQTransformation can now be deployed by applying
 the provided manifest in `/config/200-deployment.yaml`.
```cmd
kubectl apply -f /config/200-deployment.yaml
```

### Interacting with the JQTransformation
The JQTransformation object will accept any event with a valid JSON payload and transform it using the JQ expression provided in the spec.
If it was deployed with the example JQ expression, one can try the following example event:
```
curl -v "http://jqtransformations-hello-jq.default.34.133.226.173.sslip.io" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.transform.me" \
       -H "Ce-Source: dev.knative.samples/helloworldsource" \
       -H "Content-Type: application/json" \
       -d '{"foo":"richard@triggermesh.com"}'
```
and expect a response of:
```
"richard@triggermesh.com"
```
