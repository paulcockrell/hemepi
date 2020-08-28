# [HE]avy [ME]tals Raspberry [PI] - Hemepi

Precious Metals prices tracker with Raspberry Pi + InkyPhat epaper display
![](./hemepi.jpg)

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

Register with [goldapi](https://www.goldapi.io) and get an API key. They offer
a limited free API key.

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

The Makefile has a command for deploying, modify the target device user and
host as required.
```
make deploy
```

You can also copy the binary to a suitable location on your RPI, e.g:

```
scp hemepi pi@raspberrypi.local:/usr/local/bin
```

## Usage

Execute the binary, passing in at least the mandatory api key flag as follows:

```
sudo ./hemepi -metal=XAU -currency=GBP -apikey=<your.api.key>
```

You have various configuration flags to alter the connection to the display or
the frequency data is collected from the external API

## Automatically running

You can setup a cronjob so that the program will run at set times/intervals. The following example will get Gold every even minute, and Silver every odd minute. If you do this you must make sure your Goldapi.io account type has the required request limits.

To edit your cron type the following (must be run as sudo)
```
sudo crontab -e
```

Add the following for switching between Gold and Silver every minute
```
*/2 * * * * sudo $HOME/hemepi/hemepi -metal=XAG -currency=GBP -apikey=<your.api.key>
1-59/2 * * * * sudo $HOME/hemepi/hemepi -metal=XAU -currency=GBP -apikey=<your.api.key>
```

## Licences

Periph is carries the following [license](https://github.com/google/periph/blob/master/LICENSE)
