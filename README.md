# SHORTURL

This project aims to be a cloud native service that consists on the following architecture:
* valkey -> used for caching of the existing keys
* mysql/psql/sqlite -> used for long term storage depending on the package was build using tags either of the 3 db's can be used
* shorturl -> main service running that shortens the url, specially created to be run on a kubernetes cluster or as a standalone docker container

## Development
Running `nix develop` would create a shell with the dependencies that made the development of this project possible

## Build using nix
In order to build the project there's 2 possible ways:
* `nix build .#dockerImage.<architecture> # architecture could be x86_64-linux`
* `nix build`

## Installation via helm
```(bash)
git clone git@github.com:FKouhai/shorturl.git
cd shorturl
# update the values inside of urlshort/
helm install -f urlshort/values.yaml urlshort ./urlshort
```
## OTEL TRACING
As of right now only the OTEL grpc endpoint is supported
### traces
![image](https://github.com/user-attachments/assets/96ec4ada-528e-4b28-9fd4-fed3f7da2265)
### spans
![image](https://github.com/user-attachments/assets/c6196324-fb05-4dc1-bfd1-99ef01910b49)


## TODO
* docker-compose that uses the default db's 
* github workflow to build the project
    * Ideally speaking the workflow would run a nix build command to build the binary, build the docker image(s) and build the charts
