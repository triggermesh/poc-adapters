curl -v  "localhost:8080" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.csnf.aquasec.transformation" \
       -H "Ce-Source: projects/jeff-dev-env/topics/triggermeshcsnf" \
       -H "Content-Type: application/json" \
       -d '{
  "event": {
    "guid": "5f38a4a3-8047-4b63-adf5-5608f2a9f6eb",
    "name": "Container.Engine",
    "severity": "HIGH",
    "shortDescription": "an aquasec short description goes here",
    "startTime": 1627587060,
    "status": "OPEN"
  },
  "image": "httpd:latest",
  "provider": {
    "accountId": "aks-default-15484652-vmss000001.tpe5bzjk4yoevknn00ux3h31kb.ax.internal.cloudapp.net",
    "name": "Google PubSub",
    "providerId": "aks-default-15484652-vmss000001.tpe5bzjk4yoevknn00ux3h31kb.ax.internal.cloudapp.net",
    "providerType": "aquasec"
  },
  "resource": {
    "identifier": "aks-default-15484652-vmss000001.tpe5bzjk4yoevknn00ux3h31kb.ax.internal.cloudapp.net",
    "name": "5f38a4a3-8047-4b63-adf5-5608f2a9f6eb",
    "region": "us-1-c",
    "type": "Container",
    "zone": "us"
  },
  "source": {
    "sourceId": "5a4da19ff2703ad2f48db4d55e563a37828dccc0f11fb6c00e60a271ab3c37cb",
    "sourceName": "Aquasec"
  }
}'


curl -v  "localhost:8080" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.csnf.azure.defender.transformation" \
       -H "Ce-Source: projects/jeff-dev-env/topics/triggermeshcsnf" \
       -H "Content-Type: application/json" \
       -d '{
  "event": {
    "guid": "/subscriptions/97e01fd4-3326-41f4-b9e3-3cfd6809e10f/resourceGroups/Sample-RG/providers/Microsoft.Security/locations/centralus/alerts/2517538088722812591_ee939333-c75e-461d-83c2-712ba9abfadb",
    "name": "SIMULATED_SQL.VM_BruteForce",
    "severity": "High",
    "shortDescription": "THIS IS A SAMPLE ALERT: A successful login occurred after an apparent brute force attack on your resource",
    "startTime": "2022-03-28T18:26:47.490371Z",
    "status": "Active"
  },
  "provider": {
    "accountId": "97e01fd4-3326-41f4-b9e3-3cfd6809e10f"
  },
  "providerId": "Sample-VM",
  "providerType": "Microsoft.Security/Locations/alerts",
  "resource": {
    "identifier": "ee939333-c75e-461d-83c2-712ba9abfadb",
    "name": "Sample-VM",
    "region": "us-east-1",
    "type": "Microsoft.Security/Locations/alerts",
    "zone": "us"
  },
  "source": {
    "sourceId": "Microsoft",
    "sourceName": "Azure Defender"
  }
}'



curl -v  "localhost:8080" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.csnf.cloudguard.transformation" \
       -H "Ce-Source: projects/jeff-dev-env/topics/triggermeshcsnf" \
       -H "Content-Type: application/json" \
       -d '{
  "event": {
    "guid": "ocid1.cloudguardproblem.oc1.iad.1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e",
    "name": "SCANNED_HOST_VULNERABILITY",
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
    "identifier": "ocid1.instance.oc1.iad.1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e1q2w3e",
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
