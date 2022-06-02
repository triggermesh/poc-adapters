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


export CONFLUENT_BOOTSTRAP_SERVERS=0.rp-4260ba7.e539449.byoc.vectorized.cloud:30684
export CONFLUENT_TOPIC=test.topic
export SSL_CA_LOCATION=./ca.crt
export SSL_CLIENT_CERT=./client.crt
export SSL_CLIENT_KEY=./client.key
export USERNAME=cxnpl-poc-sa
export PASSWORD='aYupifBYaT;CP_loz-6dt-'
