# libserial cmd

Demo command line app using [goiiot/libserial](https://github.com/goiiot/libserial)

## Usage

example command line usage

- use device `/dev/ttyUSB0`
- baud rate `115200`
- show data read in hex format
- append suffix `\r\n` when writing to serial port
- show output timestamp

```bash
libserial \
    -dev /dev/ttyUSB0 \
    -b 115200 \
    -o-mode hex \
    -suffix '\r\n' \
    -show-ts
```
