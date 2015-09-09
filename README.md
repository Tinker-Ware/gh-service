# Github Service

This is a microservice for the infrastructure as a service environment.

## Development

This uses `godep` as to manage the project dependencies, to install it go get the package using `go get github.com/tools/godep`. Then move to the project directory and use `godep restore`, if a new dependency is added use `godep save`.

## Configuration

A configuration file must be provided, the default route for the file is located at `/etc/gh-service.conf`, the template is the next:

````yaml
  ---
  clientID: asdfg
  clientSecret: aoihcou
  port: 1000
  scopes:
  - "user:email"
  - repo
````

The default path can be override using the flag `--conf`.
