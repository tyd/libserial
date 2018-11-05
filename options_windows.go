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
)

// WithDataBits set the data bits for SerialPort
// available values are {5, 6, 7, 8, 9}
// default is 8
func WithDataBits(d int) Option {
	return func(c *SerialPort) error {
		switch d {
		case 5, 6, 7, 8:
			c.dataBits = byte(d)
			return nil
		default:
			return fmt.Errorf("invalid data bits: %v", d)
		}
	}
}

// WithParity set parity mode
// available values are {ParityNone, ParityOdd, ParityEven}
// default is ParityNone
func WithParity(p Parity) Option {
	return func(c *SerialPort) error {
		switch p {
		case ParityNone, ParityOdd, ParityEven, ParityMark, ParitySpace:
			c.parity = byte(p)
			return nil
		default:
			return fmt.Errorf("invalid parity mode: %v", p)
		}
	}
}

// WithStopBits set stop bits for SerialPort port
// available values are {StopBitOne, StopBitOneHalf, StopBitTwo}
// default is StopBitOne
func WithStopBits(s StopBit) Option {
	return func(c *SerialPort) error {
		switch s {
		case StopBitOne, StopBitOneHalf, StopBitTwo:
			c.stopBits = byte(s)
			return nil
		default:
			return fmt.Errorf("invalid stop bits: %v", s)
		}
	}
}

// WithSoftwareFlowControl (WIP)
func WithSoftwareFlowControl(enable bool) Option {
	return func(c *SerialPort) error {
		// TODO: implement software flow control in windows
		return nil
	}
}

// WithHardwareFlowControl (WIP)
func WithHardwareFlowControl(enable bool) Option {
	return func(c *SerialPort) error {
		// TODO: implement hardware flow control in windows
		return nil
	}
}
