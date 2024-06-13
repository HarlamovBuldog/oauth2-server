private-key:
	openssl genrsa -out private.key 2048
public-key:
	openssl rsa -in private.key -pubout -out public.key
run-local:
	go run cmd/oauth2-server/main.go
docker-build:
	docker build -t oauth2-server:0.0.1 .
docker-run:
	docker run --rm -p 8080:8080 oauth2-server:0.0.1
