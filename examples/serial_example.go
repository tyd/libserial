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

package examples

import (
	"sync"
	"time"

	"github.com/goiiot/libserial"
)

// ExampleSerial presents
func ExampleSerial() {
	conn, err := libserial.Open(
		// set serial device to use
		libserial.WithDevice("/dev/serial0"),
		// set baud rate
		libserial.WithBaudRate(9600),
	)

	if err != nil {
		panic("hmm, how cloud it failed")
	}

	defer conn.Close()

	rCh := make(chan byte, 1<<10)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer func() {
			close(rCh)
			wg.Done()
		}()

		for {
			buf := make([]byte, 1)
			_, err := conn.Read(buf[:])
			if err != nil {
				conn.Close()
				return
			}

			rCh <- buf[0]
		}
	}()

	wg.Add(1)
	go func() {
		t := time.NewTicker(time.Second)
		defer func() {
			t.Stop()
			wg.Done()
		}()

		for {
			select {
			case data, more := <-rCh:
				if !more {
					return
				}

				conn.Write([]byte{data})
			case <-t.C:
				println("one second passed")
			}
		}
	}()

	wg.Wait()
}
