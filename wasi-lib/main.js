import { readFile } from 'node:fs/promises';
import { WASI } from 'wasi';
import { argv, env } from 'node:process';
import fs from "fs"

(async function () {
  const wasi = new WASI({
    version: 'preview1',
    args: argv,
    env,
    returnOnExit: true
    // preopens: {
    //   '/sandbox': '/root/now/wasm-runtime',
    // },
  });
  const wasm = await WebAssembly.compile(
    await readFile(new URL(process.argv[2], import.meta.url)),
  );

  console.log("start ")

  let instance

  const hostio = {
    "_gotest": //func get_preimage_len
    {
      get_preimage_len: (keyPtr) => {
        let key = new Uint8Array(instance.exports.memory.buffer,keyPtr, 32)
        // console.log("key is:", key.toString())
  
        //read preimage from file descriptor
        let PClientRFd = 5
        let PClientWFd = 6
        fs.writeSync(PClientWFd, Buffer.from(key))
  
        //write to go-wasm
        let lenBuf = Buffer.alloc(8)
        fs.readSync(PClientRFd,lenBuf,0,8)
        // console.log("lenBuf====>",lenBuf)
        let len = parseInt(lenBuf.toString("hex"),16)
        // console.log("len js:", len)
        return len
      },
  
      //func getKeyFromOracle() []byte
      get_preimage_from_oracle: (keyPtr,offset,len) => {
        let mem = new DataView(instance.exports.memory.buffer)
        let PClientRFd = 5
        let key = new Uint8Array(instance.exports.memory.buffer,keyPtr, 32)
        // console.log("key is:", key.toString())

        let data = Buffer.alloc(len)
        let readed_len = fs.readSync(PClientRFd,data)
        // console.log("read length",readed_len)
        // console.log("read data:",  data.subarray(0,32))
  
        //send data back to go-wasm
        for(let i=0; i< readed_len; i++){
          mem.setUint8(offset,data[i],true)
          offset = offset + 1
        }
        return readed_len
      
      },
  
      "hint_oracle": (retBufPtr, retBufSize) => {
        //load hintstr
        let hintArr = new Uint8Array(instance.exports.memory.buffer,retBufPtr, retBufSize)
        let HClientWFd = 4
        fs.writeSync(HClientWFd, Buffer.from(hintArr))
      },
    }
  }

  let max_mem = 0

  const wasi_imports = wasi.getImportObject()
  const previous_clock_time_get = wasi_imports.wasi_snapshot_preview1.clock_time_get
  wasi_imports.wasi_snapshot_preview1.clock_time_get =  (clockId, precision, time) => {
    let res = previous_clock_time_get(clockId, precision, time)
    if (max_mem < instance.exports.memory.buffer.byteLength) {
      max_mem = instance.exports.memory.buffer.byteLength
    }
    return res
  }

  instance = await WebAssembly.instantiate(wasm, {...wasi_imports,...hostio});
  wasi.start(instance);

  process.on("exit", ()=> {
    console.log("\nmaximum memory usage==========================>",max_mem)
  })

})()

