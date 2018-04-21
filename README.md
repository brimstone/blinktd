blinktd
=======

Simple network server for running a [blinkt!](https://shop.pimoroni.com/products/blinkt)

Made possible by [@alexellis](https://github.com/alexellis/blinkt_go_examples)

Usage
-----
```
$ docker run -d --restart=always -v /sys:/sys --name blinkt brimstone/blinktd
```
