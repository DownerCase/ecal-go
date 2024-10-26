# Go-eCAL

Go bindings for [eCAL](https://github.com/eclipse-ecal/ecal) 6 (currently unreleased).
Inspired from [Blutkoete/golang-ecal](https://github.com/Blutkoete/golang-ecal).

## Usage

Run the demo:

```sh
go run .
```

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

