# KafkaSource
A Knative Source that reads from a list of Kafka topics.

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
export CONFLUENT_BOOTSTRAP_SERVERS=pkc-419q3.us-east4.gcp.confluent.cloud:9092
export CONFLUENT_SASL_USERNAME=NX5F6PNUER4GLXFE
export CONFLUENT_SASL_PASSWORD=yndA/y2s9n3fBjjSju9OJaC09pPfIODoQGMUXJWzyUIpjd5dedIXSj6XYCN8wA3+
export CONFLUENT_TOPIC=triggermesh
