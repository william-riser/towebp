# towebp

Converts a directory of images to `.webp` in parallel.

## Requirements

- Go 1.22+
- [`github.com/chai2010/webp`](https://github.com/chai2010/webp)

## Install

```sh
go mod init towebp
go get github.com/chai2010/webp
go build -o towebp towebp.go
```

## Usage

```
./towebp [flags]

  -dir      directory to scan (default: current directory, recursive)
  -quality  webp quality 0–100 (default: 80)
  -workers  parallel workers (default: number of CPUs)
```

**Example**

```sh
./towebp -dir ./photos -quality 85 -workers 8
```

Output files are written alongside the originals with a `.webp` extension.

## Supported formats

`.jpg`, `.jpeg`, `.png`

> **RAW / DNG files** are not supported by Go's standard image decoder. Pre-convert them to TIFF or PNG first with `dcraw`, `libraw`, or `magick input.dng output.png`.
