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
	"flag"
	"fmt"
	"io"
	"testing"
	"time"
)

var (
	inputPty    string
	outputPty   string
	readTimeout time.Duration
	testRWData  = []byte("goiiot/libserial")
	baseOptions = []Option{
		WithDataBits(8),
		WithParity(ParityNone),
		WithHardwareFlowControl(false),
		WithSoftwareFlowControl(false),
	}
)

func init() {
	flag.StringVar(&inputPty, "i", "", "input pty file path")
	flag.StringVar(&outputPty, "o", "", "output pty file path")
	flag.DurationVar(&readTimeout, "t", 0, "read timeout")
	baud := flag.Int("b", 0, "")

	flag.Parse()

	if inputPty == "" {
		panic("input pty is nil")
	}

	if outputPty == "" {
		panic("output pty is nil")
	}

	baseOptions = append(baseOptions, WithBaudRate(*baud))
}

func getSerialPort(options []Option) (reader, writer *SerialPort) {
	w, err := Open(inputPty, options...)
	if err != nil {
		panic(fmt.Sprintf("fatal err: %v", err))
	}

	r, err := Open(outputPty, options...)
	if err != nil {
		panic(fmt.Sprintf("fatal err: %v", err))
	}

	return r, w
}

func TestSerialPort_Readtimeout(t *testing.T) {
	options := append([]Option{WithReadTimeout(2 * time.Second)}, baseOptions...)
	r, w := getSerialPort(options)
	defer func() {
		r.Close()
		w.Close()
	}()

	start := time.Now()
	i, err := r.Read(make([]byte, 128))
	if err != nil && err != io.EOF || i != 0 {
		t.Errorf("read timeout failed: err = %v, i = %v", err, i)
	}

	if duration := time.Now().Sub(start); duration < time.Second {
		t.Errorf("read timeout not correct")
	}
}

func TestSerialPort_ReadWrite(t *testing.T) {
	r, w := getSerialPort(baseOptions)
	defer func() {
		r.Close()
		w.Close()
	}()

	_, err := w.Write(testRWData)
	if err != nil {
		t.Errorf("write data failed: %v", err)
	}

	time.Sleep(time.Second)

	buf := make([]byte, len(testRWData))
	if _, err = r.Read(buf); err != nil {
		t.Errorf("read data failed: %v", err)
	}

	if !bytes.Equal(testRWData, buf) {
		t.Errorf("target: %v, result: %v", testRWData, buf)
	}
}

func TestSerialPort_Flush(t *testing.T) {
	options := append([]Option{WithReadTimeout(time.Second)}, baseOptions...)
	r, w := getSerialPort(options)
	defer func() {
		r.Close()
		w.Close()
	}()

	if _, err := w.Write(testRWData); err != nil {
		t.Errorf("write to pty faild: %v", err)
	}

	time.Sleep(time.Second)

	if err := r.Flush(); err != nil {
		t.Errorf("flush port failed: %v", err)
	}

	time.Sleep(time.Second)

	buf := make([]byte, 128)
	if n, err := r.Read(buf); err == nil {
		t.Errorf("flush port failed, data still there: %v", string(buf[:n]))
	}
}
