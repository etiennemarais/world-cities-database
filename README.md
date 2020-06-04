# World Cities Database

![Go](https://github.com/etiennemarais/world-cities-database/workflows/Go/badge.svg?branch=master)

A cli command to visit the [worldcitiesdb.com/country/list](http://www.worldcitiesdb.com/country/list) and compile the country/regional data into usable sql migrations for projects that are already normalised.

## Installation

```sh
go get github.com/etiennemarais/world-cities-database
```

## Usage

```sh
go run main.go

Usage:
  world-cities-database [command]

Available Commands:
  generate    Generate a resource in the specified format
  help        Help about any command
  list        List the specified resource

Flags:
  -h, --help   help for world-cities-database

Use "world-cities-database [command] --help" for more information about a command.
```