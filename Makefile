.PHONY: build

build:
	sam build

sync:
	sam sync --config-env dev --parameter-overrides "ParameterKey=Environment,ParameterValue=dev"
