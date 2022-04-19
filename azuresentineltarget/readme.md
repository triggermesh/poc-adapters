curl -v  "localhost:8080" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.sendgrid.email.send" \
       -H "Ce-Source: dev.knative.samples/helloworldsource" \
       -H "Content-Type: application/json" \
       -d'{
  "event": {
    "guid": "ocid1.cloudguardroblem.oc1.iad.1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e",
    "name": "sas",
    "severity": "CRITICAL",
    "shortDescription": "Prerequisite: Create a Host Scan Recipe and a Host Scan Target in the Scanning service. The Scanning service scans compute hosts to identify known cybersecurity vulnerabilities related to applications, libraries, operating systems, and services. This detector triggers a problem when the Scanning service has reported that an instance has one or more CRITICAL (or lower severity, based on the Input Settings within the detector config) vulnerabilities.",
    "startTime": "2022-01-30T11:14:29.130Z",
    "status": "OPEN"
  },
  "provider": {
    "accountId": "ocid1.tenancy.oc1..1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e"
  },
  "providerId": "1",
  "providerType": "CSP",
  "resource": {
    "identifier": "ocid1.instance.c1.iad.1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e",
    "name": "aDocker",
    "region": "us-ashburn-1",
    "type": "HostAgentScan",
    "zone": "Comp-Name"
  },
  "source": {
    "sourceId": "none",
    "sourceName": "Cloud Guard"
  }
}'
