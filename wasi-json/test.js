import fs from "node:fs"
let preimages = fs.readFileSync("/root/now/optimism/op-program/bin/preimages.json")
// preimages = fs.readFileSync("/root/now/wasm-runtime/package.json")
let preimages_json = JSON.parse(preimages)
for(let key in preimages_json) {
    console.log("keys", key)
    console.log("value length", preimages_json[key].length)
}