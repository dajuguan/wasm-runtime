// Assume add.wasm file exists that contains a single function adding 2 provided arguments
// const fs = require('fs');

import fs from "fs"

const wasmBuffer = fs.readFileSync('/root/now/wasm-runtime/go-test/main.wasm');
WebAssembly.instantiate(wasmBuffer, importObject).then(wasmModule => {
  // Exported function live under instance.exports
  const { _start } = wasmModule.instance.exports;
  _start();
//   const sum = add(5, 6);
//   console.log(sum); // Outputs: 11
});

importObject = {
    wasi_snapshot_preview1: {
        // https://github.com/WebAssembly/WASI/blob/main/phases/snapshot/docs.md#fd_write
        fd_write: function(fd, iovs_ptr, iovs_len, nwritten_ptr) {
            let nwritten = 0;
            if (fd == 1) {
                for (let iovs_i=0; iovs_i<iovs_len;iovs_i++) {
                    let iov_ptr = iovs_ptr+iovs_i*8; // assuming wasm32
                    let ptr = mem().getUint32(iov_ptr + 0, true);
                    let len = mem().getUint32(iov_ptr + 4, true);
                    nwritten += len;
                    for (let i=0; i<len; i++) {
                        let c = mem().getUint8(ptr+i);
                        if (c == 13) { // CR
                            // ignore
                        } else if (c == 10) { // LF
                            // write line
                            let line = decoder.decode(new Uint8Array(logLine));
                            logLine = [];
                            console.log(line);
                        } else {
                            logLine.push(c);
                        }
                    }
                }
            } else {
                console.error('invalid file descriptor:', fd);
            }
            mem().setUint32(nwritten_ptr, nwritten, true);
            return 0;
        },
        fd_close: () => 0,      // dummy
        fd_fdstat_get: () => 0, // dummy
        fd_seek: () => 0,       // dummy
        "proc_exit": (code) => {
            if (global.process) {
                // Node.js
                process.exit(code);
            } else {
                // Can't exit in a browser.
                throw 'trying to exit with code ' + code;
            }
        },
        random_get: (bufPtr, bufLen) => {
            crypto.getRandomValues(loadSlice(bufPtr, bufLen));
            return 0;
        },
    },
    env: {
        // func ticks() float64
        "runtime.ticks": () => {
            return timeOrigin + performance.now();
        },

        // func sleepTicks(timeout float64)
        "runtime.sleepTicks": (timeout) => {
            // Do not sleep, only reactivate scheduler after the given timeout.
            setTimeout(this._inst.exports.go_scheduler, timeout);
        },

        // func finalizeRef(v ref)
        "syscall/js.finalizeRef": (sp) => {
            // Note: TinyGo does not support finalizers so this should never be
            // called.
            console.error('syscall/js.finalizeRef not implemented');
        },

        // func stringVal(value string) ref
        "syscall/js.stringVal": (ret_ptr, value_ptr, value_len) => {
            const s = loadString(value_ptr, value_len);
            storeValue(ret_ptr, s);
        },

        // func valueGet(v ref, p string) ref
        "syscall/js.valueGet": (retval, v_addr, p_ptr, p_len) => {
            let prop = loadString(p_ptr, p_len);
            let value = loadValue(v_addr);
            let result = Reflect.get(value, prop);
            storeValue(retval, result);
        },

        // func valueSet(v ref, p string, x ref)
        "syscall/js.valueSet": (v_addr, p_ptr, p_len, x_addr) => {
            const v = loadValue(v_addr);
            const p = loadString(p_ptr, p_len);
            const x = loadValue(x_addr);
            Reflect.set(v, p, x);
        },

        // func valueDelete(v ref, p string)
        "syscall/js.valueDelete": (v_addr, p_ptr, p_len) => {
            const v = loadValue(v_addr);
            const p = loadString(p_ptr, p_len);
            Reflect.deleteProperty(v, p);
        },

        // func valueIndex(v ref, i int) ref
        "syscall/js.valueIndex": (ret_addr, v_addr, i) => {
            storeValue(ret_addr, Reflect.get(loadValue(v_addr), i));
        },

        // valueSetIndex(v ref, i int, x ref)
        "syscall/js.valueSetIndex": (v_addr, i, x_addr) => {
            Reflect.set(loadValue(v_addr), i, loadValue(x_addr));
        },

        // func valueCall(v ref, m string, args []ref) (ref, bool)
        "syscall/js.valueCall": (ret_addr, v_addr, m_ptr, m_len, args_ptr, args_len, args_cap) => {
            const v = loadValue(v_addr);
            const name = loadString(m_ptr, m_len);
            const args = loadSliceOfValues(args_ptr, args_len, args_cap);
            try {
                const m = Reflect.get(v, name);
                storeValue(ret_addr, Reflect.apply(m, v, args));
                mem().setUint8(ret_addr + 8, 1);
            } catch (err) {
                storeValue(ret_addr, err);
                mem().setUint8(ret_addr + 8, 0);
            }
        },

        // func valueInvoke(v ref, args []ref) (ref, bool)
        "syscall/js.valueInvoke": (ret_addr, v_addr, args_ptr, args_len, args_cap) => {
            try {
                const v = loadValue(v_addr);
                const args = loadSliceOfValues(args_ptr, args_len, args_cap);
                storeValue(ret_addr, Reflect.apply(v, undefined, args));
                mem().setUint8(ret_addr + 8, 1);
            } catch (err) {
                storeValue(ret_addr, err);
                mem().setUint8(ret_addr + 8, 0);
            }
        },

        // func valueNew(v ref, args []ref) (ref, bool)
        "syscall/js.valueNew": (ret_addr, v_addr, args_ptr, args_len, args_cap) => {
            const v = loadValue(v_addr);
            const args = loadSliceOfValues(args_ptr, args_len, args_cap);
            try {
                storeValue(ret_addr, Reflect.construct(v, args));
                mem().setUint8(ret_addr + 8, 1);
            } catch (err) {
                storeValue(ret_addr, err);
                mem().setUint8(ret_addr+ 8, 0);
            }
        },

        // func valueLength(v ref) int
        "syscall/js.valueLength": (v_addr) => {
            return loadValue(v_addr).length;
        },

        // valuePrepareString(v ref) (ref, int)
        "syscall/js.valuePrepareString": (ret_addr, v_addr) => {
            const s = String(loadValue(v_addr));
            const str = encoder.encode(s);
            storeValue(ret_addr, str);
            setInt64(ret_addr + 8, str.length);
        },

        // valueLoadString(v ref, b []byte)
        "syscall/js.valueLoadString": (v_addr, slice_ptr, slice_len, slice_cap) => {
            const str = loadValue(v_addr);
            loadSlice(slice_ptr, slice_len, slice_cap).set(str);
        },

        // func valueInstanceOf(v ref, t ref) bool
        "syscall/js.valueInstanceOf": (v_addr, t_addr) => {
             return loadValue(v_addr) instanceof loadValue(t_addr);
        },

        // func copyBytesToGo(dst []byte, src ref) (int, bool)
        "syscall/js.copyBytesToGo": (ret_addr, dest_addr, dest_len, dest_cap, source_addr) => {
            let num_bytes_copied_addr = ret_addr;
            let returned_status_addr = ret_addr + 4; // Address of returned boolean status variable

            const dst = loadSlice(dest_addr, dest_len);
            const src = loadValue(source_addr);
            if (!(src instanceof Uint8Array || src instanceof Uint8ClampedArray)) {
                mem().setUint8(returned_status_addr, 0); // Return "not ok" status
                return;
            }
            const toCopy = src.subarray(0, dst.length);
            dst.set(toCopy);
            setInt64(num_bytes_copied_addr, toCopy.length);
            mem().setUint8(returned_status_addr, 1); // Return "ok" status
        },

        // copyBytesToJS(dst ref, src []byte) (int, bool)
        // Originally copied from upstream Go project, then modified:
        //   https://github.com/golang/go/blob/3f995c3f3b43033013013e6c7ccc93a9b1411ca9/misc/wasm/wasm_exec.js#L404-L416
        "syscall/js.copyBytesToJS": (ret_addr, dest_addr, source_addr, source_len, source_cap) => {
            let num_bytes_copied_addr = ret_addr;
            let returned_status_addr = ret_addr + 4; // Address of returned boolean status variable

            const dst = loadValue(dest_addr);
            const src = loadSlice(source_addr, source_len);
            if (!(dst instanceof Uint8Array || dst instanceof Uint8ClampedArray)) {
                mem().setUint8(returned_status_addr, 0); // Return "not ok" status
                return;
            }
            const toCopy = src.subarray(0, dst.length);
            dst.set(toCopy);
            setInt64(num_bytes_copied_addr, toCopy.length);
            mem().setUint8(returned_status_addr, 1); // Return "ok" status
        },
    }
};

