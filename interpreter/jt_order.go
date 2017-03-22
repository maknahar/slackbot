package interpreter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/maknahar/go-utils"
	"github.com/nlopes/slack"
)

func GetSlackMessage() slack.PostMessageParameters {
	msg := slack.PostMessageParameters{}
	msg.AsUser = true
	msg.Attachments = append(msg.Attachments, slack.Attachment{
		Color:      "#9932CC",
		AuthorName: "Justickets Bot",
		AuthorLink: "https://github.com/maknahar/jtbot",
		AuthorIcon: ":robot_face:",
		Footer:     "Always in your service",
		FooterIcon: "http://cconnect.s3.amazonaws.com/wp-content/uploads/2017/02/2017-Funko-Pop-Mystery-Science-Theater-3000-Crow-T-Robot-e1486480774184.jpg",
	})
	return msg
}

type Order struct {
	SessionID string `json:"sessionId"`
	BlockCode string `json:"blockCode"`
	Bill      struct {
		Total       float64 `json:"total"`
		BoBlockCode string  `json:"boBlockCode"`
		BlockCode   string  `json:"blockCode"`
		Breakups    []struct {
			SeatID             string  `json:"seatId"`
			SeatClass          string  `json:"seatClass"`
			Total              float64 `json:"total"`
			PriceCardID        string  `json:"priceCardId"`
			PriceCardLineItems []struct {
				Code        string  `json:"code"`
				Name        string  `json:"name"`
				Price       float64 `json:"price"`
				PriceType   string  `json:"priceType"`
				Description string  `json:"description"`
			} `json:"priceCardLineItems"`
			BookingChargeID        string `json:"bookingChargeId"`
			BookingChargeLineItems []struct {
				Code        string  `json:"code"`
				Name        string  `json:"name"`
				Price       float64 `json:"price"`
				PriceType   string  `json:"priceType"`
				Description string  `json:"description"`
			} `json:"bookingChargeLineItems"`
			OfferID        string      `json:"offerID"`
			OfferLineItems interface{} `json:"offerLineItems"`
		} `json:"breakups"`
	} `json:"bill"`
	Paid      bool `json:"paid"`
	Confirmed bool `json:"confirmed"`
	Notified  bool `json:"notified"`
	Cancelled bool `json:"cancelled"`
	Refunded  bool `json:"refunded"`
	Payments  []struct {
		ID          string `json:"ID"`
		BlockCode   string `json:"BlockCode"`
		RedirectURL string `json:"RedirectURL"`
		PostURL     string `json:"PostURL"`
		IframeURL   string `json:"IframeURL"`
		Finalized   bool   `json:"Finalized"`
		Successful  bool   `json:"Successful"`
		Mode        string `json:"Mode"`
	} `json:"payments"`
	BookingCode     string `json:"bookingCode"`
	Name            string `json:"name"`
	Mobile          string `json:"mobile"`
	Email           string `json:"email"`
	UserID          string `json:"userId"`
	Channel         string `json:"channel"`
	AssistedOrderID struct {
		String string `json:"String"`
		Valid  bool   `json:"Valid"`
	} `json:"assistedOrderID"`
	RedirectURL string `json:"redirectURL"`
}

func GetOrder(msg string) (*Order, error) {
	orderID := ""
	for _, v := range strings.Split(msg, " ") {
		if go_utils.IsValidUUIDV4(v) {
			orderID = v
			break
		}
	}
	if orderID == "" {
		return nil, fmt.Errorf("You are almost there. I can feel it. Try something like `Show me the order details of <order id>`")
	}
	res, err := http.Get("https://pm.justickets.co/orders?block_code=" + orderID)
	if err != nil {
		return nil, fmt.Errorf("Error in getting order data %v", err)
	}
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error in reading order data %v", err)
	}

	order := new(Order)
	err = json.Unmarshal(d, order)
	if err != nil {
		log.Println("Error in decoding jt order", err, go_utils.JsonPrettyPrint(string(d), "", "\t"))
		return order, nil
	}

	if order.SessionID == "" {
		return order, nil
	}

	return order, nil
}
