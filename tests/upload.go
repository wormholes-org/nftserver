package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func Upload(userAddr, collections, name, token, meta string, workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "upload"
	datam := make(map[string]string)
	datam["asset_sample"] = "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/2wBDAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAARCAEAAQADAREAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwD+/igAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgBiSJIu+N0kXLLuRg67kYqw3KSMqwKsM5DAg4IIo7eaTXmmk013TTTT2aaa0Dv5Np+TW6fmuqH0AFACZAIBIBPQZ5OOTj1wKP03D9dhaACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKAP50/8Agn3+0D8T4PBPhX9mrwP45+GHwWs9N0L9vH9o2bx38TfC914yuPHkej/8FBf2iPCOu+H/AAlo3/Cb+ALGy0X4ZWyaP4i+JmvtqWs3Vhb+OPBkDWWiW11/aGo3Xrr/AFXy/MauMwmGjwn4KfR9xH1KukqlbDZ14TYGtWz7F1pV6csPkuV1shpZZWnGiqdTG5nT9pj8POlSw+JwoUbcRZjhKeHxFepxR4t+M+DVek244evkfF2VrC5ZQpKlN18zzanxDjMZRpymnHC5NWdChiubEzwXpuv/APBRn4/eIvgD8bfjlpXiX4HfAGb9n/8AYE+F/wC1pq/hP4meC9a8X/8ACf8AiX4peCfiL4k017O7b4meAb3w58OjrXw/j8LeGpIrPWNb8Qavr6wS3VteWdvpuo8PE1atkmXcS5hh6MqtbL8/ybIMvyuteNehjM1ynIsfQwePq8sXPM8wzHiP/V/LqKpYeksz4azVJY5V5Usv7uH6dLOM14fy2tXpUKGYPO8RiMzUorD4zK8m4lx2T1s7wSnUVOlk+GynLKfEmY4r2+Jo0Mr4jymq8ZQpUo4nHdN8Rv2//wBoPRvhh+1b8ZfDs3wr0x/2SPCnwJ1C1+But+Edb1Lxj+0d4h+K/wAF/hx8UdM0zwtr9p41sr7wnL8VfGHjy9+C3wasNJ8I+OLvUfHHhm6Q/wDCQ6lPceGtO9qjh6c+IMPgJU69TD43xeq+GscuouH9qYPBriTJsipY2U52p1sxqZZnFPi2lh6lDDYepk9TBOcqeGryzE4MJNYnJMBjK+MwuBq1PDOHHeMzOvZZZRrvCcQ4l0J06lajKjl0MRkSy2tWqYmFSOMrYmFKdRUaDq/rZ4t12exuvCF6ml2i3dzY+Jr1Bqdt5t9pc9r4XudQEMLxyoIJWkjFrfbS/mRCSJWXIcdGGw0eTOKccRKcKFGmlPDztRxMVmeEoqUk4tzoyU3Vpp2amqcntZ/P4/HTcuFa9TBU6VbHYucp0sXTc8RgKksgzPFyjSk/ZypYiEqX1arOUE3SnWpuEXP3ebvfGXjSx0vVNSa70Gb+y/AmleNXh/si8j8+S8+2tPpav/az+XB5dk/l3RRpVeRN0TKjbvTp5ZllTGUcLyYuKrZ5PJ1L6xSbjDnpU44hr6ur1FKom6aahKKaunY+drcRcQ0srxGPVbK5TocGx4qcHgcQoyrRo1azwS/25tUKkaUk6r/e05crSnFtKXWfife6frIj060Oqaa0/iew8hrW2tZJdQ8N+Hr3WLqCxuBrE2pXUkN3Zx2Ny7+H4bI/bEMN2ZUijuccPklOrhJzrVPYVlh8LioVOec4KhiswoYKnUrQ+rKhTpuFWpVSeN9tehO9L2blKl2Yvi2vSzLD08JSeKw0sfictrUPY0acni8Lk2KzGvh8NXeOeKr4ilWpYenOSytYJU8XTi8TGtyRqrc/EDW9KiW7nutF1qOf4dX3jWO306ze2eGSK+0W2gbzpNWmSXS1j1G5Z5ZPs7SC1dzcRKjhbWT4WrUqUIxxOGlSzzB5S61apGUJQrQxrqScFQhy1m8NTcIqcox9ootS5kzD/WnMcPRoYurVwGPhieDs14nWEwlCpTq06mEqZQqNONV4qq6uFjHH14VZulCcpUeZShZwXbeE9V8Q6hcahFrVtDFBHbadd6fPnSIbqWO8W4EnmWmk+IvEcQtWa3EtlevcwfaVeaJYW+ymaXzMww+Eowpyw8p+09tiKNaD+syhF0lSklz4jBYKSrx9ry1qKhPkXspuS9ryR+jyTHZni6taOOpwVB4TB4vDVv8AhPp1KkcTLER5o0sBm+bQnhKioc+GxUqlJVZKvSjGfsHUl2teWfRBQAUAFABQAUAFABQAUAFABQAUAFABQAUAFABQAUAFABQAUAFABQAUAeU+I/gP8DvGOgaT4U8XfBr4U+KvC2g63f8AiXQ/DXiP4d+ENc0DRvEeqXd9f6n4g0nR9T0e607Ttb1G+1PUr2/1Wztob+8u9Qvrm4nkmu7h5CPu1sJiI+7iMBhqWDwNeOlbBYOhDDU6GFwlVWnhsNRp4LB06VCjKFKnDCYaEIRjQpKDu1DE002qeMxEsZi6d3yYrFzdaU8ViY/DXxE5YnESlXqqVWTr1m5N1Z83zh+1r+xTpv7WMUmh6/4t8MeG/CmreELvwDrV2fgt8PfFfxY8N+E9dkuLTxrafBT4t6/C+p/CbUvHfhi9vPC/iDVpNE8ZyWFobTVPCMHhjXLZ9RueZ4PC4jFVamY0VjsLiHQ+uYVTq4SpmVCi60pZVmeOw1SGLxGSYqdWSxmDoTwmMrYfF5rho5lTpY+H1Taliq+EjhZYCUMPiMFN1cFUq0qeKw+AxUXTlh80wOCqpUKWb4GdOjiMDiK/1nA08XgMrrV8txNPC4ihjvpM/A/4NSa34I8T3Xwq+HmpeK/hrpFhoPw/8X6t4O8P6t4u8GaRplubWxsPDXinUtPutf0aC3heREWw1CDJlmZiXlkZvRrYqvWxuNzGc2sbmM69TG4iko0amIeIlUlWjUdJQ5qc/bVU6VvZqNSUFFRbRwUMLRw+BweXQhzYPAUsPSwlCq3VhRjhVSVCUVPmSqQdClNVElP2lOE780U16XNa21wUa4t4J2jEojM0McpjE0ZimCF1YoJYmaKULgSRsUfKkisYznFSUZSippKSjJpSSkpJSSa5kpRjJJ3SlFPdJms6VOo4SqU4TlSk503OEZOnNwlTcoOSbhJ05zg5Rs3Ccot8smnG1hYOkkb2Vo8ctslnLG1tCySWabtlrIpQq9sm99kDAxLvbCDcc0q1ZSUlVqKSqe2UlOSkq10/ap3uql0nzp810nfRGbwuFcHTeGoOnKh9VlTdGm4SwtnH6s4uPK6HK2vYtezs2uWzZD/ZGli5lvo9PsYdQlYO2oRWVoL3zVgkto5jcNCzvLFBLJDG0m8CJ3iIMbshr6xX9mqPt63skpJUvaT9mlOaqSShzcqUpxjOStZzipO7SZCwOCVZ4lYPCrEOVOTxCw9JVnKlTlRpSdXk53KnRnOlTbleFOUoRajJp85ofg220nVJ9Xlksp7qXTn0xEsdHstItjBPcR3N7c3UFqXF3qGoSW9mLq4Zo4NlpElvaQBpfM7MTmM6+GWFjGpGk68a8lVxFTENOnTnSo0qLqa0qFGNWs4QbnUk6r9pVnyQ5fJy/IaeCx8sxqVaNXELC1cHCVDA4bA80MRXo18ViMV9XVsTjMTPC4V1KqjRpR9k/Y4el7Sd+nsdN07S4Wt9MsLLToHkMrQWNrBaQtKyqrSNHbxxoZCqIpcqWKooJwoA4qtetXkpV61WtKMeWMqtSdSSjdy5U5uTUbyk7J2u27XbPXw2DwmChKng8LhsJTnLnnDDUKVCEp8qjzyjSjBSlyxjHmab5YpXskXayOkKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgD/9k="
	//datam["asset_sample"] = models.Default_image
	datam["categories"] = "Art"
	datam["collections"] = collections
	datam["count"] = "1"
	datam["creator_addr"] = userAddr
	datam["desc"] = "multi thread test."
	datam["hide"] = "true"
	datam["md5"] = "0a399bf22855b727ad488997db53cc7f"
	datam["meta"] = meta
	datam["name"] = name
	datam["nft_contract_addr"] = ""
	datam["nft_token_id"] = ""
	datam["owner_addr"] = userAddr
	datam["royalty"] = "200"
	datam["source_url"] = "http://116.236.41.244:9000/ipfs/QmeCZAgHfigGsAek4bEj8P1LyohaK3JGzxfp1U7QgyCbR2"
	datam["user_addr"] = userAddr

	datas, _ := json.Marshal(&datam)
	data, err := HttpSendSign(datas, workKey)
	if err != nil {
		fmt.Println("sign err=", err)
		return err
	}
	b, err := HttpSendRev(url, data, token)
	if err != nil {
		fmt.Println("modifyUserInfo() err=", err)
		return err
	}
	b = DelDataItem(b)
	var revData ResponseLogin
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("get resp failed, err", err)
		return err
	}
	if revData.Code != "200" {
		return errors.New(revData.Msg)
	}
	return nil
}

func UploadWithImage(userAddr, collections, name, token, meta, image string, workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "upload"
	datam := make(map[string]string)
	datam["asset_sample"] = image
	//datam["asset_sample"] = models.Default_image
	datam["categories"] = "Art"
	datam["collections"] = collections
	datam["count"] = "1"
	datam["creator_addr"] = userAddr
	datam["desc"] = "multi thread test."
	datam["hide"] = "true"
	datam["md5"] = "0a399bf22855b727ad488997db53cc7f"
	datam["meta"] = meta
	datam["name"] = name
	datam["nft_contract_addr"] = ""
	datam["nft_token_id"] = ""
	datam["owner_addr"] = userAddr
	datam["royalty"] = "200"
	datam["source_url"] = "http://116.236.41.244:9000/ipfs/QmeCZAgHfigGsAek4bEj8P1LyohaK3JGzxfp1U7QgyCbR2"
	datam["user_addr"] = userAddr

	datas, _ := json.Marshal(&datam)
	data, err := HttpSendSign(datas, workKey)
	if err != nil {
		fmt.Println("sign err=", err)
		return err
	}
	b, err := HttpSendRev(url, data, token)
	if err != nil {
		fmt.Println("modifyUserInfo() err=", err)
		return err
	}
	b = DelDataItem(b)
	var revData ResponseLogin
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("get resp failed, err", err)
		return err
	}
	if revData.Code != "200" {
		return errors.New(revData.Msg)
	}
	return nil
}

func attachField(bodyWriter * multipart.Writer, keyname, keyvalue string) error {
	if err := bodyWriter.WriteField(keyname, keyvalue); err != nil {
		log.Printf("Cannot WriteField: %s, err: %v", keyname, err)
		return err
	}
	return nil
}

func attachFile(bodyWriter * multipart.Writer, formname, filename string) error {
	fullname := filepath.Join(".", filename)
	file, err := os.Open(fullname)
	if err != nil {
		log.Printf("Cannot open file: %s , err: %v", fullname, err)
		return err
	}
	defer file.Close()

	// MD5
	md5hash := md5.New()
	if _, err = io.Copy(md5hash, file); err != nil {
		log.Printf("Cannot open md5 hash: %s , err: %v", fullname, err)
		return err
	}

	keyname  := filename + ".md5cksum"
	keyvalue := hex.EncodeToString(md5hash.Sum(nil)[:16])
	if err = attachField(bodyWriter, keyname, keyvalue); err != nil {
		log.Printf("Cannot WriteField: %s, err: %v", keyname, err)
		return err
	}

	// file
	part, err := bodyWriter.CreateFormFile(formname, filename)
	if err != nil {
		log.Printf("Cannot CreateFormFile for: %s , err: %v", filename, err)
		return err
	}
	file.Seek(0, 0)
	_, err = io.Copy(part, file)
	if err != nil {
		log.Printf("Cannot Copy file: %s , err: %v", fullname, err)
		return err
	}

	return nil
}

func UploadImage(userAddr, collections, name, token, meta string, workKey *ecdsa.PrivateKey) error {
	url := SrcUrl + "uploadNftImage"
	datam := make(map[string]string)
	datam["asset_sample"] = "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/2wBDAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAARCAEAAQADAREAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwD+/igAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgBiSJIu+N0kXLLuRg67kYqw3KSMqwKsM5DAg4IIo7eaTXmmk013TTTT2aaa0Dv5Np+TW6fmuqH0AFACZAIBIBPQZ5OOTj1wKP03D9dhaACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKAP50/8Agn3+0D8T4PBPhX9mrwP45+GHwWs9N0L9vH9o2bx38TfC914yuPHkej/8FBf2iPCOu+H/AAlo3/Cb+ALGy0X4ZWyaP4i+JmvtqWs3Vhb+OPBkDWWiW11/aGo3Xrr/AFXy/MauMwmGjwn4KfR9xH1KukqlbDZ14TYGtWz7F1pV6csPkuV1shpZZWnGiqdTG5nT9pj8POlSw+JwoUbcRZjhKeHxFepxR4t+M+DVek244evkfF2VrC5ZQpKlN18zzanxDjMZRpymnHC5NWdChiubEzwXpuv/APBRn4/eIvgD8bfjlpXiX4HfAGb9n/8AYE+F/wC1pq/hP4meC9a8X/8ACf8AiX4peCfiL4k017O7b4meAb3w58OjrXw/j8LeGpIrPWNb8Qavr6wS3VteWdvpuo8PE1atkmXcS5hh6MqtbL8/ybIMvyuteNehjM1ynIsfQwePq8sXPM8wzHiP/V/LqKpYeksz4azVJY5V5Usv7uH6dLOM14fy2tXpUKGYPO8RiMzUorD4zK8m4lx2T1s7wSnUVOlk+GynLKfEmY4r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Mbshr6xX9mqPt63skpJUvaT9mlOaqSShzcqUpxjOStZzipO7SZCwOCVZ4lYPCrEOVOTxCw9JVnKlTlRpSdXk53KnRnOlTbleFOUoRajJp85ofg220nVJ9Xlksp7qXTn0xEsdHstItjBPcR3N7c3UFqXF3qGoSW9mLq4Zo4NlpElvaQBpfM7MTmM6+GWFjGpGk68a8lVxFTENOnTnSo0qLqa0qFGNWs4QbnUk6r9pVnyQ5fJy/IaeCx8sxqVaNXELC1cHCVDA4bA80MRXo18ViMV9XVsTjMTPC4V1KqjRpR9k/Y4el7Sd+nsdN07S4Wt9MsLLToHkMrQWNrBaQtKyqrSNHbxxoZCqIpcqWKooJwoA4qtetXkpV61WtKMeWMqtSdSSjdy5U5uTUbyk7J2u27XbPXw2DwmChKng8LhsJTnLnnDDUKVCEp8qjzyjSjBSlyxjHmab5YpXskXayOkKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgAoAKACgD/9k="
	datam["categories"] = "Art"
	datam["collections"] = collections
	datam["count"] = "1"
	datam["creator_addr"] = userAddr
	datam["desc"] = "multi thread test."
	datam["hide"] = "true"
	datam["md5"] = "0a399bf22855b727ad488997db53cc7f"
	datam["meta"] = meta
	datam["name"] = name
	datam["nft_contract_addr"] = ""
	datam["nft_token_id"] = ""
	datam["owner_addr"] = userAddr
	datam["royalty"] = "200"
	datam["source_url"] = "http://116.236.41.244:9000/ipfs/QmeCZAgHfigGsAek4bEj8P1LyohaK3JGzxfp1U7QgyCbR2"
	datam["user_addr"] = userAddr

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	if err := attachField(bodyWriter, "asset_sample", datam["asset_sample"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "categories", datam["categories"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "collections", datam["collections"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "count", datam["count"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "creator_addr", datam["creator_addr"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "desc", datam["desc"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "hide", datam["hide"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "md5", datam["md5"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "name", datam["name"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "meta", datam["meta"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "nft_contract_addr", datam["nft_contract_addr"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "nft_token_id", datam["nft_token_id"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "owner_addr", datam["owner_addr"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "royalty", datam["royalty"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "source_url", datam["source_url"]); err != nil {
		return err
	}
	if err := attachField(bodyWriter, "user_addr", datam["user_addr"]); err != nil {
		return err
	}

	sigData := datam["user_addr"] + datam["creator_addr"] + datam["owner_addr"] + datam["md5"] + datam["name"] +
		datam["desc"] + datam["meta"] + datam["source_url"] + datam["nft_contract_addr"] + datam["nft_token_id"] +
		datam["categories"] + datam["collections"] + /*datam["asset_sample"] +*/ datam["hide"] + datam["royalty"] + datam["count"]
	sig, err := Sign([]byte(sigData), workKey)
	if err != nil {
		fmt.Println("sign err=", err)
		return err
	}
	datam["sig"] = sig
	if err := attachField(bodyWriter, "sig", datam["sig"]); err != nil {
		return err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err := http.NewRequest("POST", url, bodyBuffer)
	if err != nil {
		log.Printf("Cannot NewRequest: %s , err: %v", url, err)
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("token", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var revData ResponseLogin
	err = json.Unmarshal([]byte(b), &revData)
	if err != nil {
		fmt.Println("get resp failed, err", err)
		return err
	}
	if revData.Code != "200" {
		return errors.New(revData.Msg)
	}
	return nil
}
/*func main() {
	testCount := 1
	tKey, tokens := Mlogin(testCount)
	fmt.Println("upload() login end.")
	fmt.Println("start Test Upload.")
	wd := sync.WaitGroup{}
	for i := 0; i < testCount; i++ {
		wd.Add(1)
		go func(i int) {
			defer wd.Done()
			userAddr := crypto.PubkeyToAddress(tKey[i].LogKey.PublicKey).String()
			//err := NewCollect("collect_test_" + strconv.Itoa(i), userAddr, tokens[i], tKey[i].WorkKey)
			err := Upload(userAddr, "collect_test_" + strconv.Itoa(i), "upload_test_" + strconv.Itoa(i), tokens[i], tKey[i].WorkKey)
			if err != nil {
				fmt.Println("upload() err=", err)
			}
		}(i)
	}
	wd.Wait()
	fmt.Println("end test upload().")
}*/