kubectl create secret generic cxnpl-redpanda-kafka-poc \
--from-literal=protocol=SASL_SSL \
--from-literal=sasl.mechanism=SCRAM-SHA-256 \
--from-file=ca.crt=certs/ca.crt \
--from-file=certsclient.crt=certs/client.crt \
--from-file=client.key=certs/client.key \
--from-literal=user=cxnpl-poc-sa \
--from-literal=password="aYupifBYaT;CP_loz-6dt-"

kubectl create secret --namespace cn generic kafkacerts \
  --from-literal=protocol=SSL \
  --from-file=ca.crt=certs/ca.crt \
  --from-file=user.crt=certs/client.crt \
  --from-file=user.key=certs/client.key
