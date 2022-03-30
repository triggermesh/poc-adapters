import {v4 as uuid} from 'uuid';
import axios from "axios";
import * as log4js from 'log4js';
const logger = log4js.getLogger('azure-sentinel-dispatcher');
logger.level = process.env.LOG_LEVEL || 'info';
const express = require("express");
const { HTTP } = require("cloudevents");
const app = express();
const msal = require('@azure/msal-node');


const AZURE_SUBSCRIPTION_ID =  process.env.AZURE_SUBSCRIPTION_ID;
const AZURE_RESOURCE_GROUP =  process.env.AZURE_RESOURCE_GROUP;
const AZURE_WORKSPACE =  process.env.AZURE_WORKSPACE;
const AZURE_CLIENT_SECRET =  process.env.AZURE_CLIENT_SECRET;
const AZURE_CLIENT_ID =  process.env.AZURE_CLIENT_ID;
const AZURE_TENANT_ID =  process.env.AZURE_TENANT_ID;


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

  console.log(event);


  const incident = {
      "properties": {
          "severity": "High",
          "status": "Active",
          "title": event.data.event.name,
          "description": event.data.event.shortDescription,
          "additionalData": {
              "alertProductNames": [
                event.data.provider.accountId,
                event.data.resource.region,
                event.data.resource.type + ':' + event.data.resource.name + ':' + event.data.event.guid
              ]
          },
          "labels": [
              {
                  "labelName": event.data.event.name,
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
