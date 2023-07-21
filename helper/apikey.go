package helper

import (
	"math/rand"
	"time"
)

func RandomApiKey() string {
	apiKey := []string{"6IT85H1YWNXHJYD7IN3RFJ52YGPGMEY8WY", "6IT85H1YWNXHJYD7IN3RFJ52YGPGMEY8WY"}

	rand.Seed(time.Now().UnixNano())
	return apiKey[rand.Intn(len(apiKey))]
}

//0x39e64ff40caca9f8c579ca8573fc607a12cabd70
