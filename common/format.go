// Copyright 2025 The packetd Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

func HumanizeBit(size float64) string {
	prefix := "b"

	switch {
	case size > GB:
		size = size / GB
		prefix = "Gib"
	case size > MB:
		size = size / MB
		prefix = "Mib"
	case size > KB:
		size = size / KB
		prefix = "Kib"
	}

	return fmt.Sprintf("%.4g%s", size*8, prefix)
}

func ParseBytes(s string) (int, error) {
	s = strings.ToUpper(s)
	switch {
	case strings.HasSuffix(s, "KB"):
		f, err := strconv.ParseFloat(s[:len(s)-2], 64)
		if err != nil {
			return 0, err
		}
		return int(f * KB), nil

	case strings.HasSuffix(s, "MB"):
		f, err := strconv.ParseFloat(s[:len(s)-2], 64)
		if err != nil {
			return 0, err
		}
		return int(f * MB), nil

	case strings.HasSuffix(s, "B"):
		f, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return 0, err
		}
		return int(f), err

	default:
		return 0, errors.New("unknown unit")
	}
}

func ShouldLog(total, idx int) bool {
	return total <= 10 || (idx+1)%(total/10) == 0 || idx == total-1
}
