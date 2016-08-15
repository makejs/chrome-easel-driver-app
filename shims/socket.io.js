/*
  This file is a shim for usage of socket.io within the EaselDriver
*/

// listen expects a server (shims/http.js)
// NewSocketServer is provided by the gopherjs build output
exports.listen = httpServer => NewSocketServer(httpServer)
