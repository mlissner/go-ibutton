package w1

import (
	"fmt"
	"time"
)

// Status represents an iButton status. The iButton's status is saved in two register pages (0x0200-0x0263)
type Status struct {
	bytes []byte
}

// Time the time
func (s *Status) Time() time.Time {

	return parseTime(s.bytes[0x00:0x06])

}

// MissionTimestamp the current mission timestamp
func (s *Status) MissionTimestamp() time.Time {

	return parseTime(s.bytes[0x19:0x1F])

}

// decodeTemp gives the temperature encoded in the given byte slice
func (s *Status) decodeTemp(bytes []byte) (temp Temperature) {

	switch len(bytes) {
	case 1:
		temp = Temperature(float32(bytes[0])/2 + devices[s.DeviceId()].offset)
	case 2:
		temp = Temperature(float32(bytes[0])/2 + devices[s.DeviceId()].offset + float32(bytes[1])/512)
	}

	return
}

// parseTime parses a time object from the given bytes
func parseTime(bytes []byte) time.Time {

	year := int(2000) + int(bytes[5]&0x0f) + int(bytes[5]>>4)*10
	month := int(bytes[4]&0x0f) + int(bytes[4]>>4)*10
	day := int(bytes[3]&0x0f) + int(bytes[3]>>4)*10
	hour := int(bytes[2]&0x0f) + int(bytes[2]>>4)&3*10
	minute := int(bytes[1]&0x0f) + int(bytes[1]>>4)*10
	second := int(bytes[0]&0x0f) + int(bytes[0]>>4)*10

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

// SampleCount count of recorded samples since last mission start
func (s *Status) SampleCount() uint32 {

	return uint32(s.bytes[0x22])<<16 + uint32(s.bytes[0x21])<<8 + uint32(s.bytes[0x20])
}

// MissionInProgress true if a mission is running
func (s *Status) MissionInProgress() bool {

	return s.bytes[0x15]&(0x01<<1) > 0
}

// HighResolution true if the chip is in 16bit (0.0625°C) mode
func (s *Status) HighResolution() bool {

	return s.bytes[0x13]&(0x01<<2) > 0
}

// SampleRate return the currently set sample rate
func (s *Status) SampleRate() (duration time.Duration) {

	// first read in the raw rate
	rate := uint32(s.bytes[0x06]) + uint32(s.bytes[0x07])<<8

	// decide on minutes or seconds
	if s.bytes[0x12]>>1 == 1 {
		duration = time.Duration(rate) * time.Second
	} else {
		duration = time.Duration(rate) * time.Minute
	}

	return
}

// DeviceId the device identifier byte
func (s *Status) DeviceId() (model deviceId) {

	return deviceId(s.bytes[0x26])
}

// correctionFactors returns the temperature correction factors for this device
func (s *Status) correctionFactors() (a Temperature, b Temperature, c Temperature) {

	// get chip-hardcoded correction values
	tr1 := devices[s.DeviceId()].tr1
	tr2 := s.decodeTemp(s.bytes[0x40:0x42])
	tc2 := s.decodeTemp(s.bytes[0x42:0x44])
	tr3 := s.decodeTemp(s.bytes[0x44:0x46])
	tc3 := s.decodeTemp(s.bytes[0x46:0x48])

	// calculate correction factors
	err2 := tc2 - tr2
	err3 := tc3 - tr3
	err1 := err2

	// formula stuff from DS1922L data sheet (p.50)
	b = (tr2*tr2 - tr1*tr1) * (err3 - err1) / ((tr2*tr2-tr1*tr1)*(tr3-tr1) + (tr3*tr3-tr1*tr1)*(tr1-tr2))
	a = b * (tr1 - tr2) / (tr2*tr2 - tr1*tr1)
	c = err1 - a*tr1*tr1 - b*tr1

	return
}

// Name the device model's name
func (s *Status) Name() string {

	device, ok := devices[s.DeviceId()]
	if ok {
		return device.name
	}

	return fmt.Sprintf("Unknown Device (deviceId:%x)", s.DeviceId())
}
