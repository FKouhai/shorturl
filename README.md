# SHORTURL

This project aims to be a cloud native service that consists on the following architecture:
* valkey -> used for caching of the existing keys
* mysql/psql/sqlite -> used for long term storage depending on the package was build using tags either of the 3 db's can be used
* shorturl -> main service running that shortens the url, specially created to be run on a kubernetes cluster or as a standalone docker container

## Development
Running `nix develop` would create a shell with the dependencies that made the development of this project possible

## TODO
* Include in the nix flake the build of the container image for the main service
* helm chart for the project that provides:
    * valkey config
    * mysql config
    * psql config
    * sqlite config
* docker-compose that uses the default db's 
* Instrumentation using the OTEL sdk
* github workflow to build the project
    * Ideally speaking the workflow would run a nix build command to build the binary, build the docker image(s) and build the charts
