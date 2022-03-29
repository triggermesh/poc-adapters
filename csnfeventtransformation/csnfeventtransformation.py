# Copyright 2022 TriggerMesh Inc.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from flask import Flask, request, make_response
from cloudevents.http import CloudEvent, to_binary, from_http, to_structured
import requests
import simplejson
import json
import os
from onug_decorator import onug


app = Flask(__name__)

aquasecType = os.environ.get('AQUASEC_TYPE')
cloudGuardType = os.environ.get('CLOUDGUARD_TYPE')
azureType = os.environ.get('AZURE_TYPE')


@app.route("/", methods=["POST"])
def home():
    event = from_http(request.headers, request.get_data())
    message = json.loads(request.data.decode('utf-8'))
    provider = ""

    if aquasecType == event['type']:
        provider = "Aquasec"

    if cloudGuardType == event['type']:
        provider = "CloudGuard"

    if azureType == event['type']:
        provider = "Azure"

    if provider == "":
        return make_response("Unknown event type", 400)


    url = 'https://objectstorage.us-ashburn-1.oraclecloud.com/n/orasenatdpltsecitom01/b/HammerPublic/o/file.json'
    onugMSG = onug(url,message, provider)

    attributes = {
        "type": "io.triggermesh.csnf.cloudguard",
        "source": "https://example.com/event-producer",
        "specversion": "1.0",
        "content-type": "application/cloudevents+json",
    }

    data = simplejson.dumps(onugMSG.get_finding())
    print("data---")
    print(data)
    print("data---")
    revent = CloudEvent(attributes, onugMSG.get_finding())
    headers, body = to_structured(revent)
    sink = os.environ.get('K_SINK')

    if sink != "" :
       r = requests.post(sink, data=body, headers=headers)
       print(f"{r}")
        # return "", 200

    return (body, 200)

if __name__ == "__main__":
    app.run(port=8080)
