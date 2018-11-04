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

// WithBaudRate set serial baud rate
// default is 9600
func WithBaudRate(rate int) Option {
	return func(c *SerialPort) error {
		baudRate, ok := validBaudRates[rate]
		if !ok {
			return fmt.Errorf("invalid baud rate: %v", rate)
		}

		c.baudRate = uint64(baudRate)
		return nil
	}
}

// WithReadTimeout set timeout timer for read operations
func WithReadTimeout(timeout time.Duration) Option {
	return func(s *SerialPort) error {
		if timeout > 0 {
			s.readTimeout = timeout
		}
		return nil
	}
}
