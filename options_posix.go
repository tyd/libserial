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

	"golang.org/x/sys/unix"
)

// WithDataBits set the data bits for SerialPort
// available values are {5, 6, 7, 8, 9}
// default is 8
func WithDataBits(d int) Option {
	return func(c *SerialPort) error {
		// clear flags
		c.controlOptions &= ^uint64(unix.CS5 | unix.CS6 | unix.CS7 | unix.CS8)

		switch d {
		case 5:
			c.controlOptions |= unix.CS5
		case 6:
			c.controlOptions |= unix.CS6
		case 7:
			c.controlOptions |= unix.CS7
		case 8:
			c.controlOptions |= unix.CS8
		default:
			return fmt.Errorf("invalid data bits: %v", d)
		}
		return nil
	}
}

// WithParity set parity mode
// available values are {ParityNone, ParityOdd, ParityEven}
// default is ParityNone
func WithParity(p Parity) Option {
	return func(c *SerialPort) error {
		// clear flags
		c.controlOptions &= ^uint64(unix.PARODD | unix.PARMRK | unix.PARENB)

		switch p {
		case ParityNone:
			// do nothing default is None
		case ParityOdd, ParityEven, ParityMark, ParitySpace:
			c.controlOptions |= uint64(p) | unix.PARENB
		default:
			return fmt.Errorf("invalid parity mode: %v", p)
		}

		return nil
	}
}

// WithStopBits set stop bits for SerialPort port
// available values are {StopBitOne, StopBitOneHalf, StopBitTwo}
// default is StopBitOne
func WithStopBits(s StopBit) Option {
	return func(c *SerialPort) error {
		// clear flags
		c.controlOptions &= ^uint64(unix.CSTOPB)

		switch s {
		case StopBitOne:
			// do nothing
		case StopBitTwo:
			c.controlOptions |= uint64(s)
		default:
			return fmt.Errorf("invalid stop bits: %v", s)
		}

		return nil
	}
}

// WithSoftwareFlowControl enable software flow control
func WithSoftwareFlowControl(enable bool) Option {
	return func(c *SerialPort) error {
		// clear flags
		c.inputOptions &= ^uint64(unix.IXON | unix.IXOFF | unix.IXANY)

		if enable {
			c.inputOptions |= uint64(unix.IXON | unix.IXOFF | unix.IXANY)
		}

		return nil
	}
}

// WithHardwareFlowControl enable hardware flow control
func WithHardwareFlowControl(enable bool) Option {
	return func(c *SerialPort) error {
		// clear flags
		c.controlOptions &= ^uint64(unix.CRTSCTS)

		if enable {
			c.controlOptions |= uint64(unix.CRTSCTS)
		}

		return nil
	}
}
