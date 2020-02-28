SAM_DIRECTORY=web-functions
S3_BUCKET=sam-templates-robekoc
STACK_NAME=guardian
REGION=ap-southeast-2

.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./bin
	
install:
	go install

test:
	go test ./...

pre-build:
	mkdir -p bin

build-cli: pre-build
	go build -o ./bin/guardian

build-api: pre-build
	GOOS=linux GOARCH=amd64 go build -o ./bin/vpn/create ./${SAM_DIRECTORY}/vpn-create
	GOOS=linux GOARCH=amd64 go build -o ./bin/vpn/list ./${SAM_DIRECTORY}/vpn-list
	GOOS=linux GOARCH=amd64 go build -o ./bin/vpn/delete ./${SAM_DIRECTORY}/vpn-delete

build-all: pre-build build-api build-cli

deploy: 
	sam validate
	sam package --s3-bucket ${S3_BUCKET}-${REGION} --template-file template.yaml --output-template-file packaged.yaml
	sam deploy --template-file packaged.yaml --stack-name ${STACK_NAME} --capabilities CAPABILITY_IAM --region ${REGION}

dev:
	sam validate
	sam local