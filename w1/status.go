package w1

import (
	"time"
)

// Status represents an iButton status. The iButton's status is saved in two register pages (0x0200-0x0263)
type Status struct {
	bytes []byte
}

// Date the date
func (s *Status) Time() (t time.Time) {

	year   := int(2000) + int(s.bytes[5] & 0x0f) + int(s.bytes[5] >> 4) * 10
	month  := int(s.bytes[4] & 0x0f) + int(s.bytes[4] >> 4) * 10
	day    := int(s.bytes[3] & 0x0f) + int(s.bytes[3] >> 4) * 10
	hour   := int(s.bytes[2] & 0x0f) + int(s.bytes[2] >> 4) & 3 * 10
	minute := int(s.bytes[1] & 0x0f) + int(s.bytes[1] >> 4) * 10
	second := int(s.bytes[0] & 0x0f) + int(s.bytes[0] >> 4) * 10

	t = time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)

	return
}
