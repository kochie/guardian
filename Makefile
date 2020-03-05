SAM_DIRECTORY=web-functions
S3_BUCKET=sam-templates-robekoc
STACK_NAME=guardian
REGION=ap-southeast-2

.PHONY: deps clean build

fmt:
	gofmt -w .

pre-commit: fmt

deps:
	go get -u ./...

clean: 
	rm -rf ./bin
	
install:
	go install

test-cli:
	go test ./cmd/...

test-api:
	go test ./web-functions/...

test: test-cli test-api

pre-build:
	mkdir -p bin

build-cli: pre-build
	go build -o ./bin/guardian

build-api: pre-build
	GOOS=linux GOARCH=amd64 go build -o ./bin/vpn/create ./${SAM_DIRECTORY}/vpn-create
	GOOS=linux GOARCH=amd64 go build -o ./bin/vpn/list ./${SAM_DIRECTORY}/vpn-list
	GOOS=linux GOARCH=amd64 go build -o ./bin/vpn/delete ./${SAM_DIRECTORY}/vpn-delete

build-auth: pre-build
	GOOS=linux GOARCH=amd64 go build -o ./bin/cognito/challenge ./${SAM_DIRECTORY}/cognito-challenge
	GOOS=linux GOARCH=amd64 go build -o ./bin/cognito/define ./${SAM_DIRECTORY}/cognito-define
	GOOS=linux GOARCH=amd64 go build -o ./bin/cognito/verify ./${SAM_DIRECTORY}/cognito-verify
	GOOS=linux GOARCH=amd64 go build -o ./bin/cognito/presignup ./${SAM_DIRECTORY}/cognito-presignup
	GOOS=linux GOARCH=amd64 go build -o ./bin/cognito/post_authentication ./${SAM_DIRECTORY}/cognito-post-authentication

build-all: pre-build build-api build-cli build-auth

deploy: 
	sam validate
	sam package --s3-bucket ${S3_BUCKET}-${REGION} --template-file template.yaml --output-template-file packaged.yaml
	sam deploy --template-file packaged.yaml --stack-name ${STACK_NAME} --capabilities CAPABILITY_IAM --region ${REGION}

dev:
	sam validate
	sam local