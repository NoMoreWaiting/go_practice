BINPATH = ../bin
TAGET = stockbase

PWD := $(shell pwd)
export GOPATH=$(PWD)/../../public/gopath/:$(PWD)/../../../

all: clean
	mkdir -p $(BINPATH)
	go build -o $(BINPATH)/$(TAGET)

clean:
	rm -f $(BINPATH)/$(TAGET)
