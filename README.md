# [HE]avy [ME]tals Raspberry [PI] - Hemepi

Raspberry Pi + InkyPhat epaper display the heavy metals Gold and Silver daily prices
Built using [periph](https://periph.io/)

## Prerequisits

### Hardware

1. [Pi Zero](https://shop.pimoroni.com/products/raspberry-pi-zero-w) (preferably the wireless version)
2. [InkyPhat display](https://shop.pimoroni.com/products/inky-phat?variant=12549254217811)

### Hardware setup

SPI must be enabled on the Raspberry PI:

```
sudo raspi-config nonint do_spi 0
```

You MUST reboot after setting this option

### API key

Visit and register with <XXX> and get an API key

## Development

It is expected that you will be developing on your host machine, not the PI
directly. This means you have nothing to setup on the PI appart from the
hardware and hardware configurations as described in the prerequisits.
## Testing

```
make test
```

## Building

```
make build
```

## Deploying

The Makefile has a command for deploying, modify the target device user and host as required
```
make deploy
```

You can also copy the binary to a suitable location on your RPI, e.g:

```
scp hemepi pi@raspberrypi.local:/usr/local/bin
```

## Running

Execute the binary, passing in at least the mandatory api key flag as follows:

```
sudo hemepi --api=<api.key.here>
```

You have various configuration flags to alter the connection to the display or
the frequency data is collected from the external API

## Licences

Periph is carries the following [license](https://github.com/google/periph/blob/master/LICENSE)
