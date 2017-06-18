export GOPATH := $(PWD)
export GOBIN := $(GOPATH)/bin
export PACKAGES := $(shell env GOPATH=$(GOPATH) go list ./src/staticd/...)

all: clean predependencies dependencies build

clean:
	rm -vf bin/*

build: build-macos build-linux build-windows

build-macos: build-macos-i386 build-macos-amd64

build-linux: build-linux-i386 build-linux-amd64

build-windows: build-windows-i386 build-windows-amd64

build-macos-i386:
	GOOS=darwin GOARCH=386 CGO=0 go build -o bin/staticd-darwin-i386 staticd/cmd

build-macos-amd64:
	GOOS=darwin GOARCH=amd64 CGO=0 go build -o bin/staticd-darwin-amd64 staticd/cmd

build-linux-i386:
	GOOS=linux GOARCH=386 CGO=0 go build -o bin/staticd-linux-i386 staticd/cmd

build-linux-amd64:
	GOOS=linux GOARCH=amd64 CGO=0 go build -o bin/staticd-linux-amd64 staticd/cmd

build-windows-i386:
	GOOS=windows GOARCH=386 CGO=0 go build -o bin/staticd-windows-i386.exe staticd/cmd

build-windows-amd64:
	GOOS=windows GOARCH=amd64 CGO=0 go build -o bin/staticd-windows-amd64.exe staticd/cmd

dependencies:
	cd src && trash

predependencies:
	go get -u github.com/rancher/trash

sign:
	gpg --detach-sign --digest-algo SHA512 --no-tty --batch --output bin/staticd-darwin-i386.sig 				bin/staticd-darwin-i386
	gpg --detach-sign --digest-algo SHA512 --no-tty --batch --output bin/staticd-darwin-amd64.sig 			bin/staticd-darwin-amd64
	gpg --detach-sign --digest-algo SHA512 --no-tty --batch --output bin/staticd-linux-i386.sig					bin/staticd-linux-i386
	gpg --detach-sign --digest-algo SHA512 --no-tty --batch --output bin/staticd-linux-amd64.sig 				bin/staticd-linux-amd64
	gpg --detach-sign --digest-algo SHA512 --no-tty --batch --output bin/staticd-windows-i386.exe.sig 	bin/staticd-windows-i386.exe
	gpg --detach-sign --digest-algo SHA512 --no-tty --batch --output bin/staticd-windows-amd64.exe.sig 	bin/staticd-windows-amd64.exe

verify:
	gpg --verify bin/staticd-darwin-i386.sig 				bin/staticd-darwin-i386
	gpg --verify bin/staticd-darwin-amd64.sig 			bin/staticd-darwin-amd64
	gpg --verify bin/staticd-linux-i386.sig					bin/staticd-linux-i386
	gpg --verify bin/staticd-linux-amd64.sig 				bin/staticd-linux-amd64
	gpg --verify bin/staticd-windows-i386.exe.sig 	bin/staticd-windows-i386.exe
	gpg --verify bin/staticd-windows-amd64.exe.sig 	bin/staticd-windows-amd64.exe
