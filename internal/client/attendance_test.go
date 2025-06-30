package client

import (
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	testCase := []struct {
		t        time.Time
		expected string
	}{
		{
			t:        time.Date(2025, 1, 1, 13, 0, 0, 0, time.Local),
			expected: "Wed 1 Jan 2025",
		},
		{
			t:        time.Date(2025, 2, 4, 1, 0, 0, 0, time.Local),
			expected: "Tue 4 Feb 2025",
		},
		{
			t:        time.Date(2025, 3, 4, 1, 0, 0, 0, time.Local),
			expected: "Tue 4 Mar 2025",
		},
		{
			t:        time.Date(2025, 4, 1, 0, 0, 0, 0, time.Local),
			expected: "Tue 1 Apr 2025",
		},
		{
			t:        time.Date(2024, 8, 1, 1, 0, 0, 0, time.Local),
			expected: "Thu 1 Aug 2024",
		},
		{
			t:        time.Date(2024, 9, 9, 0, 0, 0, 0, time.Local),
			expected: "Mon 9 Sept 2024",
		},
		{
			t:        time.Date(2024, 10, 3, 0, 0, 0, 0, time.Local),
			expected: "Thu 3 Oct 2024",
		},
		{
			t:        time.Date(2024, 11, 5, 0, 0, 0, 0, time.Local),
			expected: "Tue 5 Nov 2024",
		},
	}
	for _, tt := range testCase {
		fd := formatDate(tt.t)
		if fd != tt.expected {
			t.Errorf("%v!=%v error", fd, tt.expected)
		}
	}
}
