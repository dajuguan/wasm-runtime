# wasm-runtime
requirement: node.js version >=20.5.1

## build client
```
cd go-test/client/cmd
tinygo build -o client.wasi -target=wasm main.go 
```
## run node-host + wasm-client
```
cd ../..
go run main.go
```