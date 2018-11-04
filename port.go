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
	"os"
	"time"
)

// SerialPort of serial
type SerialPort struct {
	f *os.File

	// options
	dev            string
	baudRate       uint64
	readTimeout    time.Duration
	dataBits       byte
	stopBits       byte
	parity         byte
	inputOptions   int64
	controlOptions int64
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
