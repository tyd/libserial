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
	"bytes"
	"testing"
	"time"
)

var (
	testRWData = []byte("goiiot/libserial")
)

func TestReadWrite(t *testing.T) {
	r, w := getSerialPort()

	_, err := w.Write(testRWData)
	if err != nil {
		t.Errorf("write data failed: %v", err)
	}

	time.Sleep(time.Second)

	buf := make([]byte, len(testRWData))
	_, err = r.Read(buf)
	if err != nil {
		t.Errorf("read data failed: %v", err)
	}

	if !bytes.Equal(testRWData, buf) {
		t.Errorf("target: %v, result: %v", testRWData, buf)
	}
}
