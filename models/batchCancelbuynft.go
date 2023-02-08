package models

import (
	"encoding/json"
	"log"
	"strings"
)

func (nft NftDb) BatchCancelBuy(userAddr, offerlist string) error {
	userAddr = strings.ToLower(userAddr)

	offerList := OfferList{}
	if offerlist == "" {
		log.Println("BatchCancelBuy() input offerList is null ")
		return ErrDataFormat
	}
	uerr := json.Unmarshal([]byte(offerlist), &offerList)
	if uerr != nil {
		log.Println("BatchCancelBuy() unmarshal offerList err = ", uerr)
		return ErrDataFormat
	}
	for _, s := range offerList.Snft {
		uerr := nft.CancelBuy(userAddr, s.ContractAddr, s.TokenId, "", "")
		if uerr != nil {
			log.Println("BatchCancelBuy() CancelBuy err = ", uerr)
			return uerr
		}
	}
	return nil
}
