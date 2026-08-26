.PHONY: sealer pit test

sealer:
	cd sealer && go test ./... && go build -o pit-sealer .

pit:
	cd pit && go test ./... && go build -o pit ./cmd/pit

test:
	cd pit && go test ./...
	cd sealer && go test ./...
