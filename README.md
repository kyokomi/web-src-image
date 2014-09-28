web-src-image
====

## Description

## Usage

```bash
$ web-src-image {outputDirPath} {webURL}
```

- Args1: images outputDirPath
- Args2: scraping web URL

## Install

To install, use `go get`:

```bash
$ go get -d github.com/kyokomi/web-src-image
```

## Memo

```bash
$ gox -osarch="windows/amd64" -output="_obj/web-src-image"
$ zip _obj/web-src-image.zip _obj/web-src-image.exe
$ ghr -u kyokomi -r web-src-image v0.1.0 _obj/web-src-image.zip
```

## Author

[kyokomi](https://github.com/kyokomi)
