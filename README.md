## cf-check
Check an Host is Owned by CloudFlare.

## Install

1. Grab from [releases page](https://github.com/dwisiswant0/cf-check/releases), or
2. If you have [Go1.18+](https://go.dev/dl/) compiler installed & configured:

```console
$ go install github.com/dwisiswant0/cf-check@latest
```

## Usage

```console
$ echo "uber.com" | cf-check
34.98.127.226
```

or

```console
$ cf-check -d <FILE>
```

### Flags

```console
$ cf-check -h
Usage of cf-check:
  -c int
        Set the concurrency level (default: 20)
  -cf
        Show CloudFlare only
  -d    Print domains instead of IP addresses
```

## Workaround

The goal is that you don't need to do a port scan if it's proven that the IP is owned by Cloudflare.

```console
$ subfinder -silent -d uber.com | filter-resolved | cf-check -d | anew | naabu -silent -verify | httpx -silent
```

## License

`cf-check` is distributed under Apache License 2.0.