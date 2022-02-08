## XML to JSON Dataweave Transformation Sample

After deploying the files in this directory, we can send it the following curl request:

```
curl --location --request POST 'http://dataweavetransformations-hello-dw.default.35.238.200.185.sslip.io' \
--header 'Ce-Specversion: 1.0' \
--header 'Ce-Type: io.triggermesh.sample.event' \
--header 'Ce-Id: 123123' \
--header 'Ce-Source: ser' \
--header 'Content-Type: application/xml' \
--data-raw '<?xml version='\''1.0'\'' encoding='\''UTF-8'\''?>
<order>
  <product>
    <price>5</price>
    <model>MuleSoft Connect 2016</model>
  </product>
  <item_amount>3</item_amount>
  <payment>
    <payment-type>credit-card</payment-type>
    <currency>USD</currency>
    <installments>1</installments>
  </payment>
  <buyer>
    <email>mike@hotmail.com</email>
    <name>Michael</name>
    <address>Koala Boulevard 314</address>
    <city>San Diego</city>
    <state>CA</state>
    <postCode>1345</postCode>
    <nationality>USA</nationality>
  </buyer>
  <shop>main branch</shop>
  <salesperson>Mathew Chow</salesperson>
</order>'
```

and expect an event like this in the event-display:
```
☁️  cloudevents.Event
Context Attributes,
  specversion: 1.0
  type: io.triggermesh.sample.event
  source: ser
  id: 123123
  time: 2022-02-03T20:55:02.202047009Z
  datacontenttype: application/json
Data,
  {
    "address1": "Koala Boulevard 314",
    "city": "San Diego",
    "country": "USA",
    "email": "mike@hotmail.com",
    "name": "Michael",
    "postalCode": "1345",
    "stateOrProvince": "CA"
  }
```
