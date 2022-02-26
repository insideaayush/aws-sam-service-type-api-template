.PHONY: build deps start test db format

deps:
	go install github.com/golang/mock/mockgen@v1.5.0
	go mod download

build:
	sam build

devtools:
	go build -o ./build/bin/sqs-emu cmd/sqs-emu/main.go

start:
	make build && sam local start-api --port 8080 --region ap-southeast-1 --parameter-overrides "ParameterKey=EnvType,ParameterValue=dev" --env-vars local_env.json

test:
	EnvType=dev go test ./... -cover -covermode=atomic -coverprofile=coverage.txt

db:
	docker run -d -p 8000:8000 amazon/dynamodb-local

cover:
	go tool cover -html=coverage.txt

format:
	go fmt ./...

clean:
	rm ./build/bin/*
	