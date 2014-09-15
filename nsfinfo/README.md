# nsfinfo

command-line tool that gives the full header details of a [NSF file](http://wiki.nesdev.com/w/index.php/NSF).

## Install
run `go get` to download via git.

    go get github.com/gbranchaud/go-nes/nsfinfo

## Usage
    $ nsfinfo path-to.nsf

## Example
    $ nsfinfo loz.nsf
        Details of loz.nsf
          song name: The Legend of Zelda
          TODO: more fields!

## todo
* listing of all the header fields (currently, only a partial listing is produced)
* add tests
