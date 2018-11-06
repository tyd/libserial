// +build windows

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
	"unsafe"
)

type Parity byte

const (
	ParityNone  Parity = 0
	ParityOdd   Parity = 1
	ParityEven  Parity = 2
	ParityMark  Parity = 3
	ParitySpace Parity = 4
)

type StopBit byte

const (
	StopBitOne     StopBit = 0
	StopBitOneHalf StopBit = 1
	StopBitTwo     StopBit = 2
)

type _dcb struct {
	DCBLength, BaudRate                            uint32
	flags                                          [4]byte
	wReserved, XOnLim, XOffLim                     uint16
	ByteSize, Parity, StopBits                     byte
	XonChar, XOffChar, ErrorChar, EofChar, EvtChar byte
	wReserved1                                     uint16
}

const _dcbSize = uint32(unsafe.Sizeof(_dcb{}))

var validBaudRates = map[int]uint32{
	0:       0,
	50:      50,
	75:      75,
	110:     110,
	134:     134,
	150:     150,
	200:     200,
	300:     300,
	600:     600,
	1200:    1200,
	1800:    1800,
	2400:    2400,
	4800:    4800,
	9600:    9600,
	19200:   19200,
	38400:   38400,
	57600:   57600,
	115200:  115200,
	230400:  230400,
	460800:  460800,
	500000:  500000,
	576000:  576000,
	921600:  921600,
	1000000: 1000000,
	1152000: 1152000,
	1500000: 1500000,
	2000000: 2000000,
	2500000: 2500000,
	3000000: 3000000,
	3500000: 3500000,
	4000000: 4000000,
}
