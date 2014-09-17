# nsfinfo

command-line tool that gives the full header details of a [NSF file](http://wiki.nesdev.com/w/index.php/NSF).

## Install
run `go get` to download via git.

    go get github.com/gbranchaud/nsfinfo

## Usage
    $ nsfinfo path-to.nsf

## Example
    $ nsfinfo samples/loz.nsf
    name                   : The Legend of Zelda
    artist                 : Koji Kondo
    copyright holder       : 1987 Nintendo
    total # of songs       : 8
    first song             : 1
    region                 : NTSC
    play speed (µs)        : 16666
    ----------------
    nsf version            : 1
    uses bankswitching     : true
    expansion chips in use : none
    load address           : 0x8D60
    init address           : 0xA003
    play address           : 0xA000

## License
This project is licensed under the MIT license. Consult `LICENSE` to get all the details.

## todo
* add NSFE support (http://wiki.nesdev.com/w/index.php/NSFe)
