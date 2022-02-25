# KafkaSource
the KafkaSource exposes a service that allows the user to transform a JSON payload
by using a [JQ](https://devdocs.io/jq/) expression.

## Deploying with Koby

### Prerequisites
* Ensure that you have installed [Koby](https://github.com/triggermesh/koby) on the target cluster.

### Configuring the KafkaSource CRD with Koby
The KafkaSource CRD can be configured with [Koby](https://github.com/triggermesh/koby) by applying the provided manifest in `/config/100-registration.yaml`
```cmd
kubectl apply -f /config/100-registration.yaml
```

### Deploying an instance of the KafkaSource
After updating the manifest with valid information, the KafkaSource can now be deployed by applying the provided manifest in `/config/200-deployment.yaml`.

```cmd
kubectl apply -f /config/200-deployment.yaml
```
