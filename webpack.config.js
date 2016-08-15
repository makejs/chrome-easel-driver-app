const path = require("path")
const webpack = require("webpack")

module.exports = {
  entry: "./iris-lib/iris",
  output: {
    filename: "app/bundle.js"
  },
  resolve: {
    alias: {
      readline: "readline-browserify",
      "socket.io": require.resolve("./shims/socket.io"),
      http: require.resolve("./shims/http"),
      fs: require.resolve("./shims/fs"),
      serialport: require.resolve("./shims/serialport")
    }
  },
  module: {
    loaders: [
      {test: /\.js$/, exclude: /node_modules|iris-lib/, loader: "babel"},
      {test: /\.json$/, loader: "json"}
    ]
  },
  plugins: [
    new webpack.DefinePlugin({
      "process.platform": JSON.stringify("Chrome")
    })
  ]
}
