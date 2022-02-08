## Rename JSON Keys Sample

This DataWeave example renames some keys in a JSON object, while retaining the names of all others in the output.

After deploying the DataweaveTransformation, we can send it the following curl request:

```
curl --location --request POST 'http://dataweavetransformations-hello-dw.default.35.238.200.185.sslip.io' \
--header 'Ce-Specversion: 1.0' \
--header 'Ce-Type: io.triggermesh.sample.event' \
--header 'Ce-Id: 123123' \
--header 'Ce-Source: ser' \
--header 'Content-Type: application/json' \
--data-raw '{
  "flights":[
  {
  "availableSeats":45,
  "airlineName":"Ryan Air",
  "aircraftBrand":"Boeing",
  "aircraftType":"737",
  "departureDate":"12/14/2017",
  "origin":"BCN",
  "destination":"FCO"
  },
  {
  "availableSeats":15,
  "airlineName":"Ryan Air",
  "aircraftBrand":"Boeing",
  "aircraftType":"747",
  "departureDate":"08/03/2017",
  "origin":"FCO",
  "destination":"DFW"
  }]
}'
```

and expect an event like this in the event-display:
```
☁️  cloudevents.Event
Context Attributes,
  specversion: 1.0
  type: io.triggermesh.sample.event
  source: ser
  id: 123123
  time: 2022-02-03T19:53:54.562188279Z
  datacontenttype: application/json
Data,
  [
    {
      "emptySeats": 45,
      "airline": "Ryan Air",
      "aircraftBrand": "Boeing",
      "aircraftType": "737",
      "departureDate": "12/14/2017",
      "origin": "BCN",
      "destination": "FCO"
    },
    {
      "emptySeats": 15,
      "airline": "Ryan Air",
      "aircraftBrand": "Boeing",
      "aircraftType": "747",
      "departureDate": "08/03/2017",
      "origin": "FCO",
      "destination": "DFW"
    }
  ]
```
