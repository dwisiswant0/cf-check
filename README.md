## cf-check
Check an Host is Owned by CloudFlare.

## Install
```
▶ go get -u github.com/dwisiswant0/cf-check
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
▶ subfinder -silent -d uber.com | filter-resolved | cf-check | anew | naabu -silent -verify | httpx -silent
```