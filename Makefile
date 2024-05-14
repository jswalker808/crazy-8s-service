.PHONY: build

build:
	sam build

unit-tests:
	cd crazy-8s && go test -short -v ./... && cd ..

integ-tests:
	cd crazy-8s && go test -run Integration -v ./... && cd ..

sync:
	sam sync --config-env dev --parameter-overrides "ParameterKey=Environment,ParameterValue=dev"
