generate:
	go run github.com/steebchen/prisma-client-go generate

db-push:
	go run github.com/steebchen/prisma-client-go db push

start:
	air

.PHONY: generate db-push start
