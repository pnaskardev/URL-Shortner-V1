# URL-Shortner-V1

## Run All Services Together

```bash
make -j2
```

## Protobuf / Contracts

Proto files live in [`contracts/proto/`](contracts/proto/). The generated Go code
in `contracts/gen/go/` is committed to the repo, so whenever you edit a `.proto`
file you must regenerate it.

Run from the `contracts/` directory:

```bash
cd contracts

# (optional) lint before generating
buf lint

# regenerate Go code into contracts/gen/go/
buf generate
```

The workspace (`go.work`) points all services at the local `contracts` module, so
after `buf generate` the new/changed types are picked up immediately — just rebuild:

```bash
# from repo root
cd contracts && go build ./... && cd ../core && go build ./... && cd ../shortener-service && go build ./...
```

Requires the `buf` CLI and the `protoc-gen-go` plugin on your `PATH`.
