## cf-check
Check an Host is Owned by CloudFlare.

## Install

1. Grab from [releases page](https://github.com/dwisiswant0/cf-check/releases), or
2. If you have [Go1.13+](https://go.dev/dl/) compiler installed & configured:

```
▶ go install github.com/dwisiswant0/cf-check@latest
```

## Usage

```
▶ echo "uber.com" | cf-check
```

or

```
▶ cf-check -d <FILE>
```

### Flags

```
Usage of cf-check:
  -c int
        Set the concurrency level (default: 20)
  -cf
        Show CloudFlare only
  -d    Print domains instead of IP addresses
```

## Workaround

The goal is that you don't need to do a port scan if it's proven that the IP is owned by Cloudflare.

```
▶ subfinder -silent -d uber.com | filter-resolved | cf-check -d | anew | naabu -silent -verify | httpx -silent
```