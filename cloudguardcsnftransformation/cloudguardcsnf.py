from flask import Flask, request, make_response
from cloudevents.http import CloudEvent, to_binary, from_http, to_structured


import requests
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
        "type": event['type'] + ".response",
        "source": "https://example.com/event-producer",
    }

    # data = simplejson.dumps(onugMSG)
    revent = CloudEvent(attributes, onugMSG)
    headers, body = to_binary(revent)
    sink = os.environ.get('K_SINK')

    if sink != "" :
        requests.post(sink, data=body, headers=headers)
        return "", 200

    return (body, 200)

if __name__ == "__main__":
    app.run(port=8080)
