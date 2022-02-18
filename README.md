# poc-adapters
This repo contains the following poc adapters:
* MongoDBTarget
* JQTransformation
* DataweaveTransformation
* JSONToXMLTransformation


# Usage
## Installation
These adapters can be applied to a kubernetes cluster running Koby by executing the following command:
```cmd
make apply
```

Sample usage can be found in the `config/samples` folder nested within each adapter.

## Removal
To remove the adapters from a kubernetes cluster, execute the following command:
```cmd
make delete
```

## Building the images from scratch
To build the images from scratch, execute the following command:
```cmd
make build
```

## Building the binaries
To build the binaries, execute the following command:
```cmd
make build
```

## Clean build artifacts
To clean build artifacts, execute the following command:
```cmd
make clean
```

## Run tests
To run tests, execute the following command:
```cmd
make test
```

# Deploy the Samples
Each one of the POC adapters, except MongoDB, has a sample deployment/flow that can be used to test out the adapter and see how it works. The sample lives in the `config/samples` folder of each adapter. It can be deployed by executing the following command:
```cmd
kubectl apply -f <adapter>/config/sample/
```
Each one of the samples runs on its own, with data coming from a PingSource object and ending up in an event-display sink. To verify the results, view the logs of the event-display pod.


# Devlopment
To create a new adapter, execute the following command:
```
./hack/scaffold.bash <adapter-name>
```

This will scafold out a new golang adapter project.
