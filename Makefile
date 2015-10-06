all: clean deps build

clean:
	rm -f $(GOPATH)/bin/packer-provisioner-host-command

deps:
	go get github.com/flynn/go-shlex
	go get github.com/mitchellh/packer
	go get github.com/mitchellh/multistep
	go get github.com/mitchellh/go-fs
	go get code.google.com/p/go.crypto/ssh

build:
	go build -o $(GOPATH)/bin/packer-provisioner-host-command -a

test:
	go test
