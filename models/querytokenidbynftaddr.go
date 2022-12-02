package models

import "strings"

func (nft NftDb) QueryTokenIdByNftaddr(nftaddr string) (string, error) {
	nftaddr = strings.ToLower(nftaddr)

	var nftRecord Nfts
	err := nft.db.Where("nftaddr = ?", nftaddr).First(&nftRecord)
	if err.Error != nil {
		return "", ErrNftNotExist
	}

	return nftRecord.Tokenid, nil
}
