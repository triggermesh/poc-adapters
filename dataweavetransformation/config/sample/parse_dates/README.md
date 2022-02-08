## Parse Dates Sample

After deploying the files in this directory, we can send it the following curl request:

```
curl --location --request POST 'http://dataweavetransformations-hello-dw.default.35.238.200.185.sslip.io' \
--header 'Ce-Specversion: 1.0' \
--header 'Ce-Type: io.triggermesh.sample.event' \
--header 'Ce-Id: 123123' \
--header 'Ce-Source: ser' \
--header 'Content-Type: application/xml' \
--data-raw '      <dates>
        <date>26-JUL-16</date>
        <date>27/JUL/16</date>
        <date>28.JUL.16</date>
      </dates>'
```
and expect an event like this in the event-display:

```
☁️  cloudevents.Event
Context Attributes,
  specversion: 1.0
  type: io.triggermesh.sample.event
  source: ser
  id: 123123
  time: 2022-02-03T20:56:53.27449787Z
  datacontenttype: application/xml
Data,

<?xml version='1.0' encoding='UTF-8'?>
<dates>
  <normalized_as_string>26-JUL-16</normalized_as_string>
  <normalized_as_string>27-JUL-16</normalized_as_string>
  <normalized_as_string>28-JUL-16</normalized_as_string>
</dates>
```
