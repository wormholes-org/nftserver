package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strings"
)

const (
	default_ipfs_hash = "QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o"
)

func ParseBase64Type(image_base64 string) (string, string, error) {
	i := strings.Index(image_base64, "data:image/")
	if i == -1 {
		fmt.Println("default_image() image_base64 error.")
		return "", "", errors.New("base64 data error.")
	}
	offset := i + len("data:image/")
	j := strings.Index(image_base64[offset:], ";")
	if i == -1 {
		fmt.Println("default_image() image_base64 error.")
		return "", "", errors.New("base64 data error.")
	}
	imagetype := image_base64[offset : offset+j]
	i = strings.Index(image_base64, "base64,")
	if i == -1 {
		fmt.Println("default_image() image_base64 error.")
		return "", "", errors.New("base64 data error.")
	}
	img := image_base64[i+len("base64,"):]
	return imagetype, img, nil
}

func base64toJpeg(file, data string) error {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, formatString, err := image.Decode(reader)
	if err != nil {
		fmt.Println("base64toJpeg() Decode() err=", err)
		return err
	}
	bounds := m.Bounds()
	fmt.Println("base64toJpg", bounds, formatString)
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("base64toJpeg() OpenFile() err=", err)
		return err
	}
	err = jpeg.Encode(f, m, nil)
	if err != nil {
		fmt.Println("base64toJpeg() jpeg.Encode() err=", err)
		return err
	}
	i := strings.LastIndex(file, "jpeg")
	if i != -1 {
		file = file[:i] + "jpg"
	} else {
		file = file[:i] + "jpeg"
	}
	f, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("base64toJpeg() OpenFile() err=", err)
		return err
	}
	err = jpeg.Encode(f, m, nil)
	if err != nil {
		fmt.Println("base64toJpeg() jpeg.Encode() err=", err)
		return err
	}
	return err
}

func SavePortrait(path, user_addr, image_base64 string) error {
	newPath := path + "/user/" + strings.ToLower(user_addr) + "/"
	_, err := os.Stat(newPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println("SavePortrait() create dir err=", err)
			return err
		}
	}
	var file, imagetype, img string
	if image_base64 != "" {
		imagetype, img, err = ParseBase64Type(image_base64)
		if err != nil {
			fmt.Println("SavePortrait() ParseBase64Type() err=", err)
			return err
		}
		file = newPath + "portrait." + imagetype
	} else {
		imagetype, img, err = ParseBase64Type(Default_image)
		if err != nil {
			fmt.Println("SavePortrait() ParseBase64Type() err=", err)
			return err
		}
		file = newPath + "portrait." + imagetype
	}
	if img == "" || imagetype == "" {
		fmt.Println("SavePortrait() imagetype error.")
		return errors.New("SavePortrait() imagetype error.")
	}
	switch imagetype {
	case "jpeg", "jpg":
		err = base64toJpeg(file, img)
		if err != nil {
			fmt.Println("SavePortrait() base64toJpeg() err=", err)
			return err
		}
	default:
		fmt.Println("SavePortrait() imagetype error.")
		return errors.New("SavePortrait() imagetype error.")
	}
	return err
}

func SaveBackground(path, user_addr, image_base64 string) error {
	newPath := path + "/user/" + strings.ToLower(user_addr) + "/"
	_, err := os.Stat(newPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println("SaveBackground() create dir err=", err)
			return err
		}
	}
	var file, imagetype, img string
	if image_base64 != "" {
		imagetype, img, err = ParseBase64Type(image_base64)
		if err != nil {
			fmt.Println("SaveBackground() ParseBase64Type() err=", err)
			return err
		}
		file = newPath + "background." + imagetype
	} else {
		imagetype, img, err = ParseBase64Type(Default_image)
		if err != nil {
			fmt.Println("SaveBackground() ParseBase64Type() err=", err)
			return err
		}
		file = newPath + "background." + imagetype
	}
	if img == "" || imagetype == "" {
		fmt.Println("SaveBackground() imagetype error.")
		return errors.New("SaveBackground() imagetype error.")
	}
	switch imagetype {
	case "jpeg", "jpg":
		err = base64toJpeg(file, img)
		if err != nil {
			fmt.Println("SaveBackground() base64toJpeg() err=", err)
			return err
		}
	default:
		fmt.Println("SaveBackground() imagetype error.")
		return errors.New("SaveBackground() imagetype error.")
	}
	return err
}

func SaveCollectionsImage(path, user_addr, name, image_base64 string) error {
	newPath := path + "/user/" + strings.ToLower(user_addr) + "/collections/"
	_, err := os.Stat(newPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println("SaveCollectionsImage() create dir err=", err)
			return err
		}
	}
	var file, imagetype, img string
	if image_base64 != "" {
		imagetype, img, err = ParseBase64Type(image_base64)
		if err != nil {
			fmt.Println("SaveCollectionsImage() ParseBase64Type() err=", err)
			return err
		}
		//hexname := hex.EncodeToString([]byte(name))
		var hexname string
		for _, c := range name {
			hexname += fmt.Sprintf("%02x", c)
		}
		file = newPath + hexname + "." + imagetype
	} else {
		fmt.Println("SaveCollectionsImage() image_base64==0 error.")
		return err
	}
	if img == "" || imagetype == "" {
		fmt.Println("SaveCollectionsImage() imagetype error.")
		return err
	}
	switch imagetype {
	case "jpeg", "jpg":
		err = base64toJpeg(file, img)
		if err != nil {
			fmt.Println("SaveCollectionsImage() base64toJpeg() err=", err)
			return err
		}
	default:
		fmt.Println("SaveCollectionsImage() imagetype error.")
		return errors.New("SaveCollectionsImage() imagetype error.")
	}
	return err
}

func SaveSnftCollectionsImage(path, user_addr, name, image_base64 string) error {
	newPath := path + "/user/" + strings.ToLower(user_addr) + "/snftcollections/"
	_, err := os.Stat(newPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println("SaveSnftCollectionsImage() create dir err=", err)
			return err
		}
	}
	var file, imagetype, img string
	if image_base64 != "" {
		imagetype, img, err = ParseBase64Type(image_base64)
		if err != nil {
			fmt.Println("SaveSnftCollectionsImage() ParseBase64Type() err=", err)
			return err
		}
		//hexname := hex.EncodeToString([]byte(name))
		var hexname string
		for _, c := range name {
			hexname += fmt.Sprintf("%02x", c)
		}
		file = newPath + hexname + "." + imagetype
	} else {
		fmt.Println("SaveSnftCollectionsImage() image_base64==0 error.")
		return err
	}
	if img == "" || imagetype == "" {
		fmt.Println("SaveSnftCollectionsImage() imagetype error.")
		return err
	}
	switch imagetype {
	case "jpeg", "jpg":
		err = base64toJpeg(file, img)
		if err != nil {
			fmt.Println("SaveSnftCollectionsImage() base64toJpeg() err=", err)
			return err
		}
	default:
		fmt.Println("SaveCollectionsImage() imagetype error.")
		return errors.New("SaveCollectionsImage() imagetype error.")
	}
	return err
}

func SaveNftImage(path, contract_addr, token_id, image_base64 string) error {
	newPath := path + "/nft/" + strings.ToLower(contract_addr) + "/image/"
	_, err := os.Stat(newPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println("SaveImage() create dir err=", err)
			return err
		}
	}
	var file, imagetype, img string
	if image_base64 != "" {
		imagetype, img, err = ParseBase64Type(image_base64)
		if err != nil {
			fmt.Println("SaveImage() ParseBase64Type() err=", err)
			return err
		}
		file = newPath + token_id + "." + imagetype
	} else {
		fmt.Println("SaveImage() image_base64==0 error.")
		return err
	}
	if img == "" || imagetype == "" {
		fmt.Println("SaveImage() imagetype error.")
		return err
	}
	switch imagetype {
	case "jpeg", "jpg":
		err = base64toJpeg(file, img)
		if err != nil {
			fmt.Println("SaveImage() base64toJpeg() err=", err)
			return err
		}
	default:
		fmt.Println("SaveImage() imagetype error.")
		return errors.New("SaveImage() imagetype error.")
	}
	return err
}

func SaveSNftImage(path, contract_addr, token_id, image_base64 string) error {
	newPath := path + "/snft/" + strings.ToLower(contract_addr) + "/image/"
	_, err := os.Stat(newPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println("SaveImage() create dir err=", err)
			return err
		}
	}
	var file, imagetype, img string
	if image_base64 != "" {
		imagetype, img, err = ParseBase64Type(image_base64)
		if err != nil {
			fmt.Println("SaveImage() ParseBase64Type() err=", err)
			return err
		}
		file = newPath + token_id + "." + imagetype
	} else {
		fmt.Println("SaveImage() image_base64==0 error.")
		return err
	}
	if img == "" || imagetype == "" {
		fmt.Println("SaveImage() imagetype error.")
		return err
	}
	switch imagetype {
	case "jpeg", "jpg":
		err = base64toJpeg(file, img)
		if err != nil {
			fmt.Println("SaveImage() base64toJpeg() err=", err)
			return err
		}
	default:
		fmt.Println("SaveImage() imagetype error.")
		return errors.New("SaveImage() imagetype error.")
	}
	return err
}