// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// "Blinkenlights" for the CASIO LK-S250 keyboard.
package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/writer"

	driver "gitlab.com/gomidi/rtmididrv"
)

var (
	port = flag.String("port", "CASIO USB-MIDI",
		"Substring of the port name to match to find a port.")
)

func main() {
	flag.Parse()

	drv, err := driver.New()
	if err != nil {
		log.Panic(err)
	}
	defer drv.Close()

	var out midi.Out
	unique := true
	outs, err := drv.Outs()
	if err != nil {
		log.Panic(err)
	}
	for _, p := range outs {
		log.Printf("Available port: %s\n", p.String())
		if strings.Contains(p.String(), *port) {
			log.Printf("Matched %s.\n", p.String())
			if out != nil {
				unique = false
			}
			out = p
		}
	}
	if out == nil || !unique {
		log.Panicf("No unique port. See above for valid ports to choose from via the --port option.")
	}

	out.Open()
	runLights(writer.New(out))
}

func ping(wr writer.ChannelWriter) {
	writer.SysEx(wr, []byte{
		0x44, 0x7E, 0x7E, 0x7F,
		0x00, 0x03,
	})
}

func what(wr writer.ChannelWriter) {
	writer.SysEx(wr, []byte{
		0x44, 0x7E, 0x7E, 0x7F,
		0x00, 0x06, 0x00,
	})
}

func lightOn(wr writer.ChannelWriter, note int) {
	writer.SysEx(wr, []byte{
		0x44, 0x7E, 0x7E, 0x7F,
		0x02, 0x00, byte(note), 0x01,
	})
}

func lightOff(wr writer.ChannelWriter, note int) {
	writer.SysEx(wr, []byte{
		0x44, 0x7E, 0x7E, 0x7F,
		0x02, 0x00, byte(note), 0x00,
	})
}

const (
	minKey = 60 - 12*2
	maxKey = 60 + 12*3
)

func whiteKey(i int) int {
	white := []int{0, 2, 4, 5, 7, 9, 11}
	return (i/7)*12 + white[i%7]
}

func runLights(wr writer.ChannelWriter) {
	minWhiteKey, maxWhiteKey := 127, 0
	for i := 0; i <= 127; i++ {
		if whiteKey(i) >= minKey && whiteKey(i) <= maxKey {
			if minWhiteKey == 127 {
				minWhiteKey = i
			}
			maxWhiteKey = i
		}
	}

	ping(wr)
	pos := minWhiteKey
	dir := +1
	size := 4
	for {
		log.Printf("Cycle...\n")
		ping(wr)
		lightOff(wr, whiteKey(pos-size*dir))
		lightOn(wr, whiteKey(pos))
		pos += dir
		if pos > maxWhiteKey {
			dir = -1
			pos = maxWhiteKey - size
		}
		if pos < minWhiteKey {
			dir = +1
			pos = minWhiteKey + size
		}
		time.Sleep(67 * time.Millisecond)
	}
}
