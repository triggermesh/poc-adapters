```
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

```
curl -v "http://copybooktransformations-cb.default.tmkongdemo.triggermesh.io" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.sendgrid.email.send" \
       -H "Ce-Source: dev.knative.samples/helloworldsource" \
       -H "Content-Type: application/text" \
       -d "0200400400000000900011111971-01-21JOHN ROBERT PERIN                       +0001001-00100.101234  \n"
```
