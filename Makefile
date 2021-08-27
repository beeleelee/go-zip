.PHONY: build

TARGET=./gozip

build: clean
	go build -o ${TARGET} ./cmd

run: build
	${TARGET}


.PHONY: clean
clean:
	-rm -f ${TARGET}


.PHONY: all 

all: build