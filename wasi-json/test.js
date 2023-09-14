import fs from "node:fs"
let preimages = fs.readFileSync("/root/now/optimism/op-program/bin/preimages.json")
// preimages = fs.readFileSync("/root/now/wasm-runtime/package.json")
let preimages_json = JSON.parse(preimages)
console.log("total preimages is: ", Object.keys(preimages_json).length)
for(let key in preimages_json) {
    console.log("keys", key)
    let val = Buffer.from(preimages_json[key],'hex')
    console.log("value length", val.length)
    console.log("val is",val.subarray(0,8))
    process.exit()
}