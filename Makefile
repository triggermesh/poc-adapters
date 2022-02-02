# Create the contianer images for the adapters.
image:
	@cd jqtransformation && gcloud builds submit --tag gcr.io/ultra-hologram-297914/jqt
	@cd mongodbtarget && gcloud builds submit --tag gcr.io/ultra-hologram-297914/mongodbtarget
	@cd dataweavetransformation && gcloud builds submit --tag gcr.io/ultra-hologram-297914/dw
	@cd jsontoxmltransformation && gcloud builds submit --tag gcr.io/ultra-hologram-297914/jtx
# Apply the Koby manifests for the adapters.
apply:
	@cd jqtransformation/config && kubectl apply -f 100-registration.yaml
	@cd mongodbtarget/config && kubectl apply -f 100-registration.yaml
	@cd dataweavetransformation/config && kubectl apply -f 100-registration.yaml
	@cd jsontoxmltransformation/config && kubectl apply -f 100-registration.yaml
# Delete the Koby manifests for the adapters.
delete:
	@cd jqtransformation/config && kubectl delete -f 100-registration.yaml
	@cd mongodbtarget/config && kubectl delete -f 100-registration.yaml
	@cd dataweavetransformation/config && kubectl delete -f 100-registration.yaml
	@cd jsontoxmltransformation/config && kubectl delete -f 100-registration.yaml
# Lint the adapters.
lint:
	@cd jqtransformation/pkg/adapter && golangci-lint run
	@cd jqtransformation/cmd && golangci-lint run
	@cd mongodbtarget/pkg/adapter && golangci-lint run
	@cd mongodbtarget/cmd && golangci-lint run
	@cd dataweavetransformation/pkg/adapter && golangci-lint run
	@cd dataweavetransformation/cmd && golangci-lint run
	@cd jsontoxmltransformation/pkg/adapter && golangci-lint run
	@cd jsontoxmltransformation/cmd && golangci-lint run
# Build the adapters.
build:
	@cd jqtransformation/cmd && go build
	@cd mongodbtarget/cmd && go build
	@cd dataweavetransformation/cmd && go build
	@cd jsontoxmltransformation/cmd && go build
# Clean build artifacts.
clean:
	@cd jqtransformation/cmd && go clean
	@cd mongodbtarget/cmd && go clean
	@cd dataweavetransformation/cmd && go clean
	@cd jsontoxmltransformation/cmd && go clean
# Test the adapters.
test:
	@cd jqtransformation/pkg/adapter && go test
	@cd mongodbtarget/pkg/adapter && go test
	@cd jsontoxmltransformation/pkg/adapter && go test
