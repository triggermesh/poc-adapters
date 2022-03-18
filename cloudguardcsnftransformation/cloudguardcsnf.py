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


# create an endpoint at http://localhost:/8080/
@app.route("/", methods=["POST"])
def home():
    event = from_http(request.headers, request.get_data())
    message = json.loads(request.data.decode('utf-8'))

    url = 'https://objectstorage.us-ashburn-1.oraclecloud.com/n/orasenatdpltsecitom01/b/HammerPublic/o/file.json'
    onugMSG = onug(url,message)

    attributes = {
        "type": "io.triggermesh.csnf.cloudguard",
        "source": "https://example.com/event-producer",
        "specversion": "1.0",
        "content-type": "application/cloudevents+json",
    }

    data = simplejson.dumps(onugMSG.get_finding())
    print(data)

    revent = CloudEvent(attributes, onugMSG)
    headers, body = to_binary(revent)
    sink = os.environ.get('K_SINK')

    if sink != "" :
        requests.post(sink, data=data, headers=headers)
        return "", 200

    return (data, 200)

if __name__ == "__main__":
    app.run(port=8080)
