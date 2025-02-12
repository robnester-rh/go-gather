# go-gather

go-gather is a library for Go (golang) for downloading ("gathering") from various sources. These sources include:
 * filepaths
 * git
 * http
 * oci

go-gather simplifies the process of gathering from these sources by freeing the implementer from having to be concerned about the details of the sources.

go-gather is based heavily on [go-getter](https://github.com/hashicorp/go-getter), but designed for specific requirements by the [ec-cli](https://github.com/enterprise-contract/ec-cli) project.

## Installation and Use

Installation can be done with a normal go get:
```
$ go get github.com/conforma/go-gather
```

## Security

All efforts are made to ensure security, but gathering resources from user provided sources has an intrensic amount of danger. go-gather attempts to mitigate some of these issues but the user should still use caution in security-critical contexts.

## Examples 

See the [`examples`](examples) directory for examples on how to use this package.