package cmd

import (
	"reflect"
	"testing"
	"time"

	"cl-parse/changelog"
)

func TestFilterEntries(t *testing.T) {
	base := time.Date(2025, 9, 8, 15, 4, 5, 0, time.UTC)

	mkEntry := func(ver string, daysAgo int) changelog.ChangelogEntry {
		return changelog.ChangelogEntry{Version: ver, Date: base.AddDate(0, 0, -daysAgo)}
	}

	entries := []changelog.ChangelogEntry{
		mkEntry("v5", 0),  // today (newest)
		mkEntry("v4", 1),  // 1 day ago
		mkEntry("v3", 3),  // 3 days ago
		mkEntry("v2", 7),  // 7 days ago
		mkEntry("v1", 10), // 10 days ago (oldest)
	}

	tests := []struct {
		name      string
		last      int
		sinceDays int
		wantVers  []string
	}{
		{"no filters", 0, 0, []string{"v5", "v4", "v3", "v2", "v1"}},
		{"last only trims", 2, 0, []string{"v5", "v4"}},
		{"last > len", 10, 0, []string{"v5", "v4", "v3", "v2", "v1"}},
		{"sinceDays only inclusive", 3, 3, []string{"v5", "v4", "v3"}}, // cutoff 3 days ago includes v3
		{"sinceDays only none", 0, 0, []string{"v5", "v4", "v3", "v2", "v1"}},
		{"sinceDays window bigger", 0, 8, []string{"v5", "v4", "v3", "v2"}},
		{"last applied before sinceDays", 3, 7, []string{"v5", "v4", "v3"}}, // last=3 gives v5,v4,v3; all within 7 days
	}

	for _, tt := range tests {
		got := filterEntries(entries, tt.last, tt.sinceDays, base)
		var gotVers []string
		for _, e := range got {
			gotVers = append(gotVers, e.Version)
		}
		if !reflect.DeepEqual(gotVers, tt.wantVers) {
			var allDates []string
			for _, e := range got {
				allDates = append(allDates, e.Date.Format("2006-01-02"))
			}
			t.Fatalf("%s: expected %v, got %v (dates %v)", tt.name, tt.wantVers, gotVers, allDates)
		}
	}
}
