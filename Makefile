.PHONY: build

build:
	sam build

unit-tests:
	cd crazy-8s && go test -v ./... && cd ..

sync:
	sam sync --config-env dev --parameter-overrides "ParameterKey=Environment,ParameterValue=dev"
