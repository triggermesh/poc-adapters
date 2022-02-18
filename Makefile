# Create the contianer images for the adapters.
image:
	@cd jqtransformation && gcloud builds submit --tag gcr.io/triggermesh/jqt
	@cd mongodbtarget && gcloud builds submit --tag gcr.io/triggermesh/mongodbtarget
	@cd dataweavetransformation && gcloud builds submit --tag gcr.io/triggermesh/dw
	@cd jsontoxmltransformation && gcloud builds submit --tag gcr.io/triggermesh/jtx
	@cd javascript && gcloud builds submit --tag gcr.io/triggermesh/js
	@cd fixedwithtojson && gcloud builds submit --tag gcr.io/triggermesh/fwtojson
# Apply the Koby manifests for the adapters.
apply:
	@cd jqtransformation/config && kubectl apply -f 100-registration.yaml
	@cd mongodbtarget/config && kubectl apply -f 100-registration.yaml
	@cd dataweavetransformation/config && kubectl apply -f 100-registration.yaml
	@cd jsontoxmltransformation/config && kubectl apply -f 100-registration.yaml
	@cd javascript/config && kubectl apply -f 100-registration.yaml
	@cd fixedwithtojson/config && kubectl apply -f 100-registration.yaml
# Delete the Koby manifests for the adapters.
delete:
	@cd jqtransformation/config && kubectl delete -f 100-registration.yaml
	@cd mongodbtarget/config && kubectl delete -f 100-registration.yaml
	@cd dataweavetransformation/config && kubectl delete -f 100-registration.yaml
	@cd jsontoxmltransformation/config && kubectl delete -f 100-registration.yaml
	@cd javascript/config && kubectl delete -f 100-registration.yaml
	@cd fixedwithtojson/config && kubectl delete -f 100-registration.yaml
# Lint the adapters.
lint:
	@cd jqtransformation/pkg/adapter && golangci-lint run --deadline 2m
	@cd jqtransformation/cmd && golangci-lint run  --deadline 2m
	@cd mongodbtarget/pkg/adapter && golangci-lint run  --deadline 2m
	@cd mongodbtarget/cmd && golangci-lint run  --deadline 2m
	@cd dataweavetransformation/pkg/adapter && golangci-lint run  --deadline 2m
	@cd dataweavetransformation/cmd && golangci-lint run  --deadline 2m
	@cd jsontoxmltransformation/pkg/adapter && golangci-lint run  --deadline 2m
	@cd jsontoxmltransformation/cmd && golangci-lint run  --deadline 2m
	@cd fixedwithtojson/pkg/adapter && golangci-lint run  --deadline 2m
	@cd fixedwithtojson/cmd && golangci-lint run  --deadline 2m
# Build the adapters.
build:
	@cd jqtransformation/cmd && go build
	@cd mongodbtarget/cmd && go build
	@cd dataweavetransformation/cmd && go build
	@cd jsontoxmltransformation/cmd && go build
	@cd fixedwithtojson/cmd && go build
# Clean build artifacts.
clean:
	@cd jqtransformation/cmd && go clean
	@cd mongodbtarget/cmd && go clean
	@cd dataweavetransformation/cmd && go clean
	@cd jsontoxmltransformation/cmd && go clean
	@cd fixedwithtojson/cmd && go clean
# Test the adapters.
test:
	@cd jqtransformation/pkg/adapter && go test
	@cd mongodbtarget/pkg/adapter && go test
	@cd jsontoxmltransformation/pkg/adapter && go test
