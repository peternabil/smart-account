postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=simple_calc -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

mock:
	mockgen -source store/store.go -destination mocks/mocks.go

test:
	go test -cover -coverprofile=c.out ./...

coverage:
	go tool cover -html=c.out -o coverage.html

dropdb:
	docker exec -it postgres16 dropdb simple_bank


.PHONY: postgres createdb dropdb mock
