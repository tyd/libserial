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

package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"time"

	lib "github.com/goiiot/libserial"
)

var (
	ctx   context.Context
	exit  context.CancelFunc
	s     *lib.SerialPort
	wg    = &sync.WaitGroup{}
	wCh   = make(chan []byte, 10)
	sigCh = make(chan os.Signal, 1)
)

const configFmt = `
serial config:

	baud rate: %v
	data bits: %v
	parity:    %v
	stop bits: %v
	
	flow control: 
		software: %v
		hardware: %v
	
	other settings:
		suffix:     %v
		input mode: %v

`

var (
	inputMode     string
	outputMode    string
	suffix        string
	showTimeStamp bool
	config        = struct {
		device     string
		baudRate   int
		parityMode string
		dataBits   int
		stopBits   string
		swCtrl     bool
		hwCtrl     bool
	}{}
)

func main() {
	flag.Parse()
	fmt.Printf(configFmt,
		config.baudRate, config.dataBits, config.parityMode, config.stopBits,
		config.swCtrl, config.hwCtrl,
		suffix, inputMode)

	var err error
	s, err = lib.Open(config.device,
		lib.WithBaudRate(config.baudRate),
		lib.WithDataBits(config.dataBits),
		lib.WithParity(getParityMode()),
		lib.WithStopBits(getStopBit()),
		lib.WithSoftwareFlowControl(config.swCtrl),
		lib.WithHardwareFlowControl(config.hwCtrl),
	)

	if err != nil {
		fmt.Printf("open serial port failed: %v\n", err)
		os.Exit(1)
	}
	defer s.Close()

	ctx, exit = context.WithCancel(context.Background())

	signal.Notify(sigCh, os.Interrupt)
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-sigCh:
		case <-ctx.Done():
		}

		s.Close()
		exit()
	}()

	go readUserInput()
	go startReadFromPort()
	go startWriteToPort()

	wg.Wait()
}

func readUserInput() {
	// get user input
	handler := getInputHandler()
	suffix := getSuffix()
	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanLines)
	for {
		select {
		case <-ctx.Done():
			close(wCh)
			return
		default:
			if sc.Scan() {
				wCh <- handler(append(sc.Bytes(), suffix...))
			}
		}
	}
}

func startWriteToPort() {
	// write input to serial port
	handler := getOutputHandler()
	for {
		select {
		case data, more := <-wCh:
			if !more {
				return
			}
			s.Write(handler(data))
		case <-ctx.Done():
			return
		}
	}
}

func startReadFromPort() {
	// read and print serial output
	printer := getPrinter()
	handler := getOutputHandler()
	buf := make([]byte, 128)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := s.Read(buf[:])
			if err != nil && err != io.EOF {
				fmt.Printf("read from serial error: %v", err)
				exit()
				return
			}

			if n > 0 {
				printer(handler(buf[:n]))
			}
		}
	}
}

func init() {
	flag.StringVar(&config.device, "dev", "", "serial device name(path) (required)")
	flag.IntVar(&config.baudRate, "b", 9600, "baud rate")
	flag.IntVar(&config.dataBits, "d", 8, "data bits, one of 5, 6, 7, 8")
	flag.StringVar(&config.stopBits, "s", "1", "stop bits, one of 1, 2, 1.5 (windows only)")
	flag.StringVar(&config.parityMode, "p", "none", "parity mode, one of none, odd, even, mark, space")
	flag.BoolVar(&config.swCtrl, "cs", false, "enable software flow control")
	flag.BoolVar(&config.hwCtrl, "ch", false, "enable hardware flow control")
	flag.StringVar(&suffix, "suffix", "", "suffix append to serial write")
	flag.StringVar(&inputMode, "i-mode", "normal", "input mode, one of normal, hex")
	flag.StringVar(&outputMode, "o-mode", "normal", "output mode, one of normal, hex")
	flag.BoolVar(&showTimeStamp, "show-ts", false, "show timestamp")
}

func getParityMode() lib.Parity {
	switch config.parityMode {
	case "odd":
		return lib.ParityOdd
	case "even":
		return lib.ParityEven
	case "mark":
		return lib.ParityMark
	case "space":
		return lib.ParitySpace
	default:
		return lib.ParityNone
	}
}

func getStopBit() lib.StopBit {
	switch config.stopBits {
	case "1.5":
		// windows
		return 1
	case "2":
		return lib.StopBitTwo
	default:
		return lib.StopBitOne
	}
}

func normalHandler(data []byte) []byte {
	return data
}

func getOutputHandler() func([]byte) []byte {
	switch outputMode {
	case "hex":
		return func(output []byte) []byte {
			return []byte(hex.EncodeToString(output))
		}
	default:
		return normalHandler
	}
}

func getInputHandler() func([]byte) []byte {
	switch inputMode {
	case "hex":
		return func(input []byte) []byte {
			return []byte(hex.EncodeToString([]byte(input)))
		}
	default:
		return normalHandler
	}
}

func getPrinter() func([]byte) {
	if showTimeStamp {
		return func(data []byte) {
			fmt.Printf("[%s] %s\n", time.Now().Format(time.RFC3339), string(data))
		}
	}

	return func(data []byte) {
		fmt.Print(string(data))
	}
}

var (
	chSpChMap = map[byte][]byte{
		'a': {'\a'},
		'b': {'\b'},
		'f': {'\f'},
		'n': {'\n'},
		'r': {'\r'},
		't': {'\t'},
		'v': {'\v'},
	}
)

func getSuffix() []byte {
	return regexp.MustCompile(`\\[abfnrtv]`).ReplaceAllFunc([]byte(suffix), func(d []byte) []byte {
		if len(d) != 2 {
			return d
		}

		return chSpChMap[d[1]]
	})
}
