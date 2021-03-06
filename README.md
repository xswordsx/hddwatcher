# Hard Drive watcher

[![Go Reference](https://pkg.go.dev/badge/github.com/xswordsx/hddwatcher.svg)](https://pkg.go.dev/github.com/xswordsx/hddwatcher)

A watcher for hard-drive space which sends emails when it
starts to run out.

## Limitations

⚠ Currently only Windows is supported.

⚠ Requires [Golang](https://golang.org/) 1.16 and above.

## Installing

```console
go get -u github.com/xswordsx/hddwatcher
```

## Building

```console
go build -o hddwatcher ./
```

## Running

Copy [config.example.toml](config.example.toml) and set-up
prefered parameters.

```console
hddwatcher -c config.toml
```

## Testing

```console
go test ./...
```

## Links

- [Example configuration](config.example.toml)

## License

This repository is licensed under the MIT license.
