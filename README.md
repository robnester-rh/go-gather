# go-gather

go-gather is a library for Go (golang) for downloading ("gathering") from various sources. These sources include:
 * filepaths
 * git
 * http
 * oci

go-gather simplifies the process of gathering from these sources by freeing the implementer from having to be concerned about the details of the sources.

go-gather is based heavily on [go-getter](https://github.com/hashicorp/go-getter), but designed for specific requirements by [ec](https://github.com/conforma/cli) project.

## Features

- **Automatic source detection** — URI-based routing selects the correct gatherer (file, git, HTTP, OCI) without caller intervention
- **Hardened OCI registry matching** — anchored, escaped regex patterns with subdomain-aware matching prevent spoofing via look-alike hostnames
- **Robust URL parsing** — proper hostname extraction with case normalization and port stripping for reliable source classification
- **Safe metadata handling** — each `Gather()` call returns an independent metadata snapshot; deferred `Close()` errors are propagated with context
- **Security-conscious file operations** — source files are validated before destination creation to prevent data loss from TOCTOU races

## Installation and Use

Installation can be done with a normal go get:
```
$ go get github.com/conforma/go-gather
```

## Security

All efforts are made to ensure security, but gathering resources from user-provided sources has an intrinsic amount of danger. go-gather attempts to mitigate some of these issues but the user should still use caution in security-critical contexts.

## Examples

See the [`examples`](examples) directory for examples on how to use this package.