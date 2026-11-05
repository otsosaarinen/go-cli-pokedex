# go-cli-pokedex

![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white)
![Dependencies](https://img.shields.io/badge/dependencies-none-brightgreen)
![Data](https://img.shields.io/badge/data-PokéAPI-ef5350)

A small command-line Pokédex written in Go. Look up one or more Pokémon and see their weight and types, using data from [PokéAPI](https://pokeapi.co).

This is a learning project for getting to know Go: HTTP clients, JSON decoding, error handling and testing, using only the standard library.

## Usage

```
go run . get <pokemon> [pokemon...]
```

```
$ go run . get pikachu arceus
Name: pikachu
Weight: 6.0 kg
Types: electric

Name: arceus
Weight: 320.0 kg
Types: normal
```

Names are case-insensitive. If a lookup fails, the error goes to stderr, the remaining names still run, and the exit code is 1.

```
$ go run . get pikachu asd
Name: pikachu
Weight: 6.0 kg
Types: electric

error: asd: pokemon not found
```

## Build

```
go build -o go-cli-pokedex .
./go-cli-pokedex get pikachu
```

## Test

```
go test ./...
```

Tests run offline against a fake local server, so they don't call the real API.

## Project layout

- `main.go`: argument parsing and output (`run`)
- `pokeapi.go`: PokéAPI client (`FetchPokemon`)
- `*_test.go`: tests
