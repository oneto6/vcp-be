// Package meet handle meet db
package meet

import "time"

type Meet struct {
	Time        time.Time `json:"time"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func GetMeet() (list []Meet) {
	today := time.Now()

	list = append(list,
		Meet{
			time.Date(today.Year(), today.Month(), today.Day(), 9, 0, 0, 0, today.Location()),
			"Daily Standup",
			"Morning sync meeting",
		},
		Meet{
			time.Date(today.Year(), today.Month(), today.Day(), 11, 30, 0, 0, today.Location()),
			"Client Call",
			"Discussion with client team",
		},
		Meet{
			time.Date(today.Year(), today.Month(), today.Day(), 14, 0, 0, 0, today.Location()),
			"Backend Review",
			"Review API architecture",
		},
		Meet{
			time.Date(today.Year(), today.Month(), today.Day(), 16, 15, 0, 0, today.Location()),
			"Design Meeting",
			"UI/UX discussion",
		},
		Meet{
			time.Date(today.Year(), today.Month(), today.Day(), 18, 0, 0, 0, today.Location()),
			"Weekly Sync",
			"Team weekly sync-up",
		},
	)
	return list
}
