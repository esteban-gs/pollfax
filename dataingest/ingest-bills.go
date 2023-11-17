package dataingest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"pollfax/db"
	"regexp"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type CongressResponse struct {
	Congresses []struct {
		EndYear  string `json:"endYear"`
		Name     string `json:"name"`
		Sessions []struct {
			Chamber   string `json:"chamber"`
			Number    int    `json:"number"`
			StartDate string `json:"startDate"`
			Type      string `json:"type"`
		} `json:"sessions"`
		StartYear string `json:"startYear"`
		URL       string `json:"url"`
	} `json:"congresses"`
	Pagination struct {
		Count int    `json:"count"`
		Next  string `json:"next"`
	} `json:"pagination"`
	Request struct {
		ContentType string `json:"contentType"`
		Format      string `json:"format"`
	} `json:"request"`
}

type Bill struct {
	Congress     int `json:"congress"`
	LatestAction struct {
		ActionDate string `json:"actionDate"`
		Text       string `json:"text"`
	} `json:"latestAction"`
	Number                  string    `json:"number"`
	OriginChamber           string    `json:"originChamber"`
	OriginChamberCode       string    `json:"originChamberCode"`
	Title                   string    `json:"title"`
	Type                    string    `json:"type"`
	UpdateDate              string    `json:"updateDate"`
	UpdateDateIncludingText time.Time `json:"updateDateIncludingText"`
	URL                     string    `json:"url"`
	Created                 time.Time `json:"created"`
}

type BillsResponse struct {
	Bills []Bill `json:"bills"`
}

func GetLatestCongress() (number int64) {
	apiKey := os.Getenv("CONGRESS_API_KEY")
	params := url.Values{}
	params.Add("format", "json")
	params.Add("limit", "1")
	params.Add("api_key", apiKey)
	dataUrl := fmt.Sprintf("https://api.congress.gov/v3/congress?%s", params.Encode())

	log.Info().Msg("Latest congress query running")

	res, err := http.Get(dataUrl)
	if err != nil {
		log.Err(err).Msg("Http request error")
	}
	reqBody, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Err(err).Msg("Could not read request body")
	}

	var result CongressResponse
	if err := json.Unmarshal(reqBody, &result); err != nil {
		log.Err(err).Msg("Can not unmarshal JSON")
	}

	re := regexp.MustCompile(`(\d+)[A-Za-z]{0,4}th Congress`)
	match := re.FindStringSubmatch(result.Congresses[0].Name)

	if len(match) > 0 {
		number, err = strconv.ParseInt(match[1], 10, 0)
	} else {
		number = 0
	}
	log.Info().
		Int64("Congress", number).
		Msg("Last congress found")
	return
}

func GetBills(congress int64) (bills []Bill) {
	apiKey := os.Getenv("CONGRESS_API_KEY")
	params := url.Values{}
	params.Add("format", "json")
	params.Add("limit", "50")
	params.Add("sort", "updateDate+desc")
	params.Add("api_key", apiKey)
	dataUrl := fmt.Sprintf("https://api.congress.gov/v3/bill/%d?%s", congress, params.Encode())

	log.Info().Msg("Getting bills")

	res, err := http.Get(dataUrl)
	if err != nil {
		log.Err(err).Msg("Http request error")
	}
	reqBody, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Err(err).Msg("Could not read request body")
	}

	var response BillsResponse
	if err := json.Unmarshal(reqBody, &response); err != nil {
		log.Err(err).Msg("Can not unmarshal JSON")
	}
	log.Info().Int("count", len(response.Bills)).Msg("Bills found")
	return response.Bills
}

func Clear() {
	_db := db.Instance()
	tx := _db.MustBegin()
	log.Info().Str("dataingest", "Clear").Msg("Removing existing bills")
	tx.Exec(`TRUNCATE bills`)
	commitErr := tx.Commit()
	if commitErr != nil {
		log.Error().Err(commitErr).Msg("Error Clearing Bills")
	}
}

func Persist(bills *[]Bill) {
	_db := db.Instance()
	tx := _db.MustBegin()
	for _, bill := range *bills {
		log.Info().Str("bill", bill.Number+bill.Type).Msg("Saving bill")
		tx.Exec(`INSERT INTO bills
		(congress,
		bill_number,
		origin_chamber,
		origin_chamber_code,
		title,
		type,
		url,
		latest_action_date,
		latest_action_text,
		update_date,
		update_including_text,
		created)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			bill.Congress,
			bill.Number,
			bill.OriginChamber,
			bill.OriginChamberCode,
			bill.Title,
			bill.Type,
			bill.URL,
			bill.LatestAction.ActionDate,
			bill.LatestAction.Text,
			bill.UpdateDate,
			bill.UpdateDateIncludingText,
			time.Now().UTC(),
		)
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		log.Error().Err(commitErr).Msg("Error inserting Bills")
	}
}

func Run() {
	Clear()
	congress := GetLatestCongress()
	bills := GetBills(congress)
	Persist(&bills)
}
