image:
	@cd jqtransformation && gcloud builds submit --tag gcr.io/ultra-hologram-297914/jqt
	@cd mongodbtarget && gcloud builds submit --tag gcr.io/ultra-hologram-297914/mongodbtarget
	@cd dataweavetransformation && gcloud builds submit --tag gcr.io/ultra-hologram-297914/dw

apply:
	@cd jqtransformation/config && kubectl apply -f 100-registration.yaml
	@cd mongodbtarget/config && kubectl apply -f 100-registration.yaml
	@cd dataweavetransformation/config && kubectl apply -f 100-registration.yaml

delete:
	@cd jqtransformation/config && kubectl delete -f 100-registration.yaml
	@cd mongodbtarget/config && kubectl delete -f 100-registration.yaml
	@cd dataweavetransformation/config && kubectl delete -f 100-registration.yaml

lint:
	@cd jqtransformation/pkg/adapter && golangci-lint run
	@cd jqtransformation/cmd && golangci-lint run
	@cd mongodbtarget/pkg/adapter && golangci-lint run
	@cd mongodbtarget/cmd && golangci-lint run
	@cd dataweavetransformation/pkg/adapter && golangci-lint run
	@cd dataweavetransformation/cmd && golangci-lint run

build:
	@cd jqtransformation/cmd && go build
	@cd mongodbtarget/cmd && go build
	@cd dataweavetransformation/cmd && go build

clean:
	@cd jqtransformation/cmd && go clean
	@cd mongodbtarget/cmd && go clean
	@cd dataweavetransformation/cmd && go clean
