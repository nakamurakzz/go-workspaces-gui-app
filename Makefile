build:
	go build -o ./bin/main ./src

start: build
	./bin/main
  