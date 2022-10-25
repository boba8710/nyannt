# nyann't

Anti nyancat bootloader code

This code will read off 512 bytes from your boot disk and write it back down every half second.

Precompiled 32 and 64 bit binaries are provided for convenience.

In order to compile yourself, clone the repo

`git clone https://github.com/boba8710/nyannt`

install golang 

`apt install golang`

cd into the repository directory

`cd nyannt`

depending on the version of golang available in your package manager, you may need to init the module:

`go mod init`

then build for either 64 bit:

`CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -trimpath -buildvcs=false -o nyannt64 .`

or 32 bit

`CGO_ENABLED=0 GOARCH=386 go build -ldflags="-s -w" -trimpath -buildvcs=false -o nyannt32 .`

Get the binary onto your machine, and make sure you execute it while your MBR is still clean! This means you should run it early in the competiton.

Good luck!
