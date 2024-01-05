package handlers

import (
	"net/http"
	"pollfax/db"
	"pollfax/internal/dto"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func GetAll(c echo.Context) error {
	_db := db.Instance()
	bills := []dto.BillRes{}
	_db.Select(&bills, `SELECT title,
                              type,
													    bill_number,
													    origin_chamber,
													    url,
													    latest_action_text,
													    update_including_text
									           FROM bills`)

	return c.JSON(http.StatusOK, bills)
}

func GetAllCache(c echo.Context) error {
	return c.File("../bills.json")
}

func CreateBillSentiment(c echo.Context) error {
	sDto := new(dto.CreateBillSentiment)
	c.Bind(sDto)
	sDto.VotedOn = time.Now().UTC()

	_db := db.Instance()
	billSent := new(dto.BillSentiment)
	err := _db.QueryRowx(`WITH sentiment_cte AS (
		INSERT INTO sentiments (sentiment, voted_on)
		VALUES ($1, $2)
		RETURNING id
		),
	bill_sentiment_cte AS (
		INSERT INTO bills_sentiments (bill_id, sentiment_id)
		VALUES ($3, (SELECT s.id FROM sentiment_cte s))
		RETURNING bill_id, sentiment_id
	)

	SELECT 	bill_id, sentiment_id FROM Bill_sentiment_cte
`, sDto.Sentiment, sDto.VotedOn, sDto.BillId).StructScan(billSent)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting question")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, billSent)
}
