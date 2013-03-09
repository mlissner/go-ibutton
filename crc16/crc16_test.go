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

package crc16

import "testing"

func TestChecksum(t *testing.T) {
	var in, out = []byte("123456789"), uint16(0xbb3d)
	if x := Checksum(in); x != out {
		t.Errorf("Checksum(%v) = %v, want %v", in, x, out)
		}
}
