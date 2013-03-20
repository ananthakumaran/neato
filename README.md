# NEato

NES emulator

![Current Status](http://dl.dropbox.com/u/24494398/neato-development/status.png)

## Playable Games

* Super Mario Bros
* Donkey Kong
* Pacman
* Ballon Fight

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

## Usage

    $> ./neato super_mario.nes

## TODO

* proper time synchronization
* mappers
* audio
