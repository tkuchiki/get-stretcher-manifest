# get-stretcher-manifest
Get manifest file (See: https://github.com/fujiwara/stretcher)

# Installation

## Linux

```
$ curl -sLO https://github.com/tkuchiki/get-stretcher-manifest/releases/download/v0.0.1/get-manifest-linux-amd64.zip
```


## Mac OSX
```
$ curl -sLO https://github.com/tkuchiki/get-stretcher-manifest/releases/download/v0.0.1/get-manifest-darwin-amd64.zip
```

# Build

```
go get
go build -o get-manifest main.go
```

# Usage

```
$ get-manifest --help
usage: get-manifest --bucket=BUCKET [<flags>]

Flags:
  --help               Show help (also see --help-long and --help-man).
  -a, --all            All manifests
  -n, --num=1          N-th Manifest
  -b, --bucket=BUCKET  Bucket
  --region="ap-northeast-1"
                       Region
  -f, --file=FILE      Credentials file(Default ~/.aws/credentials, ~/.aws/config)
  --profile="default"  Profile
  --oldest             The oldest manifest
  --version            Show application version.
```
