# Publishing to the Go module proxy / pkg.go.dev

This module is **not** uploaded to a central registry. Discovery uses the public
GitHub repository plus the [Go module proxy](https://proxy.golang.org/) indexing
model described in the official docs:

- [Publishing a module](https://go.dev/doc/modules/publishing)
- [About pkg.go.dev](https://pkg.go.dev/about)
- [Managing module source](https://go.dev/doc/modules/managing-source)

## Prerequisites (before tagging)

1. Default branch `main` contains the release commit with:
   - `go.mod` module path `github.com/sherotree/image-to-pixel-art-go`
   - MIT `LICENSE` (required for full pkg.go.dev documentation)
   - Passing `gofmt`, `go vet`, `go test` / `go test -race`, `go build`
2. Repository is **public** at https://github.com/sherotree/image-to-pixel-art-go
3. Maintainer explicitly approves the **semver** and publication scope
   (do not invent tags from automation).

## Release steps (after approval)

From a clean clone of this repository on the commit to publish:

```bash
go mod tidy
go test ./...
git tag v0.1.0
git push origin v0.1.0
```

Then **request indexing** so proxy.golang.org caches the tag (pick one):

```bash
GOPROXY=proxy.golang.org go list -m github.com/sherotree/image-to-pixel-art-go@v0.1.0
```

or:

```bash
curl -sI https://proxy.golang.org/github.com/sherotree/image-to-pixel-art-go/@v/v0.1.0.info
```

or open https://pkg.go.dev/github.com/sherotree/image-to-pixel-art-go@v0.1.0
and use the **Request** button if the page is not yet populated.

Consumers then install with:

```bash
go get github.com/sherotree/image-to-pixel-art-go@v0.1.0
```

Docs URL (after index catch-up, typically a few minutes):

https://pkg.go.dev/github.com/sherotree/image-to-pixel-art-go@v0.1.0

## Notes

- Tagged module content is immutable via the checksum database (`sum.golang.org`).
  Never force-push a retagged `vX.Y.Z`; publish a new patch/minor instead.
- This repository hosts a **single module at the repo root**, so tags are plain
  `vX.Y.Z` (no subdirectory prefix).
- `v0` versions are experimental per Go module conventions; reserve `v1` for a
  stable API commitment.
