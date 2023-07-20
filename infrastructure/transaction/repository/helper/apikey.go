package helper

import (
	"math/rand"
	"time"
)

func RandomApiKey() string {
	apiKey := []string{"6IT85H1YWNXHJYD7IN3RFJ52YGPGMEY8WY", "6IT85H1YWNXHJYD7IN3RFJ52YGPGMEY8WY"}

	rand.Seed(time.Now().UnixNano()) // Initialize random seed based on current time

	min := 0
	max := len(apiKey) - 1
	return apiKey[rand.Intn(max-min+1)+min]
}
