```
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
export K_SINK=http://localhost:8081
```


```
curl -v "http://gcloudgard-transformation.default.tmkongdemo.triggermesh.io" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: as" \
       -H "Ce-Source: dev.knative.samples/helloworldsource" \
       -H "Content-Type: application/json" \
       -d '{
    "id": 442,
    "time": 1627587060,
    "date": 0,
    "type": "Container.Engine",
    "user": "",
    "action": "exec",
    "image": "httpd:latest",
    "imagehash": "sha256:73b8cfec11558fe86f565b4357f6d6c8560f4c49a5f15ae970a24da86c9adc93",
    "imageid": "",
    "container": "apache",
    "containerid": "5a4da19ff2703ad2f48db4d55e563a37828dccc0f11fb6c00e60a271ab3c37cb",
    "host": "aks-default-15484652-vmss000001.tpe5bzjk4yoevknn00ux3h31kb.ax.internal.cloudapp.net",
    "hostid": "5f38a4a3-8047-4b63-adf5-5608f2a9f6eb",
    "category": "container",
    "result": 2,
    "data": "{\"host\": \"aks-default-15484652-vmss000001.tpe5bzjk4yoevknn00ux3h31kb.ax.internal.cloudapp.net\", \"rule\": \"test-block-exec\", \"time\": 1627587060, \"image\": \"httpd:latest\", \"level\": \"block\", \"vm_id\": \"909679a6-f7a8-4d1e-ab49-ebce7eaef47d\", \"action\": \"exec\", \"hostid\": \"5f38a4a3-8047-4b63-adf5-5608f2a9f6eb\", \"hostip\": \"10.240.0.5\", \"reason\": \"Unauthorized container exec\", \"result\": 2, \"tactic\": \"Execution\", \"control\": \"Block Container Exec\", \"imageid\": \"73b8cfec11558fe86f565b4357f6d6c8560f4c49a5f15ae970a24da86c9adc93\", \"podname\": \"apache\", \"podtype\": \"container\", \"vm_name\": \"aks-default-15484652-vmss_1\", \"category\": \"container\", \"resource\": \"bash\", \"vm_group\": \"MC_rnd-aks2729-aks-rg_aks2729_westeurope\", \"container\": \"apache\", \"hostgroup\": \"aquactl-default-enforcer-group\", \"rule_type\": \"runtime.policy\", \"technique\": \"Command and Script Interpreter\", \"repository\": \"httpd\", \"containerid\": \"5a4da19ff2703ad2f48db4d55e563a37828dccc0f11fb6c00e60a271ab3c37cb\", \"k8s_cluster\": \"aqua-secure\", \"vm_location\": \"westeurope\", \"podnamespace\": \"test\"}",
    "account_id": 0
}'
```

expect

```
{"provider": {"providerId": "1", "providerType": "CSP", "name": "Oracle Cloud Infrastructure", "accountId": "ocid1.tenancy.oc1..aaaaaaaagfqbe4ujb2dgwxtp37gzinrxt6h6hfshjokfgfi5nzquxmfpzkyq"}, "source": {"sourceName": "Cloud Guard", "sourceId": "None"}, "event": {"guid": "ocid1.cloudguardproblem.oc1.iad.amaaaaaa24o7ld2qphw36wghvk44yms2u3hm4wnzhvtgcakktusurhtepevq", "name": "Public Bucket", "shortDescription": "Object Storage supports anonymous, unauthenticated access to a bucket. A public bucket that has read access enabled for anonymous users allows anyone to obtain object metadata, download bucket objects, and optionally list bucket contents.", "startTime": "2021-08-28T16:37:36.945Z", "severity": "CRITICAL", "status": "OPEN", "recommendation": "Ensure that the bucket is sanctioned for public access, and if not, direct the OCI administrator to restrict the bucket policy to allow only specific users access to the resources required to accomplish their job functions."}, "resource": {"identifier": "orasenatdpltsecitom01/AutoVinci", "type": "Bucket", "name": "AutoVinci", "region": "us-phoenix-1", "zone": "McMillan"}}
```