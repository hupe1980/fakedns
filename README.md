# fakedns
> Tiny DNS proxy for Penetration Testers and Malware Analysts

## Features
- Regular Expression based DNS server
- IPV4 & IPV6
- DNS Rebinding
- DNS Round-Robin
- Upstream DNS Resolver

## Installing
You can install the pre-compiled binary in several different ways

### homebrew tap:
```bash
brew tap hupe1980/fakedns
brew install fakedns
```

### snapcraft:
[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/fakedns)
```bash
sudo snap install fakedns
```

### scoop:
```bash
scoop bucket add fakedns https://github.com/hupe1980/fakedns-bucket.git
scoop install fakedns
```

### deb/rpm/apk:

Download the .deb, .rpm or .apk from the [releases page](https://github.com/hupe1980/fakedns/releases) and install them with the appropriate tools.

### manually:
Download the pre-compiled binaries from the [releases page](https://github.com/hupe1980/fakedns/releases) and copy to the desired location.

## Building from source
Install a [Go 1.17 compiler](https://golang.org/dl) and run the following command in the checked-out repository:
```
make build
```

(Add the appropriate .exe extension on Windows systems, of course.)

## Docker Support
```bash
git clone https://github.com/hupe1980/fakedns
cd fakedns
make docker-build

# Run fakedns for example.org on port 4444
docker run -it --rm -p 4444:53/udp fakedns example.org -v --ipv4 127.0.0.1

# Use dig to test fakedns
dig example.org @localhost -p 4444
```
## How to use
```console
Usage:
  fakedns [domains] [flags]

Examples:
- IPV4: fakedns example.org --ipv4 127.0.0.1
- Wildcards: fakedns example.* --ipv4 127.0.0.1
- RoundRobin: fakedns example.org --ipv4 127.0.0.1,10.10.10.10
- Rebind: fakedns example.org --ipv4 127.0.0.1 --rebind-v4 10.10.10.10
- Upstream: fakedns example.org --ipv4 127.0.0.1 --upstream 8.8.8.8

Flags:
  -a, --addr string            fakeDNS address (default "0.0.0.0:53")
  -h, --help                   help for fakedns
      --ipv4 strings           IPV4 address to return
      --ipv6 strings           IPV6 address to return
      --mx string              host name of mail exchange server
  -n, --net string             fakeDNS network protocol (default "udp")
      --rebind-threshold int   rebind threshold (default 1)
      --rebind-v4 string       IPV4 rebind address
      --rebind-v6 string       IPV6 rebind address
      --text strings           TXT text value
      --ttl uint32             time to live (default 60)
      --upstream string        upstream dns server
  -v, --verbose                print detailed logging messages
      --version                version for fakedns
```

## License
[MIT](LICENCE)
