package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"

	lib "github.com/goiiot/libserial"
)

var config = struct {
	device     string
	baudRate   int
	parityMode string
	dataBits   int
	stopBits   string
	swCtrl     bool
	hwCtrl     bool
}{}

func main() {
	flag.StringVar(&config.device, "dev", "", "serial device name(path)")
	flag.IntVar(&config.baudRate, "b", 9600, "baud rate")
	flag.IntVar(&config.dataBits, "d", 8, "data bits, one of: 5, 6, 7, 8")
	flag.StringVar(&config.stopBits, "s", "1", "stop bits, one of: 1, 2, 1.5 (windows only)")
	flag.StringVar(&config.parityMode, "p", "0", "parity mode, one of: none, odd, even, mark, space")
	flag.BoolVar(&config.swCtrl, "cs", false, "enable software flow control")
	flag.BoolVar(&config.hwCtrl, "ch", false, "enable hardware flow control")

	flag.Parse()
	fmt.Printf("config: %v\n", config)
	parityMode := func() lib.Parity {
		switch config.parityMode {
		case "none":
			return lib.ParityNone
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
	}()
	stopBits := func() lib.StopBit {
		switch config.stopBits {
		case "1":
			return lib.StopBitOne
		case "1.5":
			// windows
			return 1
		case "2":
			return lib.StopBitTwo
		default:
			return lib.StopBitOne
		}
	}()

	s, err := lib.Open(config.device,
		lib.WithBaudRate(config.baudRate),
		lib.WithDataBits(config.dataBits),
		lib.WithParity(parityMode),
		lib.WithStopBits(stopBits),
		lib.WithSoftwareFlowControl(config.swCtrl),
		lib.WithHardwareFlowControl(config.hwCtrl),
	)

	if err != nil {
		fmt.Printf("open serial port failed: %v", err)
		os.Exit(1)
	}
	defer s.Close()

	ctx, exit := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	defer wg.Wait()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-sigCh
		s.Close()
		exit()
	}()

	go func() {
		// print serial output
		buf := make([]byte, 128)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := s.Read(buf[:])
				if err != nil && err != io.EOF {
					fmt.Printf("read from serial error: %v", err)
					return
				}

				if n > 0 {
					print(string(buf[:n]))
				}
			}
		}
	}()

	wCh := make(chan []byte, 10)
	go func() {
		// get user input
		sc := bufio.NewScanner(os.Stdin)
		sc.Split(bufio.ScanLines)
		for {
			select {
			case <-ctx.Done():
				close(wCh)
				return
			default:
				if sc.Scan() {
					wCh <- sc.Bytes()
				}
			}
		}
	}()

	go func() {
		// write input to serial port
		for {
			select {
			case data, more := <-wCh:
				if !more {
					return
				}
				s.Write(data)
				s.Write([]byte{'\r', '\n'})
			case <-ctx.Done():
				return
			}
		}
	}()
}
