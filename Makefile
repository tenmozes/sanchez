.PNONY: all build push


all: build

build:
	GOOS=linux GOARCH=arm GOARM=6 go build -o bin/sanchez github.com/tenmozes/sanchez/cmd/sanchez

push:build
	git add bin/sanchez
	git commit -m "autocommit. update binary"
	git push