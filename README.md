Geodisk
---------
[![Go Report Card](https://goreportcard.com/badge/github.com/influx6/geodisk)](https://goreportcard.com/report/github.com/influx6/wordies)
[![Travis Build Status](https://travis-ci.org/influx6/geodisk.svg?branch=master)](https://travis-ci.org/influx6/geodisk#)

Sample project demonstrating distance calculation from HousingAnywhere.

## Install

```bash
go get  github.com/influx6/geodisk
```

## Usage

The project installs the `geodisk` CLI into your `$GOBIN` or `$GOPATH/bin` path, ensure to have this path exported, so binaries in there can be executed.

```bash
> geodisk
Usage: geodisk [flags] [command] 

⡿ COMMANDS:
	⠙ csv        Calculate geo distance from csv file.

	⠙ db        Calculate geo distance from a db.

⡿ HELP:
	Run [command] help

⡿ OTHERS:
	Run 'geodisk flags' to print all flags of all commands.

⡿ WARNING:
	Uses internal flag package so flags must precede command name. 
	e.g 'geodisk -cmd.flag=4 run'

```


#### Using CSV Files
To generate closest and farthest locations from HousingAnywhere through a CSV file ((example)[./static/geoData.csv]), run 

```bash
>  go run main.go csv ./static/geoData.csv 

Top 5 Locations closest to Housing Anywhere:
	LocationID: 442406 (0.333838 kilometers)
	LocationID: 285782 (0.528032 kilometers)
	LocationID: 429151 (0.648010 kilometers)
	LocationID: 512818 (0.740553 kilometers)
	LocationID: 25182 (0.821642 kilometers)

Top 5 Locations farthest to Housing Anywhere:
	LocationID: 382693 (1441.719154 kilometers)
	LocationID: 382582 (1758.080613 kilometers)
	LocationID: 381823 (1758.848270 kilometers)
	LocationID: 382013 (1810.117976 kilometers)
	LocationID: 7818 (8776.646278 kilometers)

```


See `csv help` for more info:

```bash
>  wordies csv help
```


## Vendor

Vendoring is done with [GoVendor](https://github.com/kardianos/govendor).


