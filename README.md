# Countbeat

Countbeat is a [Beat](https://www.elastic.co/beats/) built on libbeat that publishes an
event with an incrementing `counter` field every `countbeat.period` (default `1s`).

## Requirements

* [Go](https://go.dev/dl/) — the version in `go.mod` or newer
* [Docker](https://www.docker.com/), only for `make package`

Dependencies are vendored in `vendor/`, so builds work offline. Build tooling is
[mage](https://magefile.org/), run from `vendor/` through the `Makefile`, so you don't need to install it.

## Build

```
make
```

This builds the `countbeat` binary in the repository root. `go build .` also works.

## Run

```
./countbeat -c countbeat.yml -e
```

To send events to stdout instead of Elasticsearch:

```
./countbeat -c countbeat.yml -e -E output.elasticsearch.enabled=false -E output.console.enabled=true
```

To load the index template and data stream into Elasticsearch before the first run:

```
./countbeat setup --index-management -c countbeat.yml -e
```

## Configuration

| Setting            | Default | Description                  |
|--------------------|---------|------------------------------|
| `countbeat.period` | `1s`    | How often an event is sent.  |

The shipped configs (`countbeat.yml`, `countbeat.reference.yml`, `countbeat.docker.yml`) are
generated from `_meta/config/*.tmpl` and libbeat's templates. Edit the templates, not the
generated files.

## Test

```
make test
```

CI (`.github/workflows/ci.yml`) also does the following:

* runs `govulncheck`, weekly as well as on every PR;
* checks that generated files are up to date;
* runs an integration test against Elasticsearch.

## Update generated files

After changing `_meta/fields.yml` or `_meta/config/*.tmpl`, or after bumping libbeat, run:

```
make update
```

This regenerates `fields.yml`, `include/fields.go`, the `countbeat*.yml` configs and
`docs/reference/`.

## Upgrading libbeat

Beats releases are tagged `v8.x`/`v9.x`, but the Go module path is `github.com/elastic/beats/v7`.
You therefore have to pin a release by commit:

```
go get github.com/elastic/beats/v7@$(gh api repos/elastic/beats/commits/v9.5.4 --jq .sha)
make vendor
make update
```

Then do the following:

1. Copy the `replace` block from beats' `go.mod` at that tag into ours. Go does not apply the
   `replace` directives of dependencies.
2. Update the Elasticsearch image version in the CI integration job.
3. Re-check the entries in `.github/govulncheck-allowlist.txt`.

## Packaging

```
make package
```

This cross-compiles and builds distribution packages in `build/distributions`. It requires Docker.
Use `PLATFORMS` to limit the target platforms and `SNAPSHOT=true` for snapshot builds.

## Other targets

Run `make help` to list all targets.
