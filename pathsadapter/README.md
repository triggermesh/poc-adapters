export PATH_A_CONTINUE_IF='(event.fromEmail == "richard@triggermesh.com")'

export PATH_A_CONTINUE_PATH=http://tmdebugger.default.tmkongdemo.triggermesh.io

export PATH_A_CONTINUE_TYPE=io.triggermesh.paths.ContinuePath.a

export PATH_B_CONTINUE_IF='(event.fromEmail == "bob@triggermesh.com")'

export PATH_B_CONTINUE_PATH=http://tmdebugger.default.tmkongdemo.triggermesh.io

export PATH_B_CONTINUE_TYPE=io.triggermesh.paths.ContinuePath.b

export DEFAULT_CONTINUE_PATH=http://tmdebugger.default.tmkongdemo.triggermesh.io

export DEFAULT_CONTINUE_TYPE=io.triggermesh.paths.ContinuePath.default

path a
```
curl -v "http://paths-hello-path.default.tmkongdemo.triggermesh.io" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.sendgrid.email.send" \
       -H "Ce-Source: dev.knative.samples/helloworldsource" \
       -H "Content-Type: application/json" \
       -d '{"fromEmail":"richard@triggermesh.com"}'
```
path b
```
curl -v "http://paths-hello-path.default.tmkongdemo.triggermesh.io" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.sendgrid.email.send" \
       -H "Ce-Source: dev.knative.samples/helloworldsource" \
       -H "Content-Type: application/json" \
       -d '{"fromEmail":"bob@triggermesh.com"}'
```
default
```
curl -v "http://paths-hello-path.default.tmkongdemo.triggermesh.io" \
       -X POST \
       -H "Ce-Id: 536808d3-88be-4077-9d7a-a3f162705f79" \
       -H "Ce-Specversion: 1.0" \
       -H "Ce-Type: io.triggermesh.sendgrid.email.send" \
       -H "Ce-Source: dev.knative.samples/helloworldsource" \
       -H "Content-Type: application/json" \
       -d '{"fromEmail":"noone@triggermesh.com"}'
```

View results in TM debugger.
