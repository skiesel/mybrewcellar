package models

import (
	"fmt"
	"strings"
)

func getDurationString(from, to *Date) string {
	difference := to.Date.Sub(from.Date)

	days := int(difference.Hours() / 24)
	months := int(days / (365 / 12))
	years := int(months / 12)

	days -= int(months * (365 / 12))
	months -= int(years * 12)

	ageStr := ""

	if years > 0 {
		if years > 1 {
			ageStr = fmt.Sprintf("%s %d years", ageStr, years)
		} else {
			ageStr = fmt.Sprintf("%s %d year", ageStr, years)
		}
	}
	if months > 0 {
		if months > 1 {
			ageStr = fmt.Sprintf("%s %d months", ageStr, months)
		} else {
			ageStr = fmt.Sprintf("%s %d month", ageStr, months)
		}
	}
	if days > 0 {
		if days > 1 {
			ageStr = fmt.Sprintf("%s %d days", ageStr, days)
		} else {
			ageStr = fmt.Sprintf("%s %d day", ageStr, days)
		}

	}

	if ageStr == "" {
		ageStr = "< 1 day"
	}

	return strings.TrimSpace(ageStr)
}
