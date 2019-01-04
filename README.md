blinktd
=======
[![microbadger][1]][2] [![docker hub][3]][4]

[1]: https://images.microbadger.com/badges/image/brimstone/blinktd.svg
[2]: https://microbadger.com/images/brimstone/blinktd
[3]: https://img.shields.io/docker/automated/brimstone/blinktd.svg
[4]: https://hub.docker.com/r/brimstone/blinktd

Simple network server for running a [blinkt!](https://shop.pimoroni.com/products/blinkt)

Made possible by [@alexellis](https://github.com/alexellis/blinkt_go_examples)

Usage
-----
```
$ docker run -d --restart=always -v /sys:/sys --name blinkt brimstone/blinktd
```
