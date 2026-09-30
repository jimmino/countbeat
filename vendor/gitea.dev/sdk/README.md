# Gitea SDK for Go

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Release](https://raster.shields.io/badge/dynamic/json.svg?label=release&url=https://gitea.com/api/v1/repos/gitea/go-sdk/releases&query=$[0].tag_name)](https://gitea.com/gitea/go-sdk/releases)
[![Join the chat at https://img.shields.io/discord/322538954119184384.svg](https://img.shields.io/discord/322538954119184384.svg)](https://discord.gg/Gitea)
[![Go Report Card](https://goreportcard.com/badge/gitea.dev/sdk)](https://goreportcard.com/report/gitea.dev/sdk)
[![GoDoc](https://pkg.go.dev/badge/gitea.dev/sdk)](https://pkg.go.dev/gitea.dev/sdk)

This project acts as a client SDK implementation written in Go to interact with the Gitea API implementation. For further informations take a look at the current [documentation](https://pkg.go.dev/gitea.dev/sdk).

The SDK module path is now `gitea.dev/sdk`. If you are migrating from the previous `code.gitea.io/sdk/gitea` path, see [`docs/migrate-code-gitea-io-to-gitea-dev.md`](docs/migrate-code-gitea-io-to-gitea-dev.md). If you are migrating from direct `Client` business methods to service-based entry points, see [`docs/migrate-client-to-services.md`](docs/migrate-client-to-services.md).

Note: function arguments are escaped by the SDK.

## Use it

```go
import "gitea.dev/sdk"
```

## Version Requirements
 * go >= 1.26
 * gitea >= 1.11

## Contributing

Fork -> Patch -> Push -> Pull Request

Run common development commands from the repository root, for example `make build`,
`make test-unit`, and `make lint`.

The compatibility wrappers in `client_services_wrappers.go` are generated and
deprecated in favor of `client.<Service>.<Method>`. For the full generated
migration matrix from `Client` methods to services, see
[`docs/migrate-client-to-services.md`](docs/migrate-client-to-services.md). Run
`go generate ./...` or `make generate` after changing exported service methods.

## Authors

* [Maintainers](https://github.com/orgs/go-gitea/people)
* [Contributors](https://github.com/go-gitea/go-sdk/graphs/contributors)

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the full license text.
