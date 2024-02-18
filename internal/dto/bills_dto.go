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

type BillSentiment struct {
	BillId      int64 `db:"bill_id" json:"billId"`
	SentimentId int64 `db:"sentiment_id" json:"sentimentId"`
}

type CreateBillSentiment struct {
	Sentiment string    `db:"sentiment" json:"sentiment"`
	VotedOn   time.Time `db:"voted_on" json:"votedOn"`
	BillId    int64     `db:"bill_id" json:"billId"`
}
