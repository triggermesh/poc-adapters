## Map Objects Dataweave Transformation Sample

After deploying the files in this directory, we can send it the following curl request:

```
curl --location --request POST 'http://dataweavetransformations-hello-dw.default.35.238.200.185.sslip.io' \
--header 'Ce-Specversion: 1.0' \
--header 'Ce-Type: io.triggermesh.sample.event' \
--header 'Ce-Id: 123123' \
--header 'Ce-Source: ser' \
--header 'Content-Type: application/xml' \
--data-raw '
        {
            "books": [
              {
                "-category": "cooking",
                "title":"Everyday Italian",
                "author": "Giada De Laurentiis",
                "year": "2005",
                "price": "30.00"
              },
              {
                "-category": "children",
                "title": "Harry Potter",
                "author": "J K. Rowling",
                "year": "2005",
                "price": "29.99"
              },
              {
                "-category": "web",
                "title":  "XQuery Kick Start",
                "author": [
                  "James McGovern",
                  "Per Bothner",
                  "Kurt Cagle",
                  "James Linn",
                  "Vaidyanathan Nagarajan"
                ],
                "year": "2003",
                "price": "49.99"
              },
              {
                "-category": "web",
                "-cover": "paperback",
                "title": "Learning XML",
                "author": "Erik T. Ray",
                "year": "2003",
                "price": "39.95"
              }
            ]
        }'
```
and expect an event in the event-display.
```
Context Attributes,
  specversion: 1.0
  type: io.triggermesh.sample.event
  source: ser
  id: 123123
  time: 2022-02-04T03:33:36.859313165Z
  datacontenttype: application/json
Data,
  {
    "items": [
      {
        "book": {
          "-CATEGORY": "cooking",
          "TITLE": "Everyday Italian",
          "AUTHOR": "Giada De Laurentiis",
          "YEAR": "2005",
          "PRICE": "30.00"
        }
      },
      {
        "book": {
          "-CATEGORY": "children",
          "TITLE": "Harry Potter",
          "AUTHOR": "J K. Rowling",
          "YEAR": "2005",
          "PRICE": "29.99"
        }
      },
      {
        "book": {
          "-CATEGORY": "web",
          "TITLE": "XQuery Kick Start",
          "AUTHOR": [
            "James McGovern",
            "Per Bothner",
            "Kurt Cagle",
            "James Linn",
            "Vaidyanathan Nagarajan"
          ],
          "YEAR": "2003",
          "PRICE": "49.99"
        }
      },
      {
        "book": {
          "-CATEGORY": "web",
          "-COVER": "paperback",
          "TITLE": "Learning XML",
          "AUTHOR": "Erik T. Ray",
          "YEAR": "2003",
          "PRICE": "39.95"
        }
      }
    ]
  }
```
