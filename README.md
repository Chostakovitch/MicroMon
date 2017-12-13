# MicroMon

MicroMon is a micro application to monitor websites performance and availability created for an internship application. It watches websites at regular intervals, compute metrics for different timeframes and report them. It also warns the user when a website availability is too low.

A convenient way to use it is with a configuration file. An almost self-explained example is provided in `mm.conf`.

To run MicroMon, just `go get github.com/chostakovitch/micromon`, `go install github.com/chostakovitch/micromon/main` and run `/path/to/binary -c /path/to/conf`.

The [Wiki section](https://github.com/Chostakovitch/MicroMon/wiki) details :
* More options for configuration file
* How to use MicroMon as a library
* How to easily extend MicroMon with new features
* Details about architecture and implementation
* What improvements could be done
