```
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
export K_SINK=http://localhost:8081
```


```
curl -v "http://localhost:8080" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.sendgrid.email.send" \
       -H "Ce-Source: dev.knative.samples/helloworldsource" \
       -H "Content-Type: application/json" \
       -d '{
    "eventType" : "com.oraclecloud.cloudguard.problemdetected",
    "cloudEventsVersion" : "0.1",
    "eventTypeVersion" : "2.0",
    "source" : "CloudGuardResponderEngine",
    "eventTime" : "2021-08-28T16:37:59Z",
    "contentType" : "application/json",
    "data" : {
      "compartmentId" : "ocid1.compartment.oc1..aaaaaaaaeqm7cycpkedmlqzmitgjy5ya5urtae6zoogb7epspwrqz2kxxwta",
      "compartmentName" : "McMillan",
      "resourceName" : "Bucket is public",
      "resourceId" : "ocid1.cloudguardproblem.oc1.iad.amaaaaaa24o7ld2qphw36wghvk44yms2u3hm4wnzhvtgcakktusurhtepevq",
      "additionalDetails" : {
        "tenantId" : "ocid1.tenancy.oc1..aaaaaaaagfqbe4ujb2dgwxtp37gzinrxt6h6hfshjokfgfi5nzquxmfpzkyq",
        "status" : "OPEN",
        "reason" : "New Problem detected by CloudGuard",
        "problemName" : "BUCKET_IS_PUBLIC",
        "riskLevel" : "CRITICAL",
        "problemType" : "CONFIG_CHANGE",
        "resourceName" : "AutoVinci",
        "resourceId" : "orasenatdpltsecitom01/AutoVinci",
        "resourceType" : "Bucket",
        "targetId" : "ocid1.cloudguardtarget.oc1.iad.amaaaaaa24o7ldyaborosittsemtunkb2kkn54bysbznwa4xdu7rj3mpulza",
        "labels" : "CIS_OCI_V1.1_OBJECTSTORAGE, ObjectStorage",
        "firstDetected" : "2021-08-28T16:37:36.945Z",
        "lastDetected" : "2021-08-28T16:37:36.945Z",
        "region" : "us-phoenix-1",
        "problemDescription" : "Object Storage supports anonymous, unauthenticated access to a bucket. A public bucket that has read access enabled for anonymous users allows anyone to obtain object metadata, download bucket objects, and optionally list bucket contents.",
        "problemRecommendation" : "Ensure that the bucket is sanctioned for public access, and if not, direct the OCI administrator to restrict the bucket policy to allow only specific users access to the resources required to accomplish their job functions."
      }
    },
    "eventID" : "1642b454-725e-412f-bfc9-783803f83134",
    "extensions" : {
      "compartmentId" : "ocid1.compartment.oc1..aaaaaaaaeqm7cycpkedmlqzmitgjy5ya5urtae6zoogb7epspwrqz2kxxwta"
    }
   }'
```
