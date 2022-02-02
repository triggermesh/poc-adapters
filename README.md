This repo contains the following adapters:
* MongoDBTarget
* JQTransformation
* DataweaveTransformation
* JSONToXMLTransformation

These adapters can be applied to a kubernetes cluster running Koby by executing the following command:
```cmd
make apply
```

Sample usage can be found in the `config/samples` folder nested within each adapter.
