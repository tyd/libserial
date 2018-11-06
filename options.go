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
	"time"
)

// Option for serial conn options
type Option func(c *SerialPort) error

// WithReadTimeout set timeout timer for read operations
// if no read timeout set, use blocking read
func WithReadTimeout(timeout time.Duration) Option {
	return func(s *SerialPort) error {
		if timeout > 0 {
			s.readTimeout = timeout
		}
		return nil
	}
}

// WithBaudRate set serial baud rate
// default is 9600
func WithBaudRate(rate int) Option {
	return func(c *SerialPort) error {
		baudRate, ok := validBaudRates[rate]
		if !ok {
			return fmt.Errorf("invalid baud rate: %v", rate)
		}

		c.baudRate = uint64(baudRate)

		// clear baud rate flags and set new baud rate
		c.controlOptions &= ^maskBaudRate
		c.controlOptions |= uint64(baudRate)
		return nil
	}
}

// WithDataBits set the data bits for SerialPort
// available values are {5, 6, 7, 8, 9}
// default is 8
func WithDataBits(d int) Option {
	return func(c *SerialPort) error {
		c.dataBits = byte(d)

		// clear flags
		c.controlOptions &= ^uint64(dataBits5 | dataBits6 | dataBits7 | dataBits8)

		switch d {
		case 5:
			c.controlOptions |= dataBits5
		case 6:
			c.controlOptions |= dataBits6
		case 7:
			c.controlOptions |= dataBits7
		case 8:
			c.controlOptions |= dataBits8
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
		c.parity = byte(p)

		// clear flags
		c.controlOptions &= ^uint64(ParityOdd | ParityMark | parityEnable)

		switch p {
		case ParityNone:
			// do nothing default is None
		case ParityOdd, ParityEven, ParityMark, ParitySpace:
			c.controlOptions |= (uint64(p) | parityEnable)
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
		c.stopBits = byte(s)

		// clear flags
		c.controlOptions &= ^uint64(StopBitTwo)

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
		// TODO: implement software flow control in windows

		// clear flags
		c.inputOptions &= ^uint64(softwareCtrlFlag)

		if enable {
			c.inputOptions |= uint64(softwareCtrlFlag)
		}

		return nil
	}
}

// WithHardwareFlowControl enable hardware flow control
func WithHardwareFlowControl(enable bool) Option {
	return func(c *SerialPort) error {
		// TODO: implement hardware flow control in windows

		// clear flags
		c.controlOptions &= ^uint64(hardwareCtrlFlag)

		if enable {
			c.controlOptions |= uint64(hardwareCtrlFlag)
		}
		return nil
	}
}
