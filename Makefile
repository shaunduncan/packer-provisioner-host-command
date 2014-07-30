all: clean deps build

clean:
	rm -f $(GOPATH)/bin/packer-provisioner-host-command

deps:
	go get github.com/flynn/go-shlex
	go get github.com/mitchellh/packer

build:
	go build -o $(GOPATH)/bin/packer-provisioner-host-command -a
