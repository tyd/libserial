// +build !windows,!linux

package libserial

import (
	"golang.org/x/sys/unix"
)

type Parity uint64

const (
	ParityNone  = Parity(0)
	ParityOdd   = Parity(unix.PARODD | unix.PARENB)
	ParityEven  = ^Parity(unix.PARODD | unix.PARENB)
	ParityMark  = Parity(unix.PARMRK | unix.PARENB)
	ParitySpace = ^Parity(unix.PARMRK | unix.PARENB)
)

type StopBit uint64

const (
	StopBitOne = ^StopBit(unix.CSTOPB)
	StopBitTwo = StopBit(unix.CSTOPB)
)

const (
	serialFileFlag = unix.O_RDWR | unix.O_NOCTTY | unix.O_NONBLOCK
)

var validBaudRates = map[int]uint64{
	0:      unix.B0, // detect baud rate automatically
	50:     unix.B50,
	75:     unix.B75,
	110:    unix.B110,
	134:    unix.B134,
	150:    unix.B150,
	200:    unix.B200,
	300:    unix.B300,
	600:    unix.B600,
	1200:   unix.B1200,
	1800:   unix.B1800,
	2400:   unix.B2400,
	4800:   unix.B4800,
	9600:   unix.B9600,
	19200:  unix.B19200,
	38400:  unix.B38400,
	57600:  unix.B57600,
	115200: unix.B115200,
	230400: unix.B230400,
}
