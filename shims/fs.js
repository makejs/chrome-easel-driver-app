
exports.readFileSync = file => {
  if (file !== "/package.json") throw new Error("unknown file for shim: " + file)
  return JSON.stringify(require("../iris-lib/package.json"))
}
