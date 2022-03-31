
curl -v "localhost:8080" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.sendgrid.email.send" \
       -H "Ce-Source: dev.knative.samples/helloworldsource" \
       -H "Content-Type: application/json" \
       -d '{
  "event": {
    "guid": "apache",
    "name": "apache",
    "severity": "HIGH",
    "shortDescription": "an aquasec short description goes here",
    "startTime": 1627587060,
    "status": "OPEN"
  },
  "image": "ubuntu:latest",
  "provider": {
    "accountId": 442,
    "name": "SQS",
    "providerId": 442,
    "providerType": "aquasec"
  },
  "resource": {
    "identifier": 442,
    "name": "5f38a4a3-8047-4b63-adf5-5608f2a9f6eb",
    "region": "us-1-c",
    "type": "something",
    "zone": "us"
  },
  "source": {
    "sourceId": "none",
    "sourceName": "Aquasec"
  }
}'
