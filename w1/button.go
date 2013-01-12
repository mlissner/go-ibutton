// This file is part of ibutton.
//
// Copyright (C) 2013 Max Hille <mh@lambdasoup.com>
//
// ibutton is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// at your option) any later version.
//
// ibutton is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with ibutton.  If not, see <http://www.gnu.org/licenses/>.

// Package w1 provides access to 1-Wire devices
package w1

import (
	"github.com/maxhille/go-ibutton/crc16"
	"errors"
	"os"
	"strings"
	"time"
)

// iButton command codes
const (
	WRITE_SCRATCHPAD = 0x0F
	READ_SCRATCHPAD  = 0xAA
	READ_MEMORY      = 0x69
	CLEAR_MEMORY     = 0x96
)

// device identifiers type
type deviceId int

// device identifier byte descriptors
const (
	DS2422  deviceId = 0x00
	DS1923           = 0x20
	DS1922L          = 0x40
	DS1922T          = 0x60
	DS1922E          = 0x80
)

// device specific data
var devices = map[deviceId]struct {
	name      string
	offset    float32
	supported bool
	tr1       Temperature
}{
	DS2422:  {"DS2422", 0.0, false, 0.0},
	DS1923:  {"DS1923", 0.0, false, 0.0},
	DS1922L: {"DS1922L", -41.0, true, 60.0},
	DS1922T: {"DS1922T", -1.0, true, 90.0},
	DS1922E: {"DS1922E", 0.0, false, 0.0},
}

// 1-Wire device path
const W1_DIR = "/sys/bus/w1/devices"

// Button represents an iButton
type Button struct {
	file *os.File
}

// Sample represents a mission log sample
type Sample struct {
	Time time.Time
	Temp Temperature
}

// Temperature represents a temperature
type Temperature float32

// Status returns the current iButton status
func (b *Button) Status() (status *Status, err error) {

	status = new(Status)

	status.bytes, err = b.readMemory(0x0200, 3)
	if err != nil {
		return
	}

	return
}

// Open opens this iButton's 1-Wire session
func (b *Button) Open() (err error) {

	// open devices directory
	dir, err := os.Open(W1_DIR)
	if err != nil {
		return
	}

	// get devices directory contents
	infos, err := dir.Readdir(0)
	if err != nil {
		return
	}

	// filter familty 41 (iButton) devices
	var buttonInfo os.FileInfo
	for _, info := range infos {
		if strings.Index(info.Name(), "41") == 0 {
			if buttonInfo != nil {
				return errors.New("Multiple iButtons found - ibutton only supports working with a single device.")
			}

			buttonInfo = info
		}
	}
	if buttonInfo == nil {
		return errors.New("No iButton found.")
	}

	b.file, err = os.OpenFile(W1_DIR+"/"+buttonInfo.Name()+"/rw", os.O_RDWR, 0666)

	return err
}

// Close closes this iButton's 1-Wire session
func (b *Button) Close() (err error) {

	if b.file == nil {
		return
	}

	return b.file.Close()
}

// reset send a reset command to the 1-Wire bus
func (b *Button) reset() (err error) {

	// send empty write to reset
	data := make([]byte, 0)
	_, err = b.file.Write(data)

	return err
}

// ReadLog returns the log entries for the current mission
func (b *Button) ReadLog() (samples []Sample, err error) {

	// aquire button status
	status, err := b.Status()
	if err != nil {
		return
	}

	// make array with sample count length
	samples = make([]Sample, status.SampleCount())

	// determine temperature sample size
	var sampleBytes uint32
	if status.HighResolution() {
		sampleBytes = 2
	} else {
		sampleBytes = 1
	}

	// determine page count
	byteCount := status.SampleCount() * sampleBytes
	pages := int(byteCount / 32)
	if byteCount%32 != 0 {
		pages += 1
	}

	// read pages from device memory
	bytes, err := b.readMemory(0x1000, pages)
	if err != nil {
		return
	}

	// get temperature correction factors
	A, B, C := status.correctionFactors()

	// parse temperatures
	for index := uint32(0); index < status.SampleCount(); index++ {

		samples[index].Time = status.MissionTimestamp().Add(status.SampleRate() * time.Duration(index))

		temperatureBytes := bytes[index*sampleBytes : (index+1)*sampleBytes]

		tc := status.decodeTemp(temperatureBytes)
		samples[index].Temp = tc - (A*tc*tc + B*tc + C)

	}

	return
}

// ReadMemory reads the iButton's memory starting with the given address
func (b *Button) readMemory(address uint16, pages int) (bytes []byte, err error) {

	// send the read command
	cmd := make([]byte, 11)
	cmd[0] = READ_MEMORY
	cmd[1] = byte(address)
	cmd[2] = byte(address >> 8)
	_, err = b.file.Write(cmd)
	if err != nil {
		return
	}

	// read the initial package which has special parsing
	data := make([]byte, 34)
	_, err = b.file.Read(data)
	if err != nil {
		return
	}
	initial := make([]byte, 3+32)
	copy(initial, cmd[:3])
	copy(initial[3:], data[:32])
	checksum := 0xffff ^ (uint16(data[33])<<8 + uint16(data[32]))
	if crc16.Checksum(initial) != checksum {
		err = errors.New("crc check failed in initial read")
		return
	}
	bytes = append(bytes, data[:32]...)

	// read remaining pages
	for pages--; pages > 0; pages-- {
		data := make([]byte, 34)
		_, err = b.file.Read(data)
		if err != nil {
			return
		}
		checksum := 0xffff ^ (uint16(data[33])<<8 + uint16(data[32]))
		if crc16.Checksum(data[:32]) != checksum {
			err = errors.New("crc check failed failed in subsequent read")
			return
		}
		bytes = append(bytes, data[:32]...)
	}

	// tell the device to stop sending data
	b.reset()

	return
}
