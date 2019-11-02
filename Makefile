.PHONY: build clean funcdeploy deploy funcdeployprod deployprod

all: clean build funcdeploy

dev: clean build funcdeploy

prod: clean build funcdeployprod

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/infinitives lambdas/infinitives/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

funcdeploy: clean build
	sls deploy -f infinitives --verbose

deploy:
	sls deploy --verbose

funcdeployprod:
	sls deploy -f infinitives --stage prod --verbose

deployprod:
	sls deploy --verbose
