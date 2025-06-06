/*
 * This file is part of Chihaya.
 *
 * Chihaya is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Chihaya is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Chihaya.  If not, see <http://www.gnu.org/licenses/>.
 */

package record

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"chihaya/config"
	cdb "chihaya/database/types"
)

var (
	enabled = false // overrides default, for testing purposes only
	header  = "TorrentID,UserID,Addr,Event,Uploaded,Downloaded,Left"
	channel chan []byte
)

func filename(t time.Time) string {
	return "events/events_" + t.Format("2006-01-02T15") + ".csv"
}

var initialize = sync.OnceFunc(func() {
	if err := os.Mkdir("events", 0755); err != nil && !os.IsExist(err) {
		panic(err)
	}

	start := time.Now()
	channel = make(chan []byte)

	file, err := os.OpenFile(filename(start), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	if _, err = fmt.Fprintln(file, header); err != nil {
		panic(err)
	}

	go func() {
		for buf := range channel {
			now := time.Now()
			if now.Hour() != start.Hour() {
				start = now

				if err = file.Close(); err != nil {
					panic(err)
				}

				file, err = os.OpenFile(filename(start), os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					panic(err)
				}

				if _, err = fmt.Fprintln(file, header); err != nil {
					panic(err)
				}
			}

			if _, err = file.Write(buf); err != nil {
				panic(err)
			}
		}
	}()
})

func Record(tid, uid uint32, addr cdb.PeerAddress, event string, up, down, left uint64) {
	if enabled, _ := config.GetBool("record_announces", enabled); !enabled {
		return
	}

	if up == 0 && down == 0 {
		return
	}

	// TODO: rewrite this so it doesn't tank performance
	go func() {
		initialize()

		b := make([]byte, 0, 64)
		buf := bytes.NewBuffer(b)

		buf.WriteString(strconv.FormatUint(uint64(tid), 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatUint(uint64(uid), 10))
		buf.WriteString(",")
		buf.WriteString(addr.IPString() + ":" + strconv.FormatUint(uint64(addr.Port()), 10))
		buf.WriteString(",")
		buf.WriteString(event)
		buf.WriteString(",")
		buf.WriteString(strconv.FormatUint(up, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatUint(down, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatUint(left, 10))
		buf.WriteByte('\n')

		channel <- buf.Bytes()
	}()
}
