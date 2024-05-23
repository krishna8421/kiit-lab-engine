generate:
	go run github.com/steebchen/prisma-client-go generate

db-push:
	go run github.com/steebchen/prisma-client-go db push

start:
	go run main.go

clean:
	rm -rf db

.PHONY: generate db-push start clean
