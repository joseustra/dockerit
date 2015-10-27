# Dockerit

## A tool to execute docker containers easily from config file

Run docker containers using a simple config file. Very similar to docker-compose

## Config file example

The config file should be on: ***~/.dockerit/folder_name.yml***

```yaml
container:
  name: "name"
  port: "3000:3000"
  link: "mongodb:mongo"
  image: "myimage"
  volume: "/app"
```

## Status

The project is in a very early stage, so is not ready to be used yet.

## Installation

Get the code
```
go get -u github.com/ustrajunior/dockerit
```

## Documentation

TODO

## Tests

TODO


## License

Released under the [MIT License](https://github.com/ustrajunior/dockerit/blob/master/LICENSE).
