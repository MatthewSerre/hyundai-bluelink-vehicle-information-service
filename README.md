# Hyundai Bluelink Vehicle Information Service

A repository for an vehicle information service that uses server-side code generated from the `hyundai-bluelink-protobufs` repository [link](https://github.com/MatthewSerre/hyundai-bluelink-protobufs) to respond to requests for vehicle information via Hyundai Bluelink.

## Getting Started

* [Install Go](https://go.dev/doc/install)

## Usage

Run `go run internal/server/server.go` to start the server, which will accept requests at the specified local address. This is a standalone service even though it is grouped with other repositories associated with this project. What this means is that it can be used or modified for use in any Go application that needs an vehicle information response from Hyundai Bluelink that fits the data type specified by the protobufs contract.

## Contributing

Create an issue and/or a pull request and I will take a look.

***

This project is not affiliated with Hyundai in any way. Credit to [TaiPhamD](https://github.com/TaiPhamD) and his `bluelink_go` project [link](https://github.com/TaiPhamD/bluelink_go) for inspiration and some code snippets.