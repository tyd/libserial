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
	"os"

	"golang.org/x/sys/unix"
)

func (s *SerialPort) open() error {
	var (
		f   *os.File
		err error
	)

	f, err = os.OpenFile(s.dev, serialFileFlag, 0666)
	if err != nil {
		return err
	}

	// check sys baud rate when baud rate not present
	if s.baudRate == unix.B0 {
		if s.baudRate = uint64(sysReadBaudRate(f.Fd())); s.baudRate == unix.B0 {
			return fmt.Errorf("can't determine serial port baud rate")
		}
	}

	if err = s.sysOpen(f); err != nil {
		f.Close()
		return err
	}

	s.f = f
	return nil
}
