/*
  This file is a shim for usage of socket.io within the EaselDriver
*/

// easiest way I found to import the GopherJS-compiled module
// was to just make it global. So we just require it, and rely
// on the side-effect
var NewServer = require("../server")

// listen expects a server (shims/http.js)
exports.listen = httpServer => NewServer(httpServer)
