private-key:
	openssl genrsa -out private.key 2048
public-key:
	openssl rsa -in private.key -pubout -out public.key
run-local:
	go run cmd/oauth2-server/main.go
docker-build:
	docker build -t oauth2-server .
docker-run:
	docker run --rm oauth2-server
