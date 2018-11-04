// +build windows

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
	"math"
	"os"
	"syscall"
	"unsafe"

	win "golang.org/x/sys/windows"
)

const (
	SetCommState    = "SetCommState"
	SetCommTimeouts = "SetCommTimeouts"
	SetCommMask     = "SetCommMask"
	SetupComm       = "SetupComm"
	PurgeComm       = "PurgeComm"
)

var (
	comSyscall     = map[string]func(s *SerialPort) error{}
	comSyscallList = []string{
		SetCommState, SetupComm, SetCommTimeouts, SetCommMask, PurgeComm,
	}
)

func (s *SerialPort) sysOpen() error {
	f, err := os.OpenFile(s.dev, os.O_RDWR, 0)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			f.Close()
			return
		}
		s.f = f
	}()

	for _, name := range comSyscallList {
		if name != PurgeComm {
			if err = comSyscall[name](s); err != nil {
				return err
			}
		}
	}

	return nil
}

func init() {
	dll, err := win.LoadLibrary("kernel32.dll")

	var sys = make(map[string]func(args ...uintptr) (uintptr, error))

	if err != nil {
		panic("init kernel32.dll failed")
	}

	for _, name := range comSyscallList {
		addr, err := win.GetProcAddress(dll, name)
		if err != nil {
			sys[name] = func(args ...uintptr) (uintptr, error) {
				n := uintptr(len(args))
				r, _, err := syscall.Syscall(addr, n, args[0], args[1], args[2])
				return r, err
			}
		}
	}

	comSyscall[PurgeComm] = func(s *SerialPort) error {
		// PURGE_TXABORT | PURGE_RXABORT | PURGE_TXCLEAR | PURGE_RXCLEAR
		r, err := sys[PurgeComm](s.f.Fd(), 0x000F)
		if r == 0 {
			return err
		}
		return nil
	}

	comSyscall[SetupComm] = func(s *SerialPort) error {
		r, err := sys[SetupComm](s.f.Fd(), 64, 64)
		if r == 0 {
			return err
		}
		return nil
	}

	comSyscall[SetCommMask] = func(s *SerialPort) error {
		// EV_RXCHAR
		r, err := sys[SetCommMask](s.f.Fd(), 0x0001)
		if r == 0 {
			return err
		}
		return nil
	}

	comSyscall[SetCommTimeouts] = func(s *SerialPort) error {
		timeout := &struct {
			ReadIntervalTimeout         uint32
			ReadTotalTimeoutMultiplier  uint32
			ReadTotalTimeoutConstant    uint32
			WriteTotalTimeoutMultiplier uint32
			WriteTotalTimeoutConstant   uint32
		}{
			ReadIntervalTimeout:        math.MaxUint32,
			ReadTotalTimeoutMultiplier: math.MaxUint32,
			ReadTotalTimeoutConstant: func() uint32 {
				if s.readTimeout == 0 {
					return math.MaxUint32 - 1
				}
				return uint32(s.readTimeout.Nanoseconds() / 1e6)
			}(),
		}

		r, err := sys[SetCommTimeouts](s.f.Fd(), uintptr(unsafe.Pointer(timeout)))
		if r == 0 {
			return err
		}
		return nil
	}

	comSyscall[SetCommState] = func(s *SerialPort) error {
		DCB := &struct {
			DCBLength, BaudRate                            uint32
			flags                                          [4]byte
			wReserved, XOnLim, XOffLim                     uint16
			ByteSize, Parity, StopBits                     byte
			XonChar, XOffChar, ErrorChar, EofChar, EvtChar byte
			wReserved1                                     uint16
		}{
			flags:    [4]byte{0x11, 0x00, 0x00, 0x00},
			BaudRate: uint32(s.baudRate),
			ByteSize: s.dataBits,
			StopBits: s.stopBits,
			Parity:   s.parity,
		}

		DCB.DCBLength = uint32(unsafe.Sizeof(*DCB))

		r, err := sys[SetCommState](s.f.Fd(), uintptr(unsafe.Pointer(DCB)))
		if r == 0 {
			return err
		}
		return nil
	}
}

var validBaudRates = map[int]uint32{
	0:       0,
	50:      50,
	75:      75,
	110:     110,
	134:     134,
	150:     150,
	200:     200,
	300:     300,
	600:     600,
	1200:    1200,
	1800:    1800,
	2400:    2400,
	4800:    4800,
	9600:    9600,
	19200:   19200,
	38400:   38400,
	57600:   57600,
	115200:  115200,
	230400:  230400,
	460800:  460800,
	500000:  500000,
	576000:  576000,
	921600:  921600,
	1000000: 1000000,
	1152000: 1152000,
	1500000: 1500000,
	2000000: 2000000,
	2500000: 2500000,
	3000000: 3000000,
	3500000: 3500000,
	4000000: 4000000,
}
