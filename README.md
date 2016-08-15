# Chrome Easel Driver App

This Chrome App will act as a Easel local driver on any system that runs Chrome.

## Requirements

You will need the following available for development:

- Node 6
- Go 1.6
- `gopherjs` for compiling Go to javascript
- `xar` archive utility
- `cpio`
- `gunzip`
- EaselDriver-0.2.6.pkg

###  Installing Requirements

Node can be installed in your home directory by using [nvm](https://github.com/creationix/nvm), a bash one-liner is available on the repo page. `nvm` handles setting all of the many environment variables node and npm require for use.

Go may be installed via package or, on Linux, by extracting the tarball and setting `GOROOT` to the extracted directory. You will also need to set `GOPATH` to some directory to use for packages. Add `$GOROOT/bin` and `$GOPATH/bin` to your `PATH` variable.

`gopherjs` can be installed/updated with the following command: `go get -u github.com/gopherjs/gopherjs`

`xar` and `cpio` can generally be found in your package manager on Linux

`gunzip` -- you already have this most likely


## Building

## How It Works
