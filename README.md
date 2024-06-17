# Shrt

[![Go Reference](https://pkg.go.dev/badge/djmo.ch/go-shrt.svg)](https://pkg.go.dev/djmo.ch/go-shrt)

This is a URL shortener.
There are many like it, but this is mine.
Try putting all the vowels in "shrt" and enjoy yourself!

Also handles go-get redirects.

## HTTP Handler

This Go http.Handler module is imported as djmo.ch/go-shrt.

To use, initialize Shrt with a shrt.Config object.
Drop this Handler into your site's http.ServeMux and start serving
shortlinks and go-get redirects.

## CLI

The command line interface (CLI) in this repository is useful if you
wish to run Shrt as a standalone server.
It can is installed in the usual manner:

```
$ go install djmo.ch/go-shrt/cmd/shrt@latest
```

From there you can run "shrt help" to read the CLI documentation.

## Releases

Release artifacts created prior to December 1, 2022 can be verified
in one of two ways:

1. Each .tar.gz archive is signed inline with `signify` and can be
   verified (using the -z flag) against the public key located at
   https://www.danielmoch.com/static/djmoch-signify.pub

2. Using the .sig file corresponding to an archive, verification can
   be done with `gpg` against the public key located at
   https://www.danielmoch.com/static/gpg.asc

Archives created on or after December 1, 2022 will only be signed with
`signify`, and not with `gpg`.
To support the use of `minisign`, which is compatable with `signify`,
detached signatures will be created.
These signature files are indicated with a .minisig extension.

## License

ISC.
See the [LICENSE] file for full copyright and license details.

[LICENSE]: LICENSE
