package interpreter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"errors"

	"os"

	"github.com/maknahar/go-utils"
)

type ReportResponse struct {
	CreatedAt     string `json:"CreatedAt"`
	CurrentTime   string `json:"CurrentTime"`
	FailureReason string `json:"FailureReason"`
	FromTime      string `json:"FromTime"`
	Status        string `json:"Status"`
	ToTime        string `json:"ToTime"`
	UpdatedAt     string `json:"UpdatedAt"`
}

func GetReportStatus(staging bool) (*ReportResponse, error) {
	url := os.Getenv("REPORT_STATUS")
	if staging {
		url = os.Getenv("STAGING_REPORT_STATUS")
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error in getting order data %v", err)
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error in reading order data %v", err)
	}

	order := make([]ReportResponse, 0)
	err = json.Unmarshal(d, &order)
	if err != nil {
		log.Println("Error in decoding report response", err, go_utils.JsonPrettyPrint(string(d),
			"", "\t"))
		return nil, errors.New("Error in decoding report response")
	}

	l := len(order)
	if l == 0 {
		return nil, errors.New("Report Service might be down")
	}

	return &order[l-1], nil
}

func (r *ReportResponse) GetDelayReason() string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	fromTime, err := time.ParseInLocation("2006-01-02T15:04:05.999999", r.FromTime, loc)
	if err != nil {
		return "Unable to parse to time. Err: %+v"
	}

	if strings.Contains(r.FailureReason, "SessionNotFound") {
		return fmt.Sprintf("A session is missing from data pack since %s", fromTime)
	}

	if r.FailureReason == "" {
		if time.Since(fromTime).Minutes() > 15 {
			return fmt.Sprintf("Report sync is running smoothly. However there is delay of %s",
				time.Since(fromTime).String())
		}
		return fmt.Sprintf("Report is in sync")
	}

	return ""
}
