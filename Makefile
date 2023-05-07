build:
	go build -o ./bin/main ./src

start: build
	./bin/main
  
package: build
	fyne package -os darwin --executable bin/main -icon icon/icon.png