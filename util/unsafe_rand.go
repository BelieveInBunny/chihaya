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

package util

import (
	unsafeRandom "math/rand"
	"time"
)

var randomSource = unsafeRandom.New(unsafeRandom.NewSource(time.Now().Unix()))

func UnsafeInt() int {
	return randomSource.Int()
}

func UnsafeIntn(n int) int {
	return randomSource.Intn(n)
}

func UnsafeUint32() uint32 {
	return randomSource.Uint32()
}

func UnsafeUint64() uint64 {
	return randomSource.Uint64()
}

func UnsafeRand(min int, max int) int {
	return randomSource.Intn(max-min+1) + min
}

func UnsafeReadRand(p []byte) (n int, err error) {
	return randomSource.Read(p)
}
