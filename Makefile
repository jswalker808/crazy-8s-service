.PHONY: build

build:
	sam build

unit-tests:
	cd crazy-8s && go test -v ./... && cd ..

integ-tests:
	cd crazy-8s && go test -tags=integration ./...

sync:
	sam sync --config-env dev --parameter-overrides "ParameterKey=Environment,ParameterValue=dev"
