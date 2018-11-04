// +build !linux,!windows

/*
 * Copyright Go-IIoT (https://github.com/goiiot)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package libserial

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

func sysReadBaudRate(fd uintptr) uint64 {
	tty := &unix.Termios{}
	if _, _, err := unix.Syscall(
		unix.SYS_IOCTL, fd, unix.TIOCGETA, uintptr(unsafe.Pointer(tty)),
	); err != 0 {
		return 0
	}
	return tty.Cflag & 0x100f
}

// open serial connection
func (s *SerialPort) sysOpen(f *os.File, timeout uint8) error {
	tty := &unix.Termios{
		Cflag: unix.CREAD | unix.CLOCAL | s.baudRate |
			uint64(s.dataBits) | uint64(s.stopBits) | uint64(s.parity),
		Iflag:  unix.IGNBRK,
		Ispeed: uint64(s.baudRate),
		Ospeed: uint64(s.baudRate),
		// Lflag:  0,
		// Oflag:  0,
	}

	// non block for read
	tty.Cc[unix.VMIN] = 0
	tty.Cc[unix.VTIME] = timeout

	if _, _, err := unix.Syscall(
		unix.SYS_IOCTL, f.Fd(), unix.TIOCSETA, uintptr(unsafe.Pointer(tty)),
	); err != 0 {
		return fmt.Errorf(err.Error())
	}

	if err := unix.SetNonblock(int(f.Fd()), false); err != nil {
		return err
	}

	return nil
}

var validBaudRates = map[int]uint32{
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
