import { readFile } from 'node:fs/promises';
import { WASI } from 'wasi';
import { argv, env } from 'node:process';




(async function () {
  const wasi = new WASI({
    version: 'preview1',
    args: argv,
    env,
    preopens: {
      '/sandbox': '/root/now/wasm-runtime',
    },
  });
  const wasm = await WebAssembly.compile(
    // await readFile(new URL('./demo.wasm', import.meta.url)),
    await readFile(new URL('/root/now/wasm-runtime/go-test/main.wasi', import.meta.url)),
  );
  const instance = await WebAssembly.instantiate(wasm, wasi.getImportObject());
  
  wasi.start(instance);
})()

