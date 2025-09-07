# DBCache
A lightweight database query caching service to reduce load and improve response times.

## Overview / Motivation
The idea of a DB cache is to reduce load on the database itself by data being retrieved from the cache rather the database.

The reason I made this was because I wanted to create something using gRPC, and it had been a while since I've done any development work involving databases.

## Architecture Diagram
![screenshot](docs/architecture.jpg)


## Features
- gRPC API with a REST gateway for local development/testing.
- Pluggable DB and cache backends. Developers can use any database or cache as long as they have code that satisfy the datastore interfaces.
- Metrics via Prometheus.

## Quickstart / Installation
### Usage

- config
- testing

