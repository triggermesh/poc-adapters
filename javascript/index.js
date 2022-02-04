
const express = require('express')
const app = express()
const port = 8080
const { HTTP } = require("cloudevents");

app.post("/", (req, res) => {
  const receivedEvent = HTTP.toEvent({ headers: req.headers, body: req.body });
  console.log(receivedEvent);
 res.sendStatus(200)
});


app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})
