
curl -v "localhost:8080" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.sendgrid.email.send" \
       -H "Ce-Source: dev.knative.samples/helloworldsource" \
       -H "Content-Type: application/json" \
       -d '{
  "event": {
    "guid": "ocid1.cloudguardproblem.oc1.iad.1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e",
    "name": "$eventName",
    "severity": "CRITICAL",
    "shortDescription": "Object Storage supports anonymous, unauthenticated access to a bucket. A public bucket that has read access enabled for anonymous users allows anyone to obtain object metadata, download bucket objects, and optionally list bucket contents.",
    "startTime": "2022-03-30T23:25:02Z",
    "status": "OPEN"
  },
  "provider": {
    "accountId": "ocid1.tenancy.oc1..aaaaaaaagfqbe4ujb2dgwxtp37gzinrxt6h6hfshjokfgfi5nzquxmfpzkyq"
  },
  "providerId": "1",
  "providerType": "CSP",
  "resource": {
    "identifier": "orasenatdpltsecitom01/AutoVinci",
    "name": "AutoVinci",
    "region": "us-phoenix-1",
    "type": "Bucket",
    "zone": "Comp-Name"
  },
  "source": {
    "sourceId": "none",
    "sourceName": "Cloud Guard"
  }
}'
