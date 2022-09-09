package models

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"
)

func (nft NftDb) SetAgreement(param string, agreement multipart.File, name string) error {

	newPath := ImageDir + "/agreement/"
	_, err := os.Stat(newPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			log.Println("SetAgreement() create dir err=", err)
			return errors.New(ErrServer.Error() + "create dir err =" + err.Error())
		}
	}
	if path.Ext(name) != ".pdf" {
		log.Println("file type not pdf ,filename = ", name)
		return errors.New(ErrData.Error() + "file type err")
	}
	switch param {
	case UserAgreement:
		filename := newPath + UserAgreementFileName
		file6, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
		if err != nil {
			fmt.Println("SetAgreement creat file error")
			return errors.New("SetAgreement creat file error")
		}
		buf := &bytes.Buffer{}
		buf.ReadFrom(agreement)
		file6.Write(buf.Bytes())
	case PrivacyAgreement:
		filename := newPath + PrivacyAgreementFileName
		file6, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
		if err != nil {
			fmt.Println("SetAgreement creat file error")
			return errors.New("SetAgreement creat file error")
		}
		buf := &bytes.Buffer{}
		buf.ReadFrom(agreement)
		file6.Write(buf.Bytes())
	default:
		return errors.New(ErrData.Error() + "input agreement type error")
	}
	return nil
}

const (
	UserAgreement            = "user"
	PrivacyAgreement         = "privacy"
	UserAgreementFileName    = "useragreement.pdf"
	PrivacyAgreementFileName = "privacypolicy.pdf"
	DefaultPrivacyAgreement  = "defaultprivacy.pdf"
	DefaultUserAgreement     = "defaultuser.pdf"
)
