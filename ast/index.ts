import {v4 as uuid} from 'uuid';
import axios from "axios";
import * as log4js from 'log4js';
const logger = log4js.getLogger('azure-sentinel-dispatcher');
logger.level = process.env.LOG_LEVEL || 'info';
const express = require("express");
const { CloudEvent, HTTP } = require("cloudevents");
const app = express();
const msal = require('@azure/msal-node');


const AZURE_SUBSCRIPTION_ID = "77641a71-ffc3-4cfd-abd9-6ff8dc509a3d"
const AZURE_RESOURCE_GROUP = "sent"
const AZURE_WORKSPACE = "sent"
const AZURE_CLIENT_SECRET = "Fll7Q~fQ4uu_EkmvYAorlEO196CDJ6osTlC1C"
const AZURE_CLIENT_ID = "90f08ae5-11ce-46fb-a424-b54ffd3a25f7"
const AZURE_TENANT_ID = "f14eddee-e73b-481d-8237-17983764afcb"


const AZURE_API_VERSION = '2020-01-01';

const incidentId = uuid();

const requestUrl = `https://management.azure.com/subscriptions/${AZURE_SUBSCRIPTION_ID}/`
          + `resourceGroups/${AZURE_RESOURCE_GROUP}/`
          + `providers/Microsoft.OperationalInsights/workspaces/${AZURE_WORKSPACE}/`
          + `providers/Microsoft.SecurityInsights/incidents/${incidentId}?api-version=${AZURE_API_VERSION}`;


app.use((req, res, next) => {
  let data = "";

  req.setEncoding("utf8");
  req.on("data", function (chunk) {
    data += chunk;
  });

  req.on("end", function () {
    req.body = data;
    next();
  });
});

app.post("/", async (req, res) => {
  const msalConfig = {
    auth: {
        clientId: AZURE_CLIENT_ID,
        authority: "https://login.microsoftonline.com/" + AZURE_TENANT_ID,
        clientSecret: AZURE_CLIENT_SECRET,
    }
  };
  const cca = new msal.ConfidentialClientApplication(msalConfig);

  const tokenRequest = {
    scopes: [ 'https://management.azure.com/.default' ],
};

const getTokenResponse = await cca.acquireTokenByClientCredential(tokenRequest);
  var event

  try {
    event = HTTP.toEvent({ headers: req.headers, body: req.body });

  } catch (err) {
    console.error(err);
    res.status(415).header("Content-Type", "application/json").send(JSON.stringify(err));
  }


  const incident = {
      "properties": {
          "severity": "High",
          "status": "Active",
          "title": event.data.event.event.metadata.name,
          "description": event.data.event.event.metadata.shortDescription,
          "additionalData": {
              "alertProductNames": [
                event.data.event.event.resources[0].platform,
                event.data.event.event.resources[0].accountId,
                event.data.event.event.resources[0].region,
                event.data.event.event.resources[0].service,
                event.data.event.event.resources[0].type + ':' + event.data.event.event.resources[0].name + ':' + event.data.event.event.resources[0].guid
              ]
          },
          "labels": [
              {
                  "labelName": event.data.event.event.reporter.name,
                  "labelType": "User"
              }
          ]
      }
  };


  try {
    const result = await axios.put(requestUrl, incident,{
        headers: {
            'Authorization': `Bearer ${getTokenResponse.accessToken}`
        }
    });
    console.log(result);
    if (result.status === 200 || result.status === 201) {
        logger.debug('incident successfully dispatched');
        console.log('incident successfully dispatched');
    } else {
        logger.error('failed to dispatch incident', result);
        res.status(500).send(result);
    }
} catch (e) {
    logger.error('failed to dispatch incident', e);
    res.status(500).send(e);
}

  res.status(200).header("Content-Type", "application/json").send("ok");
});

app.listen(8080, () => {
  console.log("Example app listening on port 8080!");
});
