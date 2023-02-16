// Copyright The o11y toolkit Authors
// spdx-license-identifier: apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package period

import (
	"time"
)

type Period struct {
	Start    time.Time
	End      time.Time
	Complete bool
}

func CalculatePeriods(now time.Time, period string, lookbackTime time.Duration, includeIncompleteRanges bool) []Period {
	periods := []Period{}

	// Calculate the start and end times for each period based
	// on the specified period string
	switch period {
	case "monthly":
		firstDate := monthStart(now.Add(-lookbackTime), false)
		currentMonth := firstDate
		for d := 0; currentMonth.Before(now); d++ {
			complete := true
			end := monthStart(currentMonth, true).Add(-1 * time.Second)
			if end.After(now) {
				if !includeIncompleteRanges {
					break
				}
				complete = false
				end = now
			}

			periods = append(periods, Period{
				Start:    currentMonth,
				End:      end,
				Complete: complete,
			})

			currentMonth = monthStart(currentMonth, true)
		}
	}
	return periods
}

func monthStart(t time.Time, next bool) time.Time {
	year := t.Year()
	month := int(t.Month())
	if next {
		if month == 12 {
			month = 1
			year++
		} else {
			month++
		}
	}
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, t.Location())
}
