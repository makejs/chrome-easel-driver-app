/*
  This file is a shim for usage of the http module within the EaselDriver
*/


const EventEmitter = require("events")

// createServer just makes a new instance
exports.createServer = () => {
  return new Server
}

// Socket is a evented tcp connection wrapper
class Socket extends EventEmitter {
  constructor({socketId, localAddress, localPort, peerAddress, peerPort}) {
    super()
    this.socketId = socketId
    this.localAddress = localAddress
    this.localPort = localPort
    this.peerAddress = peerAddress
    this.peerPort = peerPort

    chrome.sockets.tcp.onReceive.addListener(this._recv = ({socketId, data})=>{
      if (socketId!==this.socketId) return
      // keep this in ArrayBuffer format (the default)
      // since it's being passed to the Go runtime which will handle
      // it as bytes
      this.emit("data", data)
    })
    chrome.sockets.tcp.onReceiveError.addListener(this._recvErr = ({socketId, error})=>{
      if (socketId!==this.socketId) return
      // TODO add UI and show errors there
      console.error("tcp.onReceiveError", error);
      this.close()
    })

    // chrome tcp streams start paused, so we need to unpause it
    chrome.sockets.tcp.setPaused(this.socketId, false);
    this.writeBuffer = []
  }

  // this may not be necessary, but it's done this way to keep sends in order
  // and one-at-a-time. If it's making things slow and there's some docs to show
  // it's not needed in chrome, we can just make the calls in .write()
  flush() {
    if (!this.writeBuffer.length) return
    if (this.closed) return
    chrome.sockets.tcp.send(this.socketId, this.writeBuffer.shift(), (code)=>{
      if (code < 0) {
        console.error("tcp.send error:", code)
        this.close()
        return
      }
      this.flush()
    })
  }
  write(data) {
    this.writeBuffer.push(data)
    this.flush()
  }


  close() {
    if (this.closed) return
    this.closed = true

    // cleanup listeners
    chrome.sockets.tcp.onReceive.removeListener(this._recv)
    chrome.sockets.tcp.onReceiveError.removeListener(this._recvErr)

    // close the socket
    // we could wait for the callback, but anyone listening is only looking
    // for close, and should stop sending now anyway
    chrome.sockets.tcp.close(this.socketId)
    this.emit("close")
  }
}


// Server is basically an evented tcp server, using the chrome API
class Server extends EventEmitter {
  constructor() {
    super()
    this.socketId = null

    // we have to create a socket, even before listening
    chrome.sockets.tcpServer.create({}, ({socketId})=>{
      if (chrome.runtime.lastError) {
        throw chrome.runtime.lastError
      }
      this.socketId = socketId

      // emit ready now, even if listen was called, so the order of events
      // is consistent
      this.emit("ready")

      // if .listen was already called, then honor it
      if (this._listen) {
        this._listen()
      }
    })

    // listen for new connections. since no code calls close, we're not handling
    // cleanup (and no close method).
    chrome.sockets.tcpServer.onAccept.addListener(({socketId, clientSocketId})=>{
      if (socketId!==this.socketId) return

      // instead of just accepting it ask for some info first, then create the Socket
      chrome.sockets.tcp.getInfo(clientSocketId, ({localAddress, localPort, peerAddress, peerPort})=>this.emit("accept", new Socket({
        localAddress, localPort, peerAddress, peerPort,
        socketId: clientSocketId
      })))
    })
  }

  // bind the listener to host:port (probably always :1338 for Easel)
  listen(port, host) {
    this.localAddress = host
    this.localPort = port
    // create a function to actually do the listen, in case we haven't received
    // the socketId yet
    const listen = () => {
      chrome.sockets.tcpServer.listen(this.socketId, host, port, result => {
        if (chrome.runtime.lastError) {
          throw chrome.runtime.lastError
        }
        if (result < 0) throw new Error(`failed to listen on: ${host}:${port} (error code: ${result})`)
        // get info about the server
        chrome.sockets.tcpServer.getInfo(this.socketId, ({localAddress, localPort})=> {
          this.localAddress = localAddress
          this.localPort = localPort
          this.emit("listening", `${host}:${port}`)
        })
      })
    }

    // if the socket isn't ready, store the function in this._listen, otherwise just call it
    if (this.socketId === null) {
      this._listen = listen
    } else {
      listen()
    }
  }
}
