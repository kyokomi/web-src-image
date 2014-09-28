web-src-image
====

## Description

## Usage

```bash
$ web-src-image -n fileName -o output-dir --url http://hogehoge/fugafuga -p 3
```

- `-n`: baseImageName
- `-o`: output dir path
- `--url`: scraping target url
- `-p`: total page count

## Install

To install, use `go get`:

```bash
$ go get -d github.com/kyokomi/web-src-image
```
To install, [relasePage](https://github.com/kyokomi/web-src-image/releases) download

## Memo

```bash
$ gox -osarch="windows/amd64" -output="_obj/web-src-image"
$ zip _obj/web-src-image.zip _obj/web-src-image.exe
$ ghr -u kyokomi -r web-src-image v0.1.0 _obj/web-src-image.zip
```

## Author

[kyokomi](https://github.com/kyokomi)
