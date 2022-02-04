include .env


clean:
		@rm -rf dist
		@mkdir -p dist

build: clean 
		GOOS=linux GOARCH=amd64 go build -o ./dist/main 

run:
		aws-sam-local local start-api

preinstall:
		curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
		unzip awscliv2.zip
		sudo ./aws/install
		sudo apt-get install zip

install:
		go get -t ./...

install-dev:
		go get -t ./...

test:
		go test ./... --cover

lint: 
		go get -u golang.org/x/lint/golint; \
		golint -set_exit_status; \
      	go vet . 

configure:
		aws s3api create-bucket \
			--bucket $(AWS_BUCKET_NAME) \
			--region $(AWS_REGION) \
			--create-bucket-configuration LocationConstraint=$(AWS_REGION)
# aws cloudformation package --template-file template.dev.yml --s3-bucket throoo-artifacts --s3-prefix dev/login --output-template-file outputtemplate.yml
package: build
		aws cloudformation package \
			--template-file $(TEMPLATE) \
			--s3-bucket $(AWS_BUCKET_NAME) \
			--s3-prefix $(AWS_BUCKET_FOLDER_PREFEX) \
			--region us-east-1 \
			--output-template-file package.yml

deploy:
		aws cloudformation deploy \
			--template-file package.yml \
			--region $(AWS_REGION) \
			--capabilities CAPABILITY_IAM \
			--stack-name $(AWS_STACK_NAME)

describe:
		aws cloudformation describe-stacks \
			--region $(AWS_REGION) \
			--stack-name $(AWS_STACK_NAME) \

outputs:
		@make describe | jq -r '.Stacks[0].Outputs'

url:
		@make describe | jq -r ".Stacks[0].Outputs[0].OutputValue" -j