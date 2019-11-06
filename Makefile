.PHONY: build clean funcdeploy deploy funcdeployprod deployprod

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/infinitives lambdas/infinitives/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/conjugations lambdas/conjugations/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

funcdeploy: clean build
	sls deploy -f infinitives --verbose
	sls deploy -f conjugations --verbose

deploy:
	sls deploy --verbose

funcdeployprod:
	sls deploy -f infinitives --stage prod --verbose
	sls deploy -f conjugations --stage prod --verbose

deployprod:
	sls deploy --stage prod --verbose
