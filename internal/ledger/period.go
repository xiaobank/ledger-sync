package ledger

import (
	"fmt"
	"time"
)

// PeriodType defines the granularity of a time period.
type PeriodType string

const (
	PeriodDaily   PeriodType = "daily"
	PeriodWeekly  PeriodType = "weekly"
	PeriodMonthly PeriodType = "monthly"
	PeriodYearly  PeriodType = "yearly"
)

// Period represents a named time window with a start and end.
type Period struct {
	Label string
	Start time.Time
	End   time.Time
}

// Contains reports whether t falls within the period [Start, End).
func (p Period) Contains(t time.Time) bool {
	return !t.Before(p.Start) && t.Before(p.End)
}

// String returns a human-readable representation of the period.
func (p Period) String() string {
	return fmt.Sprintf("%s [%s, %s)", p.Label,
		p.Start.Format(time.DateOnly),
		p.End.Format(time.DateOnly),
	)
}

// GeneratePeriods produces n consecutive periods of the given type
// starting from the period that contains 'from'.
func GeneratePeriods(from time.Time, kind PeriodType, n int) ([]Period, error) {
	if n <= 0 {
		return nil, fmt.Errorf("period: n must be positive, got %d", n)
	}

	start := periodStart(from, kind)
	periods := make([]Period, 0, n)

	for i := 0; i < n; i++ {
		end := periodEnd(start, kind)
		periods = append(periods, Period{
			Label: periodLabel(start, kind),
			Start: start,
			End:   end,
		})
		start = end
	}
	return periods, nil
}

func periodStart(t time.Time, kind PeriodType) time.Time {
	switch kind {
	case PeriodDaily:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case PeriodWeekly:
		weekday := int(t.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		day := t.AddDate(0, 0, -(weekday - 1))
		return time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, t.Location())
	case PeriodMonthly:
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	case PeriodYearly:
		return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
	default:
		return t
	}
}

func periodEnd(start time.Time, kind PeriodType) time.Time {
	switch kind {
	case PeriodDaily:
		return start.AddDate(0, 0, 1)
	case PeriodWeekly:
		return start.AddDate(0, 0, 7)
	case PeriodMonthly:
		return start.AddDate(0, 1, 0)
	case PeriodYearly:
		return start.AddDate(1, 0, 0)
	default:
		return start.AddDate(0, 0, 1)
	}
}

func periodLabel(start time.Time, kind PeriodType) string {
	switch kind {
	case PeriodDaily:
		return start.Format("2006-01-02")
	case PeriodWeekly:
		y, w := start.ISOWeek()
		return fmt.Sprintf("%d-W%02d", y, w)
	case PeriodMonthly:
		return start.Format("2006-01")
	case PeriodYearly:
		return start.Format("2006")
	default:
		return start.Format(time.DateOnly)
	}
}
