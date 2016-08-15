# Chrome Easel Driver App

This Chrome App will act as a Easel local driver on any system that runs Chrome.

## Try It

Consider this very much aplha software, it's only been tested against the *Machine Setup* & homing
features of Easel. I haven't had a chance to try carving anything. So be careful and watch your machine!

To install a pre-built release in chrome:

1. Download the latest .zip from the releases page
2. Extract it somewhere
3. Go to the extensions page
  1. Open the menu (top right, 3 dots)
  2. Go to more tools >
  3. Click extensions
2. Make sure `Developer mode` is **checked** at the top-right
3. Click `Load unpacked extension...`
4. Select the directory where the .zip was extracted
5. Click `Launch`

While the *Easel Driver* window is open, it should respond as the normal local driver would.

## Dev Requirements

You will need the following available for development:

- Node 6
- Go 1.6
- `gopherjs` for compiling Go to javascript
- `xar` archive utility
- `cpio`
- `gunzip`
- EaselDriver-0.2.6.pkg

###  Installing Needed Tools

Node can be installed in your home directory by using [nvm](https://github.com/creationix/nvm), a bash one-liner is available on the repo page. `nvm` handles setting all of the many environment variables node and npm require for use.

Go may be installed via package or, on Linux, by extracting the tarball and setting `GOROOT` to the extracted directory. You will also need to set `GOPATH` to some directory to use for packages. Add `$GOROOT/bin` and `$GOPATH/bin` to your `PATH` variable.

`gopherjs` can be installed/updated with the following command: `go get -u github.com/gopherjs/gopherjs`

`xar` and `cpio` can generally be found in your package manager on Linux

`gunzip` -- you already have this most likely


## Building

Import the iris libs by running `./import.sh EaselDriver-0.2.6.pkg`
where the filename (EaselDriver-...) is the path to the Mac version.

Run `npm install` to install all node-related dependencies

Run `gopherjs get` from the `server/` directory to install all go-related dependencies

To build, run `./build.sh` and the app will get build into the `app/` directory (that you can then load into chrome)

## How It Works

The Easel Local Driver handles enumerating serial devices, talking to them,
and parsing some basic info around the responses (like serial number, grbl info, etc...).

The code in this repo shims the network and serial calls to use `chrome` APIs instead
of the native libraries normally used, making cross-platform easier (since Chrome did the work).

It opens a tcpServer in chrome, and uses the Go runtime http server, and a socket-io compatible library.

Most of the code is around shimming the socket.io/websocket and http server stuff. As far as actually
getting and connecting serial devices in Chrome, that part took very little effort. (*hint hint*)
