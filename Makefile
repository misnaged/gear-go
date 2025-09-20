APP:=gear-go
COMMON_PATH	?= $(shell pwd)
APP_ENTRY_POINT:=cmd/gear-go.go
RUST_GEN_ENTRY=cmd/generate-rust-grpc/generate-rust-grpc.go
EXAMPLES_ENTRY_POINT=example/
GITVER_PKG:=github.com/misnaged/scriptorium/versioner
BUILD_OUT_DIR:=./
GOPRIVATE:=github.com
CARGO_DIR := lib/temp/
RUST_GRPC := lib/server_grpc/

GOOS	:=
GOARCH	:=
ifeq ($(OS),Windows_NT)
	GOOS =windows
	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		OSFLAG =amd64
	endif
	ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		OSFLAG =ia32
	endif
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		GOOS =linux
	endif
	ifeq ($(UNAME_S),Darwin)
		GOOS =darwin
	endif
		UNAME_P := $(shell uname -m)
	ifeq ($(UNAME_P),x86_64)
		GOARCH =amd64
	endif
	ifneq ($(filter %86,$(UNAME_P)),)
		GOARCH =386
	endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		GOARCH =arm
	endif
endif

TAG 		:= $(shell git describe --abbrev=0 --tags)
COMMIT		:= $(shell git rev-parse HEAD)
BRANCH		?= $(shell git rev-parse --abbrev-ref HEAD)
REMOTE		:= $(shell git config --get remote.origin.url)
BUILD_DATE	:= $(shell date +'%Y-%m-%dT%H:%M:%SZ%Z')

RELEASE :=
ifeq ($(TAG),)
	RELEASE := $(COMMIT)
else
	RELEASE := $(TAG)
endif

LDFLAGS += -X $(GITVER_PKG).ServiceName=$(APP)
LDFLAGS += -X $(GITVER_PKG).CommitTag=$(TAG)
LDFLAGS += -X $(GITVER_PKG).CommitSHA=$(COMMIT)
LDFLAGS += -X $(GITVER_PKG).CommitBranch=$(BRANCH)
LDFLAGS += -X $(GITVER_PKG).OriginURL=$(REMOTE)
LDFLAGS += -X $(GITVER_PKG).BuildDate=$(BUILD_DATE)

all: tidy build

tidy:
	go mod tidy


update:
	go get -u ./...


build:
	env CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-w -s ${LDFLAGS}" -o $(BUILD_OUT_DIR)/$(APP) $(APP_ENTRY_POINT)


test:
	go test -mod=readonly  ./...

run:
	MallocNanoZone=0 go run -race $(APP_ENTRY_POINT) serve

example-code-run:
	MallocNanoZone=0 go run -race $(EXAMPLES_ENTRY_POINT)code/example_upload_code_and_get_code_from_storage.go

example-update-code-sub:
	MallocNanoZone=0 go run -race $(EXAMPLES_ENTRY_POINT)code/upload_and_create/example_subscription_upload.go


cargo-build:
	cd $(CARGO_DIR) && cargo install subxt-cli
	cd $(CARGO_DIR) && subxt metadata -f bytes > metadata.scale
	cd $(CARGO_DIR) && cargo build

cargo-run:
	cd $(CARGO_DIR) && cargo run

generate-proto-go:
	cd lib/server_grpc/proto && protoc --go_out=paths=source_relative:. --go_opt=paths=source_relative  --go-grpc_out=paths=source_relative:. --go-grpc_opt=paths=source_relative  *.proto

generate-rust-grpc:
	cd $(RUST_GRPC) && cargo install subxt-cli
	cd $(RUST_GRPC) && subxt metadata -f bytes > metadata.scale
	cd -
	MallocNanoZone=0 go run -race $(RUST_GEN_ENTRY) generate-rust-grpc true
	MallocNanoZone=0 go run -race $(RUST_GEN_ENTRY) remove-target
	if [ -f "$(RUST_GRPC)/gear_grpc.rs" ]; then \
			MallocNanoZone=0 go run -race $(RUST_GEN_ENTRY) remove-rust; \
	fi

	if [ ! -d "$(RUST_GRPC)/target" ]; then \
		cd $(RUST_GRPC) && cargo build; \
	fi
	MallocNanoZone=0 go run -race $(RUST_GEN_ENTRY) generate-rust-grpc false

run-rust-grpc:
	cd $(RUST_GRPC) && cargo run