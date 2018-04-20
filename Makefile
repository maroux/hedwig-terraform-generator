.PHONY: test

test_setup:
	./scripts/test-setup.sh

test: clean test_setup
	./scripts/run-tests.sh

build:
	go-bindata -prefix "templates/" templates/

	env GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/hedwig-terraform-generator .
	env GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64/hedwig-terraform-generator .
	cd bin/linux-amd64 && zip hedwig-terraform-generator-linux-amd64.zip hedwig-terraform-generator; cd -
	cd bin/darwin-amd64 && zip hedwig-terraform-generator-darwin-amd64.zip hedwig-terraform-generator; cd -

clean:
	rm -rf bin bindata.go hedwig
