# Go-eCAL

Go bindings for [eCAL](https://github.com/eclipse-ecal/ecal) 6 (currently unreleased).
Inspired from [Blutkoete/golang-ecal](https://github.com/Blutkoete/golang-ecal).

## Usage

```sh
# Go binding demo
go run .

# Go version of eCAL monitor
go run ./cmd/monitor
```

## Features

- eCAL 6 compatible (unreleased)
- Pure cgo; no SWIG dependency
- Custom C interface implementation
- Direct deserialization from subscriber buffer to Go types
- Rewrite of eCAL monitor using [bubbletea](https://github.com/charmbracelet/bubbletea)

Provides Go interfaces for:
- [x] Core
- [ ] Configuration
- [x] Publisher
  - [ ] Zero Copy
- [x] Subscriber
- [x] Message Types
  - [x] Generic
  - [x] String
  - [x] Protobuf
- [x] Logging
- [ ] Services
- [x] Monitoring
  - [x] Publisher/Subscribers
  - [x] Processes
  - [ ] Server/Clients
- [ ] Registration
  - [x] Topic callbacks

## CLI Tools

### Monitor

Features:

- [x] Show Publishers and Subscribers
  - [x] Detailed topic view
  - [ ] Live message view
- [ ] Show Services
- [ ] Show Hosts
- [ ] Show Processes
- [ ] Show eCAL Logs
- [ ] Show config

## Non-system installations

If eCAL is not installed in a default search path or you wish to use a specific
install of eCAL there is a helper CMake project to generate a `package_user.go`
with the correct `cgo` flags.

```sh
cmake -S . -B build -DCMAKE_PREFIX_PATH=/path/to/cmake/install
go run .
```

## Development

To help write the C and C++ wrapper use the CMake project to generate a
`compile_commands.json`.

