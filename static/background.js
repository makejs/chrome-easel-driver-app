chrome.app.runtime.onLaunched.addListener(function() {
  // disconnect/close all serial connections and tcp server sockets
  const sockets = new Promise(resolve=>chrome.sockets.tcpServer.getSockets(resolve))
  .then(sockets=>Promise.all(sockets.map(({socketId})=>new Promise(resolve=>chrome.sockets.tcpServer.close(socketId, resolve)))))
  const sConns = new Promise(resolve=>chrome.serial.getConnections(resolve))
  .then(conns=>Promise.all(conns.map(({connectionId})=>new Promise(resolve=>chrome.serial.disconnect(connectionId, resolve)))))

  Promise.all([sockets, sConns])
  .then(()=>chrome.app.window.create('window.html', {
    'id': "",
    singleton: true,
    'outerBounds': {
      'width': 400,
      'height': 500
    }
  }))

});
