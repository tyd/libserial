// +build freebsd

package libserial

import (
	"golang.org/x/sys/unix"
)

type Parity uint32
type StopBit uint32
type termiosFlagType = uint32
type termiosSpeedType = uint32

const (
	termiosReqGet = uint(unix.TIOCGETA)
	termiosReqSet = uint(unix.TIOCSETA)
	maskBaudRate  = uint64(0)
	ParityMark    = Parity(0)
	ParitySpace   = Parity(0)
)

func mkFlushFunc(fd uintptr) func() error {
	return func() error {
		tty, err := unix.IoctlGetTermios(int(fd), termiosReqGet)
		if err != nil {
			return err
		}

		// set serial port again for input/output flush
		return unix.IoctlSetTermios(int(fd), unix.TIOCSETAF, tty)
	}
}

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
	460800: unix.B460800,
	921600:  unix.B921600,
}
