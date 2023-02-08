package models

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

func (nft NftDb) BatchMakeOffer(userAddr, offerlist string) error {
	userAddr = strings.ToLower(userAddr)

	offerList := OfferList{}
	if offerlist == "" {
		log.Println("BatchMakeOffer() input offerList is null ")
		return ErrDataFormat
	}
	uerr := json.Unmarshal([]byte(offerlist), &offerList)
	if uerr != nil {
		log.Println("BatchMakeOffer() unmarshal offerList err = ", uerr)
		return ErrDataFormat
	}
	for _, s := range offerList.Snft {
		price, _ := strconv.ParseUint(offerList.Price, 10, 64)
		deadTime, _ := strconv.ParseInt(offerList.DeadTime, 10, 64)
		uerr := nft.MakeOffer(userAddr, s.ContractAddr, s.TokenId,
			offerList.PayChannel, offerList.CurrencyType, price, "", deadTime, "", "", offerList.AuthSig)
		if uerr != nil {
			log.Println("BatchMakeOffer() makeoffer err = ", uerr)
			return uerr
		}
	}
	return nil
}
