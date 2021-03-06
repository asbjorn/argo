Argo
====

Argo is a [Signal K](http://signalk.org) compliant server written in
[Go](http://golang.org). The goal of Argo is to consume as many different
sources of data as possible and convert them to Signal K on as many different
transports as possible.

It currently supports [WebSockets](http://www.w3.org/TR/websockets/),
[ZeroMQ](http://zeromq.org/) and [MQTT](http://mqtt.org) transports and
ingestion of [NMEA 2000](http://en.wikipedia.org/wiki/NMEA_2000) data using an
Actisense NGT-1 NMEA 2000 to USB converter or a Lawicel CAN-USB adapter. It can
also read [CANboat](https://github.com/canboat/canboat) JSON from a file.

Installation
------------

If you don't have go, you'll need to install it first. See the golang [Getting
Started](http://golang.org/doc/install) guide for information. On Debian/Ubuntu
systems `apt-get install golang` may be sufficient.

After that, you can install argo. Go has some opinions on how you should manage
your source code workspaces, so use the commands below to get the paths right.

```
$ sudo apt-get install libzmq-dev
$ cd $GOPATH/src
$ git clone git@github.com:timmathews/argo github.com/timmathews/argo
$ cd github.com/timmathews/argo
$ ./build.sh
$ sudo ./build.sh --install
```

If you're not interested in hacking on Argo, you can install the libzmq-dev
library as above and then `go get github.com/timmathews/argo`. Which will
install argo in into $GOPATH/bin.

By default, argo will try to use /dev/ttyUSB0 for the Actisense NGT-1, but you
can change this by calling argo with a specific device like `argo
/dev/ttyUSB3`.

There are utility scripts in the `canusb` and `actisense` folders which will
install udev rules for those devices. You probably won't need these, but if for
some reason your distro doesn't recognize the vendor IDs for those devices,
these scripts can help.

TODO
----

* Daemonize argo and write an upstart script for it
* Better Signal K support
* Simplify commandline arguments and move more settings into the config file
* Set the default path for the config file to /etc
* Add support for reading raw CAN packet captures in Actisense and CANUSB
  formats (maybe others)
* Add SocketCAN support
* Web GUI for configuration
* CLI for configuration
