package interpreter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/maknahar/go-utils"
	"github.com/nlopes/slack"
)

type ReportResponse struct {
	CreatedAt      string `json:"CreatedAt"`
	CurrentTime    string `json:"CurrentTime"`
	FailureReason  string `json:"FailureReason"`
	FromTime       string `json:"FromTime"`
	Status         string `json:"Status"`
	ToTime         string `json:"ToTime"`
	UpdatedAt      string `json:"UpdatedAt"`
	Link           string
	MissingSession string
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
		if res.StatusCode == http.StatusServiceUnavailable {
			return nil, errors.New("Report Service is unavailable")
		}
		log.Println("Error in decoding report response", err, go_utils.JsonPrettyPrint(string(d),
			"", "\t"))
		return nil, errors.New("Error in decoding report response")
	}

	l := len(order)
	if l == 0 {
		return nil, errors.New("Report Service might be down")
	}
	activity := &order[l-1]
	activity.Link = url

	return activity, nil
}

func (r *ReportResponse) GetDelayReason() string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	fromTime, err := time.ParseInLocation("2006-01-02T15:04:05.999999", r.FromTime, loc)
	if err != nil {
		return "Unable to parse to time. Err: %+v"
	}

	if strings.Contains(r.FailureReason, "SessionNotFound") {
		s := regexp.MustCompile("[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
		r.MissingSession = s.FindString(r.FailureReason)
		return fmt.Sprintf("A Session is missing from data pack since %s", fromTime.Format("02-01-2006 03:04pm"))
	}

	if r.FailureReason == "" {
		if time.Since(fromTime).Minutes() > 15 {
			return fmt.Sprintf("Report sync is running smoothly.")
		}
		return fmt.Sprintf("Report is in sync.")
	}

	return r.FailureReason
}

func (r *ReportResponse) FormatSlackMessage(attachment *slack.Attachment) {
	cause := r.GetDelayReason()
	attachment.Pretext = cause
	attachment.Title = "Report Sync Activity Details"
	attachment.TitleLink = r.Link
	if r.MissingSession != "" {
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Missing Session ID:",
			Value: r.MissingSession,
			Short: true})
	}

	attachment.Fields = append(attachment.Fields, slack.AttachmentField{
		Title: "Status:",
		Value: r.Status,
		Short: true})

	loc, _ := time.LoadLocation("Asia/Kolkata")
	fromTime, err := time.ParseInLocation("2006-01-02T15:04:05.999999", r.FromTime, loc)
	if err == nil {
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Delay:",
			Value: time.Since(fromTime).String(),
			Short: true})
	}

	ct, err := time.ParseInLocation("2006-01-02T15:04:05.999999", r.CreatedAt, loc)
	if err == nil {
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Created At:",
			Value: ct.Format("02-01-2006 03:04pm"),
			Short: true})
	}

	ut, err := time.ParseInLocation("2006-01-02T15:04:05.999999", r.UpdatedAt, loc)
	if err == nil {
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Last Updated At:",
			Value: ut.Format("02-01-2006 03:04pm"),
			Short: true})
	}

}
