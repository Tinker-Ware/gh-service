[![Stories in Ready](https://badge.waffle.io/Tinker-Ware/gh-service.png?label=ready&title=Ready)](http://waffle.io/Tinker-Ware/gh-service)

# Github Service

This is a microservice for the infrastructure as a service environment.

## Development

This uses `godep` as to manage the project dependencies, to install it go get the package using `go get github.com/tools/godep`. Then move to the project directory and use `godep restore`, if a new dependency is added use `godep save`.

## Configuration

A configuration file must be provided, the default route for the file is located at `/etc/gh-service.conf`, the template is the next:

````yaml
  ---
  clientID: clientfromprovider
  clientSecret: secretfromprovider
  port: 1000
  salt: somesalt
  apihost: http://apihost.example
  scopes:
  - "user:email"
  - repo
````

The default path can be override using the flag `--conf`.
