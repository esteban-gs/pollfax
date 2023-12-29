package dto

import "time"

type BillRes struct {
	Id                      int64     `db:"id" json:"id"`
	LatestActionText        string    `db:"latest_action_text" json:"latestActionText"`
	Number                  string    `db:"bill_number" json:"number"`
	OriginChamber           string    `db:"origin_chamber" json:"originChamber"`
	Title                   string    `db:"title" json:"title"`
	Type                    string    `db:"type" json:"type"`
	UpdateDateIncludingText time.Time `db:"update_including_text" json:"updateDateIncludingText"`
	URL                     string    `db:"url" json:"url"`
}
