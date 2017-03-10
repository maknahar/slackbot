package interpreter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/maknahar/go-utils"
)

func GetOrder(msg string) (string, error) {
	orderID := ""
	for _, v := range strings.Split(msg, " ") {
		if go_utils.IsValidUUIDV4(v) {
			orderID = v
			break
		}
	}
	if orderID == "" {
		return "", fmt.Errorf("You are almost there. I can feel it. Try something like `Show me the order details of <order id>`")
	}
	res, err := http.Get("https://pm.justickets.co/orders?block_code=" + orderID)
	if err != nil {
		return "", fmt.Errorf("Error in getting order data %v", err)
	}
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Error in reading order data %v", err)
	}
	return go_utils.JsonPrettyPrint(string(d), "", "\t"), nil
}
