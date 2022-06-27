package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nftexchange/nftserver/models"
)

// The settings required by the token for API access
const (
	PRIVATEKEY                  = "NFTEXCHANGER.WORMHOLES.202110191729"
	DEFAULT_EXPIRE_TIME_SECONDS = 60 * 60 * 10
)

type User struct {
	UserAddr string
}
type ExchangerCustomClaims struct {
	User
	jwt.StandardClaims
}

//Server response
type HttpResponseData struct {
	Code       string      `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
	TotalCount uint64      `json:"total_count"`
}

// To facilitate parsing the filter field, this structure is defined to parse the http request data,
// The generic way map[string]string is not used
type HttpRequestFilter struct {
	Match      string                `json:"match"`
	Filter     []models.StQueryField `json:"filter"`
	Sort       []models.StSortField  `json:"sort"`
	Nfttype    string                `json:"nfttype"`
	StartIndex string                `json:"start_index"`
	Count      string                `json:"count"`
}
