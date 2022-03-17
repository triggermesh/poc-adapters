export AZURE_CLIENT_SECRET=Fll7Q~fQ4uu_EkmvYAorlEO196CDJ6osTlC1C
export AZURE_TENANT_ID=f14eddee-e73b-481d-8237-17983764afcb
export AZURE_CLIENT_ID=6fbbd6c1-a890-49ca-af2a-142bade07e7a
export AZURE_SENTINEL_SUBSCRIPTION_ID=77641a71-ffc3-4cfd-abd9-6ff8dc509a3d
export AZURE_SENTINEL_RESOURCE_GROUP=sent
export AZURE_SENTINEL_WORKSPACE=sent

curl -v "localhost:8080" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.sendgrid.email.send" \
       -H "Ce-Source: dev.knative.samples/helloworldsource" \
       -H "Content-Type: application/json" \
       -d '{
  "event": {
    "event": {
      "metadata": {
        "guid": 442,
        "name": "Block Container Tom",
        "url": null,
        "severity": "high",
        "shortDescription": "Unauthorized container exec",
        "longDescription": "Unauthorized container exec",
        "time": 1627587061
      },
      "producer": {
        "name": "aquasec"
      },
      "reporter": {
        "name": "Aqua Security"
      },
      "resources": [
        {
          "guid": "5a4da19ff2703ad2f48db4d55e563a37828dccc0f11fb6c00e60a271ab3c37cb",
          "name": "apache",
          "region": "westeurope",
          "platform": "Azure",
          "service": "Azure Kubernetes Service",
          "type": "container",
          "accountId": null,
          "package": "httpd:latest"
        }
      ]
    },
    "decoration": [
      {
        "decorator": "dockerhub",
        "timestamp": "2022-03-16T07:48:00.251Z",
        "payload": {
          "registry": "docker.io",
          "namespace": "library",
          "image": "httpd",
          "tag": "latest",
          "digests": [
            "sha256:3d88724ab878b75b231d5de35f3ea2e9368fa60718efd399767f43699c68ba85",
            "sha256:3fab949dc330282fea1238327bc2d6b932b3470e1a858aa07b46919a3d268fe0",
            "sha256:fcdcb30653211aaedaaae41855f037f37ba73565bdb88356250930d81ae666be",
            "sha256:d47883cea91afe3a030ead39598c1758d6e969a18fb27ca097d56a603d4779f3",
            "sha256:5f3ade90bd28792f4d2b39ee0cfb9947205311bd50f03c1e3c5a57ea45fc3e77",
            "sha256:82cc6bd1e1f4b29240060a3a5470999d66ba9b9b4e8389c1e74cafcd37213bd3",
            "sha256:98be26a87a6e7412c213a625e8e0910815e71b47a432621e48306c474df664ca",
            "sha256:a1e91eb1036b5ece161a774a269967b7ea2663b3c2688972dc1277c434a093f0"
          ],
          "imageLastUpdated": "2022-03-14T21:36:39.508026Z",
          "tagLastUpdated": "2022-03-14T21:36:38.99727Z",
          "description": "The Apache HTTP Server Project",
          "starCount": 3917,
          "pullCount": 3609541407
        }
      }
    ]
  },
  "sourcetype": "Aqua Security"
}'
