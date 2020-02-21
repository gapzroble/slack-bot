SHELL := /bin/bash

build: clean test
	GOARCH=amd64 GOOS=linux go build -o ./bin/translate ./translate

deps:
	GOPRIVATE=github.com go mod vendor

clean:
	ls -I*.sh ./bin | xargs -I {} rm -f ./bin/{}

zip_handlers: build
	ls -I*.zip -I*.sh ./bin | xargs -I {} zip -j ./bin/{}.zip ./bin/{}

package: zip_handlers
	sam package --template-file ./template.yml --output-template-file ./packaged.yml --s3-bucket gapz.deploys

deploy: package
	sam deploy --template-file ./packaged.yml --stack-name Translate-Bot-Test --parameter-overrides Stage=test DeployBucket=gapz.deploys --capabilities CAPABILITY_NAMED_IAM

test:
	@go test -v $$(go list ./...) >/tmp/gotesting || (grep -A 1 "FAIL:" /tmp/gotesting  && false)
	@echo PASS
