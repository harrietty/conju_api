.PHONY: build clean funcdeploy deploy

all: clean build funcdeploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/infinitives lambdas/infinitives/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

funcdeploy: clean build
	sls deploy -f infinitives --verbose

deploy:
	sls deploy --verbose
