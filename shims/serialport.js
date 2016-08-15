/*
  This file is a shim for usage of the serialport module within the EaselDriver
*/

const Events = require("events")

// util functions
const ab2str = buf => new TextDecoder("utf-8").decode(buf)
const str2ab = str => new TextEncoder("utf-8").encode(str).buffer


// list will call "cb" with (err, devices)
exports.list = cb => {
  chrome.serial.getDevices(ports=>{
      if (chrome.runtime.lastError) return cb(chrome.runtime.lastError)
      cb(null, ports.map(({vendorId, productId, path, displayName})=>({
        vendorId,
        productId,

        // Easel will callback with just comName when it's time to connect
        comName: path,

        // This is important, it's not exactly what is being asked for
        // but as long as it has "Arduino" or similar, Easel will attempt
        // to connect to it, and that's what we want
        manufacturer: displayName
      })))
  })
}

// available parsers, to replicate the serialport API
const parsers = {
  raw: () => cb => data => cb(data),
  // readline takes a stream of strings
  // and outputs a stream of lines
  readline: (separator="\r\n") => cb => {
    var buf = ""
    return data => {
      buf += data

      var idx
      // use a loop here, in case we get multiple lines in a data packet
      while ((idx = buf.indexOf(separator)) !== -1) {
        cb(buf.slice(0, idx))
        buf=buf.slice(idx+1)
      }
    }
  }
}
exports.parsers = parsers

// SerialPort is the main class for communication
// it needs to emit "open", "close", "data" (using parser), and accept .write()
exports.SerialPort = class SerialPort extends Events {
  // set some defaults
  constructor(path, {baudrate=115200, parser=parsers.raw(), errorCallback=null}) {
    super()
    this.connectionId = null
    this.writeBuffer = []

    // create our parser instance
    const parse = parser(data=>this.emit("data", data))

    // listen for errors, and save the reference for cleanup
    chrome.serial.onReceiveError.addListener(this._recvErr = ({connectionId, error}) => {
      if (connectionId!==this.connectionId) return
      this.emit("error", new Error("receive: "+error))
      this.close()
    })

    // listen for new data, passing it to the `parse` function we created
    chrome.serial.onReceive.addListener(this._recv = ({connectionId, data}) =>{
      if (connectionId!==this.connectionId) return
      // data is passed as an ArrayBuffer, but our parser expects strings
      // so we use the `ab2str` helper
      parse(ab2str(data))
    })

    // now that everything else is setup and listening, actually make the
    // connection
    chrome.serial.connect(path,{bitrate:baudrate}, info=>{
      if (chrome.runtime.lastError) {
        if (errorCallback) {
          errorCallback(chrome.runtime.lastError)
        } else {
          this.emit("error", chrome.runtime.lastError)
        }
        return
      }

      // if all goes well, set our connectionId, and call flush() in case
      // we have pending data.
      this.connectionId = info.connectionId
      this.emit("open")
      this.flush()
    })
  }

  // close will send the disconnect command
  // and clean up the references, and remove the listeners
  close() {
    if (!this.connectionId) return
    chrome.serial.disconnect(this.connectionId, result => {
      if (chrome.runtime.lastError) {
        this.emit("error", chrome.runtime.lastError)
        return
      }
      this.emit("close")
    })
    this.connectionId = null
    this.writeBuffer = null
    chrome.serial.onReceiveError.removeListener(this._recvErr)
    chrome.serial.onReceive.removeListener(this._recv)
  }

  // similar to the http server, this setup is to keep the ordering
  // of messages and ensure only one-at-a-time being sent.
  flush() {
    if (!this.connectionId) return
    if (!this.writeBuffer.length) return

    chrome.serial.send(this.connectionId, this.writeBuffer.shift(), sendInfo => {
      if (sendInfo.error) {
        this.emit("error", new Error("write: " + sendInfo.error))
        this.close()
      } else if (chrome.runtime.lastError) {
        this.emit("error", chrome.runtime.lastError)
        this.close()
      } else {
        // call flush again in case something is left
        this.flush()
      }
    })
  }
  write(data) {
    // convert data to an ArrayBuffer for the chrome API
    this.writeBuffer.push(str2ab(data))
    this.flush()
  }
}
