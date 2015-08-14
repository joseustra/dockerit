# Gozek

## A tool to execute docker containers easily from config file

Run docker containers using a simple config file. Very similar to docker-compose

## Config file example

```yaml
container:
  name: "name"
  port: "3000:3000"
  link: "mongodb:mongo"
  image: "myimage"
```

## Status

The project is in a very early stage, so is not ready to be used yet.

## Installation

Get the code
```
go get -u github.com/ustrajunior/gozek
```

## Documentation

TODO

## Tests

TODO


## License

Released under the [MIT License](https://github.com/ustrajunior/gozek/blob/master/LICENSE).
