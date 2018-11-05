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

// Open serial port
func Open(options ...Option) (*SerialPort, error) {
	// set defaults 9600 8N1
	port := &SerialPort{
		baudRate:       9600,
		dataBits:       8,
		parity:         byte(ParityNone),
		stopBits:       1,
		controlOptions: validBaudRates[9600] | uint64(ParityNone) | uint64(StopBitOne),
	}

	for _, setOption := range options {
		if err := setOption(port); err != nil {
			return nil, err
		}
	}

	if err := port.open(); err != nil {
		return nil, err
	}

	return port, nil
}
