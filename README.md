# Go-eCAL

Go bindings and CLI tools for [eCAL](https://github.com/eclipse-ecal/ecal) 6 (currently unreleased).
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
- [x] Configuration (partial, expanded as required)
- [x] Publisher & Subscriber
- [x] Message Types: Custom/binary, String, Protobuf
- [x] Logging
- [x] Monitoring ("slow" API for detailed information)
- [ ] Registration ("fast" API for basic topic name and datatype information)
  - [x] Topic event (added/removed publisher/subscriber) callbacks 
- [ ] Services

## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/string/publisher"
)

func main() {
	// Initialize eCAL
	ecal.Initialize(
		ecal.NewConfig(ecal.WithConsoleLogging(true), ecal.WithConsoleLogAll()),
		"Go eCAL!",
		ecal.CDefault,
	)
	defer ecal.Finalize() // Shutdown eCAL at the end of the program

	// Send messages
	go func() {
		publisher, err := publisher.New("string topic")
		if err != nil {
			panic("Failed to make string publisher")
		}
		defer publisher.Delete()

		for idx := 0; true; idx++ {
			_ = publisher.Send("This is message ", idx)

			time.Sleep(time.Second)
		}
	}()

	// Receive messages
	subscriber, err := ecal.NewStringSubscriber("string topic")
	if err != nil {
		panic("Failed to Create string subscriber")
	}
	defer subscriber.Delete()

	for ecal.Ok() {
		msg, err := subscriber.Receive(time.Second * 2)
		if err == nil {
			fmt.Println("Received:", msg)
		} else {
			fmt.Println(err)
		}
	}
}
```

## CLI Tools

### Monitor

Features:

All the basic functionality from the upstream TUI monitor but:

- More efficient use of space; use the _whole_ screen to view your messages!
- Easily send stop commands to eCAL processes (make `eCAL::Ok()` return `false`)

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

