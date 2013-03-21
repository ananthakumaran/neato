# NEato [![Build Status](https://secure.travis-ci.org/ananthakumaran/neato.png)](http://travis-ci.org/ananthakumaran/neato)

NES emulator

![Current Status](http://dl.dropbox.com/u/24494398/neato-development/status.png)

## Playable Games

* Super Mario Bros
* Donkey Kong
* Pacman
* Ballon Fight
* Contra

## Keys

* UP - UP ARROW
* DOWN - DOWN ARROW
* LEFT - LEFT ARROW
* RIGHT - RIGHT ARROW
* START - SPACE
* SELECT - ENTER
* A - a
* B - s

## Installation

### Mac

    $> cd path/to/neato
    $> export GOPATH=`pwd`
    $> brew install glew glfw
    $> cd src/neato
    $> go get -v
    $> go build -v

### Ubuntu

    $> cd path/to/neato
    $> export GOPATH=`pwd`
    $> sudo apt-get install libglfw-dev libglew1.6-dev libxrandr-dev
    $> cd src/neato
    $> go get -v
    $> go build -v

## Usage

    $> ./neato super_mario.nes

## TODO

* proper time synchronization
* mappers
* audio
