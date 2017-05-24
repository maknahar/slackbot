package interpreter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"strconv"

	"github.com/maknahar/go-utils"
	"github.com/nlopes/slack"
)

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

func (o *Order) FormatSlackMessage(attachment *slack.Attachment) {
	if o.SessionID == "" {
		attachment.Pretext = "No Order found for id " + o.BlockCode
		attachment.Color = "#BDB76B"
		attachment.Title = "Search On Admin"
		attachment.TitleLink = "https://admin.justickets.co/bookings?detail=" + o.BlockCode
		return
	}
	attachment.Pretext = "Found Order: " + o.BlockCode
	attachment.Title = "Movie Ticket:"
	attachment.TitleLink = "https://www.justickets.in/orders/" + o.BlockCode
	attachment.Fields = append(attachment.Fields, slack.AttachmentField{
		Title: "Name",
		Value: o.Name,
		Short: true})
	attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Email",
		Value: o.Email,
		Short: true})
	attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Mobile",
		Value: o.Mobile,
		Short: true})
	attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Paid",
		Value: strconv.FormatBool(o.Paid),
		Short: true})
	attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Confirmed",
		Value: strconv.FormatBool(o.Confirmed),
		Short: true})
	attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Failed",
		Value: strconv.FormatBool(o.Cancelled),
		Short: true})
	attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Bill Total",
		Value: strconv.FormatFloat(o.Bill.Total, 'G', 6, 64),
		Short: true})
	attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Channel",
		Value: o.Channel,
		Short: true})
	attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Session ID",
		Value: o.SessionID,
		Short: false})
	seats := ""
	for _, v := range o.Bill.Breakups {
		seats += v.SeatClass + " " + strings.TrimLeft(v.SeatID, v.SeatClass) + ", "
	}
	strings.TrimRight(seats, ",")
	attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Seats",
		Value: seats,
		Short: false})
	if o.AssistedOrderID.String != "" {
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Session ID",
			Value: o.AssistedOrderID.String,
			Short: false})
	}
	if o.UserID != "" {
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "User ID",
			Value: o.UserID,
			Short: false})
	}
	if o.Confirmed {
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Booking Code",
			Value: o.BookingCode,
			Short: true})
	}
}

func (o *Order) FormatSlackMessageForBill(attachment *slack.Attachment) {
	if o.SessionID == "" {
		attachment.Pretext = "No Order found for id " + o.BlockCode
		attachment.Color = "#BDB76B"
		attachment.Fields = append(attachment.Fields,
			slack.AttachmentField{
				Title: "Seach On Admin",
				Value: "https://admin.justickets.co/bookings?detail=" + o.BlockCode,
				Short: true})
		return
	}
	attachment.Pretext = "Found Bill for " + o.BlockCode
	attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Gross Total",
		Value: strconv.FormatFloat(o.Bill.Total, 'G', 6, 64),
		Short: true})
}
