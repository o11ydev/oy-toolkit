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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCalculatePeriods(t *testing.T) {
	testCases := []struct {
		name                    string
		period                  string
		now                     time.Time
		lookbackTime            time.Duration
		includeIncompleteRanges bool
		expectedPeriods         []Period
	}{
		{
			name:                    "monthly periods, include incomplete ranges",
			period:                  "monthly",
			now:                     time.Date(2023, 2, 16, 12, 34, 56, 789, time.UTC),
			lookbackTime:            3 * 30 * 24 * time.Hour,
			includeIncompleteRanges: true,
			expectedPeriods: []Period{
				{
					Start: time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second),
				},
				{
					Start: time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second),
				},
				{
					Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second),
				},
				{
					Start: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 2, 16, 12, 34, 56, 789, time.UTC),
				},
			},
		},
		{
			name:                    "monthly periods, include incomplete ranges",
			period:                  "monthly",
			now:                     time.Date(2023, 2, 16, 12, 34, 56, 789, time.UTC),
			lookbackTime:            2 * 30 * 24 * time.Hour,
			includeIncompleteRanges: false,
			expectedPeriods: []Period{
				{
					Start: time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second),
				},
				{
					Start: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			periods := CalculatePeriods(tc.now, tc.period, tc.lookbackTime, tc.includeIncompleteRanges)

			require.Equal(t, len(tc.expectedPeriods), len(periods), "number of periods")

			for i := range periods {
				for _, p := range periods {
					t.Logf("%+v", p)
				}
				require.Equal(t, tc.expectedPeriods[i].Start, periods[i].Start, "period %d start", i)
				require.Equal(t, tc.expectedPeriods[i].End, periods[i].End, "period %d end", i)
			}
		})
	}
}
