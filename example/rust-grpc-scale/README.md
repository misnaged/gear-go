### To avoid scale codec golang implementations which are being built heavily based on `go/reflect` 
### it's possible to encode data natively on rust. It requires rust grpc server as a proxy, golang grpc client
### and shared .proto file



## Requirements:

1.  **gear node** 
2. shared .proto file
3. rust grpc-server 
4. go grpc-client

### 1.
build manually, or download bin release from gear-tech repository https://github.com/gear-tech/gear

run node locally:
>  gear --dev

### 2.

**Generate client(go) protoc part:**
> make generate-proto

### 3.

grpc-server root location is: **lib/server_grpc/**
as it is represented in Makefile: `RUST_GRPC := lib/server_grpc/`

build rust grpc server, call (from project root):
> make generate-rust-grpc

and then start it:
> make run-rust-grpc

### 4.
*todo: add customisation to both cli and server grpc* 


