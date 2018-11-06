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
	"errors"
	"os"
	"time"
)

var (
	// ErrDeviceNameEmpty happens when opening a device with empty name
	ErrDeviceNameEmpty = errors.New("device name should not be empty")
)

// SerialPort of serial
type SerialPort struct {
	// actions performer
	f     *os.File
	flush func()

	// options
	// common options
	dev         string
	readTimeout time.Duration

	// options for windows
	baudRate uint64
	dataBits byte
	stopBits byte
	parity   byte

	// options for posix/linux
	inputOptions   uint64
	controlOptions uint64
}

// Write bytes to serial connection
func (s *SerialPort) Write(data []byte) (int, error) {
	return s.f.Write(data)
}

// Read bytes from serial connection
func (s *SerialPort) Read(data []byte) (int, error) {
	return s.f.Read(data)
}

// Close serial connection
func (s *SerialPort) Close() error {
	return s.f.Close()
}

// Flush serial input/output queue
func (s *SerialPort) Flush() {
	s.flush()
}

// Open serial port
func Open(device string, options ...Option) (*SerialPort, error) {
	if device == "" {
		return nil, ErrDeviceNameEmpty
	}

	port := &SerialPort{dev: device}

	// set defaults 9600 8N1
	defaultOptions := []Option{
		WithBaudRate(9600),
		WithDataBits(8),
		WithParity(ParityNone),
		WithStopBits(StopBitOne),
	}

	for _, setDefaultOption := range defaultOptions {
		if err := setDefaultOption(port); err != nil {
			return nil, err
		}
	}

	// set user defined options
	for _, setOption := range options {
		if err := setOption(port); err != nil {
			return nil, err
		}
	}

	// open platform specific serial port
	if err := port.open(); err != nil {
		return nil, err
	}

	return port, nil
}
