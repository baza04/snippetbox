.SILENT:

build:
	go build -o ./bin/app ./cmd/web/.

run: build
	./bin/app	
# go run ./cmd/web/app.go

# test:
