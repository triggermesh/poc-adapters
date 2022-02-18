
send this
```
curl --location --request POST 'http://fixedwithtojsontransformations-fwtojson.default.tmkongdemo.triggermesh.io' \
--header 'Ce-Specversion: 1.0' \
--header 'Ce-Type: io.triggermesh.sample.event' \
--header 'Ce-Id: 123123' \
--header 'Ce-Source: ser' \
--header 'Content-Type: text/plain' \
--data-raw 'NAME                STATE     TELEPHONE

John Smith          WA        418-Y11-4111

Mary Hartford       CA        319-Z19-4341

Evan Nolan          IL        219-532-c301
'
```

expect a response of
```
{
  "fields": [
    {
      "value": "NAME",
      "spaceLeft": 0,
      "lineNumber": 0
    },
    {
      "value": "STATE",
      "spaceLeft": 16,
      "lineNumber": 0
    },
    {
      "value": " TELEPHONE",
      "spaceLeft": 4,
      "lineNumber": 0
    },
    {
      "value": "John Smith",
      "spaceLeft": 0,
      "lineNumber": 2
    },
    {
      "value": "WA",
      "spaceLeft": 10,
      "lineNumber": 2
    },
    {
      "value": "418-Y11-4111",
      "spaceLeft": 8,
      "lineNumber": 2
    },
    {
      "value": "Mary Hartford",
      "spaceLeft": 0,
      "lineNumber": 4
    },
    {
      "value": " CA",
      "spaceLeft": 6,
      "lineNumber": 4
    },
    {
      "value": "319-Z19-4341",
      "spaceLeft": 8,
      "lineNumber": 4
    },
    {
      "value": "Evan Nolan",
      "spaceLeft": 0,
      "lineNumber": 6
    },
    {
      "value": "IL",
      "spaceLeft": 10,
      "lineNumber": 6
    },
    {
      "value": "219-532-c301",
      "spaceLeft": 8,
      "lineNumber": 6
    }
  ]
}
```
