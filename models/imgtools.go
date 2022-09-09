package models

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
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
	//i := strings.LastIndex(file, "jpeg")
	//if i != -1 {
	//	file = file[:i] + "jpg"
	//} else {
	//	file = file[:i] + "jpeg"
	//}
	//f, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0777)
	//if err != nil {
	//	fmt.Println("base64toJpeg() OpenFile() err=", err)
	//	return err
	//}
	//err = jpeg.Encode(f, m, nil)
	//if err != nil {
	//	fmt.Println("base64toJpeg() jpeg.Encode() err=", err)
	//	return err
	//}
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
		file = newPath + hexname + "." + "jpg"
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
	result := Base64AddMemory(file)
	if result != nil {
		log.Println("Base64AddMemory() err=", result)
		return result
	}
	return err
}

func SaveCollectionsBackgroundImage(path, user_addr, name, image_base64 string) error {
	newPath := path + "/user/" + strings.ToLower(user_addr) + "/collections/"
	_, err := os.Stat(newPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println("SaveCollectionsBackgroundImage() create dir err=", err)
			return err
		}
	}
	var file, imagetype, img string
	if image_base64 != "" {
		imagetype, img, err = ParseBase64Type(image_base64)
		if err != nil {
			fmt.Println("SaveCollectionsBackgroundImage() ParseBase64Type() err=", err)
			return err
		}
		//hexname := hex.EncodeToString([]byte(name))
		var hexname string
		for _, c := range name {
			hexname += fmt.Sprintf("%02x", c)
		}
		file = newPath + hexname + "background." + "jpg"
	} else {
		fmt.Println("SaveCollectionsBackgroundImage() image_base64==0 error.")
		return err
	}
	if img == "" || imagetype == "" {
		fmt.Println("SaveCollectionsBackgroundImage() imagetype error.")
		return err
	}
	switch imagetype {
	case "jpeg", "jpg":
		err = base64toJpeg(file, img)
		if err != nil {
			fmt.Println("SaveCollectionsBackgroundImage() base64toJpeg() err=", err)
			return err
		}
	default:
		fmt.Println("SaveCollectionsBackgroundImage() imagetype error.")
		return errors.New("SaveCollectionsBackgroundImage() imagetype error.")
	}
	result := Base64AddMemory(file)
	if result != nil {
		log.Println("Base64AddMemory() err=", result)
		return result
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
		file = newPath + hexname + "." + "jpg"
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
	result := Base64AddMemory(file)
	if result != nil {
		log.Println("Base64AddMemory() err=", result)
		return result
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
		file = newPath + token_id + "." + "jpg"
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
	result := Base64AddMemory(file)
	if result != nil {
		log.Println("Base64AddMemory() err=", result)
		return result
	}
	return err
}

func SavePartnerslogoImage(path, token_id, image_base64 string) error {
	newPath := path + "/partnerslogo/" + "image/"
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
		file = newPath + token_id + "." + "jpg"
	} else {
		fmt.Println("SaveImage() image_base64==0 error.")
		return err
	}
	if img == "" || imagetype == "" {
		fmt.Println("SaveImage() imagetype error.")
		return err
	}
	switch imagetype {
	case "jpg", "jpeg":
		err = base64toJpg(file, img)
		if err != nil {
			fmt.Println("SaveImage() base64toJpeg() err=", err)
			return err
		}
	default:
		fmt.Println("SaveImage() imagetype error.")
		return errors.New("SaveImage() imagetype error.")
	}
	result := Base64AddMemory(file)
	if result != nil {
		log.Println("Base64AddMemory() err=", result)
		return result
	}
	return err
}

func DelPartnerslogoImage(path, name string) error {
	newPath := path + "/partnerslogo/" + "image/"
	_, err := os.Stat(newPath)
	if os.IsNotExist(err) {
		fmt.Println("DelPartnerslogoImage() create dir err=", err)
		return err
	}
	err = os.Remove(newPath + name + ".jpg")
	if err != nil {
		fmt.Println("DelPartnerslogoImage del file err=", err)
	}
	return err
}
func SaveSnftCollectionImage(path, token_id, image_base64 string) error {
	newPath := path + "/snft/" + "image/"
	_, err := os.Stat(newPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			log.Println("SaveImage() create dir err=", err)
			return err
		}
	}
	var file, imagetype, img string
	if image_base64 != "" {
		imagetype, img, err = ParseBase64Type(image_base64)
		if err != nil {
			log.Println("SaveImage() ParseBase64Type() err=", err)
			return err
		}
		file = newPath + token_id + "." + "jpg"
	} else {
		log.Println("SaveImage() image_base64==0 error.")
		return err
	}
	if img == "" || imagetype == "" {
		log.Println("SaveImage() imagetype error.")
		return err
	}
	switch imagetype {
	case "jpg", "jpeg":
		err = base64toJpg(file, img)
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

func DelSnftDirAllImage(path string) error {
	//newPath := path + "/snft/" + "image/"
	newPath := "./snftcollect/"
	_, err := os.Stat(newPath)
	if os.IsNotExist(err) {
		fmt.Println("DelPartnerslogoImage() create dir err=", err)
		return nil
	}
	dir, err := ioutil.ReadDir(newPath)
	for _, d := range dir {
		err = os.RemoveAll(newPath + d.Name())
		if err != nil {
			fmt.Println("DelSnftDirAllImage del file err=", err)
			return err
		}
	}
	return nil
}

func base64toJpg(file, data string) error {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, formatString, err := image.Decode(reader)
	if err != nil {
		fmt.Println("base64toJpg() Decode() err=", err)
		return err
	}
	bounds := m.Bounds()
	fmt.Println("base64toJpg", bounds, formatString)
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("base64toJpg() OpenFile() err=", err)
		return err
	}
	err = jpeg.Encode(f, m, nil)
	if err != nil {
		fmt.Println("base64toJpg() jpg.Encode() err=", err)
		return err
	}

	return err
}

func IpfsTojpgbase64() (string, error) {
	fmt.Println("default worm:", DefaultWormBlue)
	url := NftIpfsServerIP + ":" + NftstIpfsServerPort
	s := shell.NewShell(url)
	s.SetTimeout(100 * time.Second)

	var err error
	var wormdata io.Reader
	for {
		wormdata, err = s.Cat(DefaultWormBlue)
		if err != nil {
			log.Printf("wormdata cat  [%v] failed! %v", DefaultWormBlue, err)
			time.Sleep(5 * time.Second)
			continue
		} else {
			break
		}
	}

	wormbody, err := ioutil.ReadAll(wormdata)
	if err != nil {
		fmt.Printf("Read http response failed! %v", err)
		return "", err
	}
	wormimg, _, err := image.Decode(bytes.NewReader(wormbody))
	if err != nil {
		fmt.Printf("image Decode  failed! %v", err)
		return "", err
	}

	emptyBuff := bytes.NewBuffer(nil)
	jpeg.Encode(emptyBuff, wormimg, nil)
	wormstr := base64.StdEncoding.EncodeToString(emptyBuff.Bytes())
	partlogo := "data:image/jpg;base64," + wormstr
	fmt.Println("IpfsTojpgbase64 ok")
	return partlogo, nil
}
