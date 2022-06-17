all: include install

install:
	@go install .

include:
	cp -r google ${GOPATH}/src/