// +build !windows

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
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	serialFileFlag = unix.O_RDWR | unix.O_NOCTTY | unix.O_NONBLOCK
)

// WithDataBits set the data bits for SerialPort
// available values are {5, 6, 7, 8, 9}
// default is 8
func WithDataBits(d int) Option {
	return func(c *SerialPort) error {
		switch d {
		case 5:
			c.controlOptions |= syscall.CS5
		case 6:
			c.controlOptions |= syscall.CS6
		case 7:
			c.controlOptions |= syscall.CS7
		case 8:
			c.controlOptions |= syscall.CS8
		default:
			return fmt.Errorf("invalid data bits: %v", d)
		}
		return nil
	}
}

type Parity int64

const (
	ParityNone  Parity = 0x000
	ParityOdd   Parity = syscall.PARODD | syscall.PARENB
	ParityEven  Parity = ^syscall.PARODD | syscall.PARENB
	ParityMark  Parity = syscall.PARMRK | syscall.PARENB
	ParitySpace Parity = ^syscall.PARMRK | syscall.PARENB
)

// WithParity set parity mode
// available values are {ParityNone, ParityOdd, ParityEven}
// default is ParityNone
func WithParity(p Parity) Option {
	return func(c *SerialPort) error {
		switch p {
		case ParityNone, ParityOdd, ParityEven, ParityMark, ParitySpace:
			c.controlOptions &= ^(syscall.PARODD | syscall.PARENB)
			c.controlOptions |= int64(p)
			return nil
		default:
			return fmt.Errorf("invalid parity mode: %v", p)
		}
	}
}

type StopBit int64

const (
	StopBitOne StopBit = ^syscall.CSTOPB
	StopBitTwo StopBit = syscall.CSTOPB
)

// WithStopBits set stop bits for SerialPort port
// available values are {StopBitOne, StopBitOneHalf, StopBitTwo}
// default is StopBitOne
func WithStopBits(s StopBit) Option {
	return func(c *SerialPort) error {
		switch s {
		case StopBitOne, StopBitTwo:
			c.controlOptions |= int64(s)
			return nil
		default:
			return fmt.Errorf("invalid stop bits: %v", s)
		}
	}
}

// WithSoftwareFlowControl enable software flow control
func WithSoftwareFlowControl(enable bool) Option {
	return func(c *SerialPort) error {
		if enable {
			c.inputOptions |= syscall.IXON | syscall.IXOFF | syscall.IXANY
		} else {
			c.inputOptions &= ^(syscall.IXON | syscall.IXOFF | syscall.IXANY)
		}
		return nil
	}
}

// WithHardwareFlowControl enable hardware flow control
func WithHardwareFlowControl(enable bool) Option {
	return func(c *SerialPort) error {
		if enable {
			c.controlOptions |= unix.CRTSCTS
		} else {
			c.controlOptions &= ^unix.CRTSCTS
		}
		return nil
	}
}
