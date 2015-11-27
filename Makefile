all: build

build:
	docker build -t pydima/go-thumbnailer .

run:
	docker run -ti --rm -v ${CURDIR}:/go/src/github.com/pydima/go-thumbnailer/ pydima/go-thumbnailer 

test:
	go test ./...
