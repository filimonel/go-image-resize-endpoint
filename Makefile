build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/main ./cmd/main.go
deploy_prod: build
	serverless deploy --stage prod
undeploy_prod:
	serverless remove --stage prod --region us-east-1