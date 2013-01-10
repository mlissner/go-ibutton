package w1

import (
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

	return s.bytes[0x15]&0x01 > 0
}

// HighResolution true if the chip is in 16bit (0.0625°C) mode
func (s *Status) HighResolution() bool {

	return s.bytes[0x13]&0x02 > 0
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

// Model the device identifier
func (s *Status) Model() (model deviceId) {

	return deviceId(s.bytes[0x26])
}
