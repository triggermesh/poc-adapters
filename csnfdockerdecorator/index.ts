import axios from "axios";
const express = require("express");
const { HTTP,CloudEvent } = require("cloudevents");
const app = express();
import {CsnfEvent} from './types';
import DockerhubDecorator from './dockerhub';

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
  var evnt = HTTP.toEvent({ headers: req.headers, body: req.body });
  var csnfEvent: CsnfEvent = evnt.data;
  const dockerhubDecorator = new DockerhubDecorator();
  const decoration = await dockerhubDecorator.decorate(csnfEvent);
  csnfEvent.decoration = decoration;
  console.log(csnfEvent)

  const ce = new CloudEvent({ type:evnt.type+".dockerhub.decorated", source:"dockerhubDecorator", data: csnfEvent });
  const message = HTTP.binary(ce); // Or HTTP.structured(ce)
  message.headers['content-type'] = 'application/json';
  var  ksink = process.env.K_SINK;

  axios.post(ksink,  message.body, {
    headers: message.headers,
    });

  res.status(200).header("Content-Type", "application/json").send(message.body);
});

app.listen(8080, () => {
  console.log("Example app listening on port 8080!");
});
