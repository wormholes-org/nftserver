package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

type CountryRec struct {
	Regionen string `json:"regionen" gorm:"type:char(50) CHARACTER SET utf8mb4;comment:'country/region English'"`
	Regioncn string `json:"regioncn" gorm:"type:char(50) CHARACTER SET utf8mb4;comment:'country/region Chinese'"`
	Domain   string `json:"domain" gorm:"type:char(10) CHARACTER SET utf8mb4;comment:'international domain name'"`
	Telecode string `json:"telecode" gorm:"type:char(10) CHARACTER SET utf8mb4 NOT NULL;comment:'phone code'"`
}

type Countrys struct {
	gorm.Model
	CountryRec
}

func (v Countrys) TableName() string {
	return "countrys"
}

func (nft NftDb) QueryCountrys() ([]CountryRec, error) {
	countrys := make([]Countrys, 0, 20)
	db := nft.db.Model(&Countrys{}).Find(&countrys)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			fmt.Println("QueryCountry() dbase err=", db.Error)
			return nil, db.Error
		}
	} else {
		countryRecs := make([]CountryRec, 0, 20)
		for _, cn := range countrys {
			countryRecs = append(countryRecs, cn.CountryRec)
		}
		return countryRecs, nil
	}
}

func (nft NftDb) ModifyCountry(Regionen, Regioncn, Domain, Telecode string) error {
	//UserSync.Lock(Adminaddr)
	//defer UserSync.UnLock(Adminaddr)
	country := Countrys{}
	db := nft.db.Model(&Countrys{}).Where("regionen = ? AND regioncn = ?", Regionen, Regioncn).First(&country)
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			country.Regionen = Regionen
			country.Regioncn = Regioncn
			country.Domain = Domain
			country.Telecode = Telecode
			db = nft.db.Model(&Countrys{}).Create(&country)
			if db.Error != nil {
				fmt.Println("ModifyCountry()->create() err=", db.Error)
				return db.Error
			}
		} else {
			fmt.Println("ModifyCountry() dbase err=", db.Error)
			return db.Error
		}
	} else {
		country := Countrys{}
		country.Domain = Domain
		country.Telecode = Telecode
		db = nft.db.Model(&Countrys{}).Where("regionen = ? AND regioncn = ?", Regionen, Regioncn).Updates(&country)
		if db.Error != nil {
			fmt.Printf("ModifyCountry()->UPdate() users err=%s\n", db.Error)
		}
	}
	return db.Error
}

func (nft NftDb) DefaultCountry() error {
	data := make([]map[string]string, 0, 20)
	json.Unmarshal([]byte(CountryList), &data)
	nft.db.Transaction(func(tx *gorm.DB) error {
		err := tx.First(&Countrys{})
		if err.Error != nil {
			countrys := make([]Countrys, 0, 20)
			for _, value := range data {
				country := Countrys{}
				country.Regionen = value["regionen"]
				country.Regioncn = value["regioncn"]
				country.Domain = value["domain"]
				country.Telecode = value["telecode"]
				countrys = append(countrys, country)
			}
			if err.Error == gorm.ErrRecordNotFound {
				err = tx.Model(&Countrys{}).Create(&countrys)
				if err.Error != nil {
					fmt.Println("DefaultCountry() error create Countrys record err=", err.Error)
					return err.Error
				}
			}
			return nil
		}
		return nil
	})
	return nil
}

const CountryList = `[
  {
    "regioncn": "阿富汗",
    "regionen": "Afghanistan",
    "domain": "AF",
    "telecode": "93"
  },
  {
    "regioncn": "奥兰群岛",
    "regionen": "Aland Islands",
    "domain": "AX",
    "telecode": "35818"
  },
  {
    "regioncn": "阿尔巴尼亚",
    "regionen": "Albania",
    "domain": "AL",
    "telecode": "355"
  },
  {
    "regioncn": "阿尔及利亚",
    "regionen": "Algeria",
    "domain": "DZ",
    "telecode": "213"
  },
  {
    "regioncn": "美属萨摩亚",
    "regionen": "American Samoa",
    "domain": "AS",
    "telecode": "1684"
  },
  {
    "regioncn": "安道尔",
    "regionen": "Andorra",
    "domain": "AD",
    "telecode": "376"
  },
  {
    "regioncn": "安哥拉",
    "regionen": "Angola",
    "domain": "AO",
    "telecode": "244"
  },
  {
    "regioncn": "安圭拉",
    "regionen": "Anguilla",
    "domain": "AI",
    "telecode": "1264"
  },
  {
    "regioncn": "南极洲",
    "regionen": "Antarctica",
    "domain": "AQ",
    "telecode": "6721"
  },
  {
    "regioncn": "安提瓜和巴布达",
    "regionen": "Antigua and Barbuda",
    "domain": "AG",
    "telecode": "1268"
  },
  {
    "regioncn": "阿根廷",
    "regionen": "Argentina",
    "domain": "AR",
    "telecode": "54"
  },
  {
    "regioncn": "亚美尼亚",
    "regionen": "Armenia",
    "domain": "AM",
    "telecode": "374"
  },
  {
    "regioncn": "阿鲁巴",
    "regionen": "Aruba",
    "domain": "AW",
    "telecode": "297"
  },
  {
    "regioncn": "澳大利亚",
    "regionen": "Australia",
    "domain": "AU",
    "telecode": "61"
  },
  {
    "regioncn": "奥地利",
    "regionen": "Austria",
    "domain": "AT",
    "telecode": "43"
  },
  {
    "regioncn": "阿塞拜疆",
    "regionen": "Azerbaijan",
    "domain": "AZ",
    "telecode": "994"
  },
  {
    "regioncn": "巴哈马",
    "regionen": "Bahamas",
    "domain": "BS",
    "telecode": "1242"
  },
  {
    "regioncn": "巴林",
    "regionen": "Bahrain",
    "domain": "BH",
    "telecode": "973"
  },
  {
    "regioncn": "孟加拉国",
    "regionen": "Bangladesh",
    "domain": "BD",
    "telecode": "880"
  },
  {
    "regioncn": "巴巴多斯",
    "regionen": "Barbados",
    "domain": "BB",
    "telecode": "1246"
  },
  {
    "regioncn": "白俄罗斯",
    "regionen": "Belarus",
    "domain": "BY",
    "telecode": "375"
  },
  {
    "regioncn": "比利时",
    "regionen": "Belgium",
    "domain": "BE",
    "telecode": "32"
  },
  {
    "regioncn": "伯利兹",
    "regionen": "Belize",
    "domain": "BZ",
    "telecode": "501"
  },
  {
    "regioncn": "贝宁",
    "regionen": "Benin",
    "domain": "BJ",
    "telecode": "229"
  },
  {
    "regioncn": "百慕大",
    "regionen": "Bermuda",
    "domain": "BM",
    "telecode": "1441"
  },
  {
    "regioncn": "不丹",
    "regionen": "Bhutan",
    "domain": "BT",
    "telecode": "975"
  },
  {
    "regioncn": "玻利维亚",
    "regionen": "Bolivia",
    "domain": "BO",
    "telecode": "591"
  },
  {
    "regioncn": "波黑",
    "regionen": "Bosnia and Herzegovina",
    "domain": "BA",
    "telecode": "387"
  },
  {
    "regioncn": "博茨瓦纳",
    "regionen": "Botswana",
    "domain": "BW",
    "telecode": "267"
  },
  {
    "regioncn": "布维岛",
    "regionen": "Bouvet Island",
    "domain": "BV",
    "telecode": "47"
  },
  {
    "regioncn": "巴西",
    "regionen": "Brazil",
    "domain": "BR",
    "telecode": "55"
  },
  {
    "regioncn": "英属印度洋领地",
    "regionen": "British Indian Ocean Territory",
    "domain": "IO",
    "telecode": "246"
  },
  {
    "regioncn": "文莱",
    "regionen": "Brunei Darussalam",
    "domain": "BN",
    "telecode": "673"
  },
  {
    "regioncn": "保加利亚",
    "regionen": "Bulgaria",
    "domain": "BG",
    "telecode": "359"
  },
  {
    "regioncn": "布基纳法索",
    "regionen": "Burkina Faso",
    "domain": "BF",
    "telecode": "226"
  },
  {
    "regioncn": "布隆迪",
    "regionen": "Burundi",
    "domain": "BI",
    "telecode": "257"
  },
  {
    "regioncn": "柬埔寨",
    "regionen": "Cambodia",
    "domain": "KH",
    "telecode": "855"
  },
  {
    "regioncn": "喀麦隆",
    "regionen": "Cameroon",
    "domain": "CM",
    "telecode": "237"
  },
  {
    "regioncn": "加拿大",
    "regionen": "Canada",
    "domain": "CA",
    "telecode": "1"
  },
  {
    "regioncn": "佛得角",
    "regionen": "Cape Verde",
    "domain": "CV",
    "telecode": "238"
  },
  {
    "regioncn": "开曼群岛",
    "regionen": "Cayman Islands",
    "domain": "KY",
    "telecode": "1345"
  },
  {
    "regioncn": "中非",
    "regionen": "Central African Republic",
    "domain": "CF",
    "telecode": "236"
  },
  {
    "regioncn": "乍得",
    "regionen": "Chad",
    "domain": "TD",
    "telecode": "235"
  },
  {
    "regioncn": "智利",
    "regionen": "Chile",
    "domain": "CL",
    "telecode": "56"
  },
  {
    "regioncn": "中国",
    "regionen": "China",
    "domain": "CN",
    "telecode": "86"
  },
  {
    "regioncn": "圣诞岛",
    "regionen": "Christmas Island",
    "domain": "CX",
    "telecode": "61"
  },
  {
    "regioncn": "科科斯（基林）群岛",
    "regionen": "Cocos (Keeling) Islands",
    "domain": "CC",
    "telecode": "61"
  },
  {
    "regioncn": "哥伦比亚",
    "regionen": "Colombia",
    "domain": "CO",
    "telecode": "57"
  },
  {
    "regioncn": "科摩罗",
    "regionen": "Comoros",
    "domain": "KM",
    "telecode": "269"
  },
  {
    "regioncn": "刚果（布）",
    "regionen": "Congo",
    "domain": "CG",
    "telecode": "243"
  },
  {
    "regioncn": "刚果（金）",
    "regionen": "Congo",
    "domain": "CD",
    "telecode": "242"
  },
  {
    "regioncn": "库克群岛",
    "regionen": "Cook Islands",
    "domain": "CK",
    "telecode": "682"
  },
  {
    "regioncn": "哥斯达黎加",
    "regionen": "Costa Rica",
    "domain": "CR",
    "telecode": "506"
  },
  {
    "regioncn": "科特迪瓦",
    "regionen": "Côte d'Ivoire",
    "domain": "CI",
    "telecode": "225"
  },
  {
    "regioncn": "克罗地亚",
    "regionen": "Croatia",
    "domain": "HR",
    "telecode": "385"
  },
  {
    "regioncn": "古巴",
    "regionen": "Cuba",
    "domain": "CU",
    "telecode": "53"
  },
  {
    "regioncn": "塞浦路斯",
    "regionen": "Cyprus",
    "domain": "CY",
    "telecode": "357"
  },
  {
    "regioncn": "捷克",
    "regionen": "Czech Republic",
    "domain": "CZ",
    "telecode": "420"
  },
  {
    "regioncn": "丹麦",
    "regionen": "Denmark",
    "domain": "DK",
    "telecode": "45"
  },
  {
    "regioncn": "吉布提",
    "regionen": "Djibouti",
    "domain": "DJ",
    "telecode": "253"
  },
  {
    "regioncn": "多米尼克",
    "regionen": "Dominica",
    "domain": "DM",
    "telecode": "1767"
  },
  {
    "regioncn": "多米尼加",
    "regionen": "Dominican Republic",
    "domain": "DO",
    "telecode": "1809"
  },
  {
    "regioncn": "厄瓜多尔",
    "regionen": "Ecuador",
    "domain": "EC",
    "telecode": "593"
  },
  {
    "regioncn": "埃及",
    "regionen": "Egypt",
    "domain": "EG",
    "telecode": "20"
  },
  {
    "regioncn": "萨尔瓦多",
    "regionen": "El Salvador",
    "domain": "SV",
    "telecode": "503"
  },
  {
    "regioncn": "赤道几内亚",
    "regionen": "Equatorial Guinea",
    "domain": "GQ",
    "telecode": "240"
  },
  {
    "regioncn": "厄立特里亚",
    "regionen": "Eritrea",
    "domain": "ER",
    "telecode": "291"
  },
  {
    "regioncn": "爱沙尼亚",
    "regionen": "Estonia",
    "domain": "EE",
    "telecode": "372"
  },
  {
    "regioncn": "埃塞俄比亚",
    "regionen": "Ethiopia",
    "domain": "ET",
    "telecode": "251"
  },
  {
    "regioncn": "福克兰群岛",
    "regionen": "Falkland Islands",
    "domain": "FK",
    "telecode": "500"
  },
  {
    "regioncn": "法罗群岛",
    "regionen": "Faroe Islands",
    "domain": "FO",
    "telecode": "298"
  },
  {
    "regioncn": "斐济",
    "regionen": "Fiji",
    "domain": "FJ",
    "telecode": "679"
  },
  {
    "regioncn": "芬兰",
    "regionen": "Finland",
    "domain": "FI",
    "telecode": "358"
  },
  {
    "regioncn": "法国",
    "regionen": "France",
    "domain": "FR",
    "telecode": "33"
  },
  {
    "regioncn": "法属圭亚那",
    "regionen": "French Guiana",
    "domain": "GF",
    "telecode": "594"
  },
  {
    "regioncn": "法属波利尼西亚",
    "regionen": "French Polynesia",
    "domain": "PF",
    "telecode": "689"
  },
  {
    "regioncn": "法属南部领地",
    "regionen": "French Southern Territories",
    "domain": "TF",
    "telecode": "-"
  },
  {
    "regioncn": "加蓬",
    "regionen": "Gabon",
    "domain": "GA",
    "telecode": "241"
  },
  {
    "regioncn": "冈比亚",
    "regionen": "Gambia",
    "domain": "GM",
    "telecode": "220"
  },
  {
    "regioncn": "格鲁吉亚",
    "regionen": "Georgia",
    "domain": "GE",
    "telecode": "995"
  },
  {
    "regioncn": "德国",
    "regionen": "Germany",
    "domain": "DE",
    "telecode": "49"
  },
  {
    "regioncn": "加纳",
    "regionen": "Ghana",
    "domain": "GH",
    "telecode": "233"
  },
  {
    "regioncn": "直布罗陀",
    "regionen": "Gibraltar",
    "domain": "GI",
    "telecode": "350"
  },
  {
    "regioncn": "希腊",
    "regionen": "Greece",
    "domain": "GR",
    "telecode": "30"
  },
  {
    "regioncn": "格陵兰",
    "regionen": "Greenland",
    "domain": "GL",
    "telecode": "299"
  },
  {
    "regioncn": "格林纳达",
    "regionen": "Grenada",
    "domain": "GD",
    "telecode": "1473"
  },
  {
    "regioncn": "瓜德罗普",
    "regionen": "Guadeloupe",
    "domain": "GP",
    "telecode": "590"
  },
  {
    "regioncn": "关岛",
    "regionen": "Guam",
    "domain": "GU",
    "telecode": "1671"
  },
  {
    "regioncn": "危地马拉",
    "regionen": "Guatemala",
    "domain": "GT",
    "telecode": "502"
  },
  {
    "regioncn": "格恩西岛",
    "regionen": "Guernsey",
    "domain": "GG",
    "telecode": "44"
  },
  {
    "regioncn": "几内亚",
    "regionen": "Guinea",
    "domain": "GN",
    "telecode": "224"
  },
  {
    "regioncn": "几内亚比绍",
    "regionen": "Guinea-Bissau",
    "domain": "GW",
    "telecode": "245"
  },
  {
    "regioncn": "圭亚那",
    "regionen": "Guyana",
    "domain": "GY",
    "telecode": "592"
  },
  {
    "regioncn": "海地",
    "regionen": "Haiti",
    "domain": "HT",
    "telecode": "509"
  },
  {
    "regioncn": "赫德岛和麦克唐纳岛",
    "regionen": "Heard Island and McDonald Islands",
    "domain": "HM",
    "telecode": "1672"
  },
  {
    "regioncn": "梵蒂冈",
    "regionen": "Holy See",
    "domain": "VA",
    "telecode": "379"
  },
  {
    "regioncn": "洪都拉斯",
    "regionen": "Honduras",
    "domain": "HN",
    "telecode": "504"
  },
  {
    "regioncn": "香港",
    "regionen": "Hong Kong",
    "domain": "HK",
    "telecode": "852"
  },
  {
    "regioncn": "匈牙利",
    "regionen": "Hungary",
    "domain": "HU",
    "telecode": "36"
  },
  {
    "regioncn": "冰岛",
    "regionen": "Iceland",
    "domain": "IS",
    "telecode": "354"
  },
  {
    "regioncn": "印度",
    "regionen": "India",
    "domain": "IN",
    "telecode": "91"
  },
  {
    "regioncn": "印度尼西亚",
    "regionen": "Indonesia",
    "domain": "ID",
    "telecode": "62"
  },
  {
    "regioncn": "伊朗",
    "regionen": "Iran",
    "domain": "IR",
    "telecode": "98"
  },
  {
    "regioncn": "伊拉克",
    "regionen": "Iraq",
    "domain": "IQ",
    "telecode": "964"
  },
  {
    "regioncn": "爱尔兰",
    "regionen": "Ireland",
    "domain": "IE",
    "telecode": "353"
  },
  {
    "regioncn": "英国属地曼岛",
    "regionen": "Isle of Man",
    "domain": "IM",
    "telecode": "44"
  },
  {
    "regioncn": "以色列",
    "regionen": "Israel",
    "domain": "IL",
    "telecode": "972"
  },
  {
    "regioncn": "意大利",
    "regionen": "Italy",
    "domain": "IT",
    "telecode": "39"
  },
  {
    "regioncn": "牙买加",
    "regionen": "Jamaica",
    "domain": "JM",
    "telecode": "1876"
  },
  {
    "regioncn": "日本",
    "regionen": "Japan",
    "domain": "JP",
    "telecode": "81"
  },
  {
    "regioncn": "泽西岛",
    "regionen": "Jersey",
    "domain": "JE",
    "telecode": "44"
  },
  {
    "regioncn": "约旦",
    "regionen": "Jordan",
    "domain": "JO",
    "telecode": "962"
  },
  {
    "regioncn": "哈萨克斯坦",
    "regionen": "Kazakhstan",
    "domain": "KZ",
    "telecode": "7"
  },
  {
    "regioncn": "肯尼亚",
    "regionen": "Kenya",
    "domain": "KE",
    "telecode": "254"
  },
  {
    "regioncn": "基里巴斯",
    "regionen": "Kiribati",
    "domain": "KI",
    "telecode": "686"
  },
  {
    "regioncn": "朝鲜",
    "regionen": "Korea",
    "domain": "KP",
    "telecode": "850"
  },
  {
    "regioncn": "韩国",
    "regionen": "Korea",
    "domain": "KR",
    "telecode": "82"
  },
  {
    "regioncn": "科威特",
    "regionen": "Kuwait",
    "domain": "KW",
    "telecode": "965"
  },
  {
    "regioncn": "吉尔吉斯斯坦",
    "regionen": "Kyrgyzstan",
    "domain": "KG",
    "telecode": "996"
  },
  {
    "regioncn": "老挝",
    "regionen": "Lao People's Democratic Republic",
    "domain": "LA",
    "telecode": "856"
  },
  {
    "regioncn": "拉脱维亚",
    "regionen": "Latvia",
    "domain": "LV",
    "telecode": "371"
  },
  {
    "regioncn": "黎巴嫩",
    "regionen": "Lebanon",
    "domain": "LB",
    "telecode": "961"
  },
  {
    "regioncn": "莱索托",
    "regionen": "Lesotho",
    "domain": "LS",
    "telecode": "266"
  },
  {
    "regioncn": "利比里亚",
    "regionen": "Liberia",
    "domain": "LR",
    "telecode": "231"
  },
  {
    "regioncn": "利比亚",
    "regionen": "Libyan Arab Jamahiriya",
    "domain": "LY",
    "telecode": "218"
  },
  {
    "regioncn": "列支敦士登",
    "regionen": "Liechtenstein",
    "domain": "LI",
    "telecode": "423"
  },
  {
    "regioncn": "立陶宛",
    "regionen": "Lithuania",
    "domain": "LT",
    "telecode": "370"
  },
  {
    "regioncn": "卢森堡",
    "regionen": "Luxembourg",
    "domain": "LU",
    "telecode": "352"
  },
  {
    "regioncn": "澳门",
    "regionen": "Macao",
    "domain": "MO",
    "telecode": "853"
  },
  {
    "regioncn": "北马其顿",
    "regionen": "Macedonia",
    "domain": "MK",
    "telecode": "389"
  },
  {
    "regioncn": "马达加斯加",
    "regionen": "Madagascar",
    "domain": "MG",
    "telecode": "261"
  },
  {
    "regioncn": "马拉维",
    "regionen": "Malawi",
    "domain": "MW",
    "telecode": "265"
  },
  {
    "regioncn": "马来西亚",
    "regionen": "Malaysia",
    "domain": "MY",
    "telecode": "60"
  },
  {
    "regioncn": "马尔代夫",
    "regionen": "Maldives",
    "domain": "MV",
    "telecode": "960"
  },
  {
    "regioncn": "马里",
    "regionen": "Mali",
    "domain": "ML",
    "telecode": "223"
  },
  {
    "regioncn": "马耳他",
    "regionen": "Malta",
    "domain": "MT",
    "telecode": "356"
  },
  {
    "regioncn": "马绍尔群岛",
    "regionen": "Marshall Islands",
    "domain": "MH",
    "telecode": "692"
  },
  {
    "regioncn": "马提尼克",
    "regionen": "Martinique",
    "domain": "MQ",
    "telecode": "596"
  },
  {
    "regioncn": "毛利塔尼亚",
    "regionen": "Mauritania",
    "domain": "MR",
    "telecode": "222"
  },
  {
    "regioncn": "毛里求斯",
    "regionen": "Mauritius",
    "domain": "MU",
    "telecode": "230"
  },
  {
    "regioncn": "马约特",
    "regionen": "Mayotte",
    "domain": "YT",
    "telecode": "262"
  },
  {
    "regioncn": "墨西哥",
    "regionen": "Mexico",
    "domain": "MX",
    "telecode": "52"
  },
  {
    "regioncn": "密克罗尼西亚联邦",
    "regionen": "Micronesia",
    "domain": "FM",
    "telecode": "691"
  },
  {
    "regioncn": "摩尔多瓦",
    "regionen": "Moldova",
    "domain": "MD",
    "telecode": "373"
  },
  {
    "regioncn": "摩纳哥",
    "regionen": "Monaco",
    "domain": "MC",
    "telecode": "377"
  },
  {
    "regioncn": "蒙古",
    "regionen": "Mongolia",
    "domain": "MN",
    "telecode": "976"
  },
  {
    "regioncn": "黑山",
    "regionen": "Montenegro",
    "domain": "ME",
    "telecode": "382"
  },
  {
    "regioncn": "蒙特塞拉特",
    "regionen": "Montserrat",
    "domain": "MS",
    "telecode": "1664"
  },
  {
    "regioncn": "摩洛哥",
    "regionen": "Morocco",
    "domain": "MA",
    "telecode": "212"
  },
  {
    "regioncn": "莫桑比克",
    "regionen": "Mozambique",
    "domain": "MZ",
    "telecode": "258"
  },
  {
    "regioncn": "缅甸",
    "regionen": "Myanmar",
    "domain": "MM",
    "telecode": "95"
  },
  {
    "regioncn": "纳米比亚",
    "regionen": "Namibia",
    "domain": "NA",
    "telecode": "264"
  },
  {
    "regioncn": "瑙鲁",
    "regionen": "Nauru",
    "domain": "NR",
    "telecode": "674"
  },
  {
    "regioncn": "尼泊尔",
    "regionen": "Nepal",
    "domain": "NP",
    "telecode": "977"
  },
  {
    "regioncn": "荷兰",
    "regionen": "Netherlands",
    "domain": "NL",
    "telecode": "31"
  },
  {
    "regioncn": "荷属安的列斯",
    "regionen": "Netherlands Antilles",
    "domain": "AN",
    "telecode": ""
  },
  {
    "regioncn": "新喀里多尼亚",
    "regionen": "New Caledonia",
    "domain": "NC",
    "telecode": "687"
  },
  {
    "regioncn": "新西兰",
    "regionen": "New Zealand",
    "domain": "NZ",
    "telecode": "64"
  },
  {
    "regioncn": "尼加拉瓜",
    "regionen": "Nicaragua",
    "domain": "NI",
    "telecode": "505"
  },
  {
    "regioncn": "尼日尔",
    "regionen": "Niger",
    "domain": "NE",
    "telecode": "227"
  },
  {
    "regioncn": "尼日利亚",
    "regionen": "Nigeria",
    "domain": "NG",
    "telecode": "234"
  },
  {
    "regioncn": "纽埃",
    "regionen": "Niue",
    "domain": "NU",
    "telecode": "683"
  },
  {
    "regioncn": "诺福克岛",
    "regionen": "Norfolk Island",
    "domain": "NF",
    "telecode": "6723"
  },
  {
    "regioncn": "北马里亚纳",
    "regionen": "Northern Mariana Islands",
    "domain": "MP",
    "telecode": "1670"
  },
  {
    "regioncn": "挪威",
    "regionen": "Norway",
    "domain": "NO",
    "telecode": "47"
  },
  {
    "regioncn": "阿曼",
    "regionen": "Oman",
    "domain": "OM",
    "telecode": "968"
  },
  {
    "regioncn": "巴基斯坦",
    "regionen": "Pakistan",
    "domain": "PK",
    "telecode": "92"
  },
  {
    "regioncn": "帕劳",
    "regionen": "Palau",
    "domain": "PW",
    "telecode": "680"
  },
  {
    "regioncn": "巴勒斯坦",
    "regionen": "Palestinian Territory",
    "domain": "PS",
    "telecode": "970"
  },
  {
    "regioncn": "巴拿马",
    "regionen": "Panama",
    "domain": "PA",
    "telecode": "507"
  },
  {
    "regioncn": "巴布亚新几内亚",
    "regionen": "Papua New Guinea",
    "domain": "PG",
    "telecode": "675"
  },
  {
    "regioncn": "巴拉圭",
    "regionen": "Paraguay",
    "domain": "PY",
    "telecode": "595"
  },
  {
    "regioncn": "秘鲁",
    "regionen": "Peru",
    "domain": "PE",
    "telecode": "51"
  },
  {
    "regioncn": "菲律宾",
    "regionen": "Philippines",
    "domain": "PH",
    "telecode": "63"
  },
  {
    "regioncn": "皮特凯恩",
    "regionen": "Pitcairn",
    "domain": "PN",
    "telecode": "64"
  },
  {
    "regioncn": "波兰",
    "regionen": "Poland",
    "domain": "PL",
    "telecode": "48"
  },
  {
    "regioncn": "葡萄牙",
    "regionen": "Portugal",
    "domain": "PT",
    "telecode": "351"
  },
  {
    "regioncn": "波多黎各",
    "regionen": "Puerto Rico",
    "domain": "PR",
    "telecode": "1787"
  },
  {
    "regioncn": "卡塔尔",
    "regionen": "Qatar",
    "domain": "QA",
    "telecode": "974"
  },
  {
    "regioncn": "留尼汪",
    "regionen": "Réunion",
    "domain": "RE",
    "telecode": "262"
  },
  {
    "regioncn": "罗马尼亚",
    "regionen": "Romania",
    "domain": "RO",
    "telecode": "40"
  },
  {
    "regioncn": "俄罗斯联邦",
    "regionen": "Russian Federation",
    "domain": "RU",
    "telecode": "7"
  },
  {
    "regioncn": "卢旺达",
    "regionen": "Rwanda",
    "domain": "RW",
    "telecode": "250"
  },
  {
    "regioncn": "圣赫勒拿",
    "regionen": "Saint Helena",
    "domain": "SH",
    "telecode": "290"
  },
  {
    "regioncn": "圣基茨和尼维斯",
    "regionen": "Saint Kitts and Nevis",
    "domain": "KN",
    "telecode": "1869"
  },
  {
    "regioncn": "圣卢西亚",
    "regionen": "Saint Lucia",
    "domain": "LC",
    "telecode": "1758"
  },
  {
    "regioncn": "圣皮埃尔和密克隆",
    "regionen": "Saint Pierre and Miquelon",
    "domain": "PM",
    "telecode": "508"
  },
  {
    "regioncn": "圣文森特和格林纳丁斯",
    "regionen": "Saint Vincent and the Grenadines",
    "domain": "VC",
    "telecode": "1784"
  },
  {
    "regioncn": "萨摩亚",
    "regionen": "Samoa",
    "domain": "WS",
    "telecode": "685"
  },
  {
    "regioncn": "圣马力诺",
    "regionen": "San Marino",
    "domain": "SM",
    "telecode": "378"
  },
  {
    "regioncn": "圣多美和普林西比",
    "regionen": "Sao Tome and Principe",
    "domain": "ST",
    "telecode": "239"
  },
  {
    "regioncn": "沙特阿拉伯",
    "regionen": "Saudi Arabia",
    "domain": "SA",
    "telecode": "966"
  },
  {
    "regioncn": "塞内加尔",
    "regionen": "Senegal",
    "domain": "SN",
    "telecode": "221"
  },
  {
    "regioncn": "塞尔维亚",
    "regionen": "Serbia",
    "domain": "RS",
    "telecode": "381"
  },
  {
    "regioncn": "塞舌尔",
    "regionen": "Seychelles",
    "domain": "SC",
    "telecode": "248"
  },
  {
    "regioncn": "塞拉利昂",
    "regionen": "Sierra Leone",
    "domain": "SL",
    "telecode": "232"
  },
  {
    "regioncn": "新加坡",
    "regionen": "Singapore",
    "domain": "SG",
    "telecode": "65"
  },
  {
    "regioncn": "斯洛伐克",
    "regionen": "Slovakia",
    "domain": "SK",
    "telecode": "421"
  },
  {
    "regioncn": "斯洛文尼亚",
    "regionen": "Slovenia",
    "domain": "SI",
    "telecode": "386"
  },
  {
    "regioncn": "所罗门群岛",
    "regionen": "Solomon Islands",
    "domain": "SB",
    "telecode": "677"
  },
  {
    "regioncn": "索马里",
    "regionen": "Somalia",
    "domain": "SO",
    "telecode": "252"
  },
  {
    "regioncn": "南非",
    "regionen": "South Africa",
    "domain": "ZA",
    "telecode": "27"
  },
  {
    "regioncn": "南乔治亚岛和南桑德韦奇岛",
    "regionen": "South Georgia and the South Sandwich Islands",
    "domain": "GS",
    "telecode": "500"
  },
  {
    "regioncn": "西班牙",
    "regionen": "Spain",
    "domain": "ES",
    "telecode": "34"
  },
  {
    "regioncn": "斯里兰卡",
    "regionen": "Sri Lanka",
    "domain": "LK",
    "telecode": "94"
  },
  {
    "regioncn": "苏丹",
    "regionen": "Sudan",
    "domain": "SD",
    "telecode": "249"
  },
  {
    "regioncn": "苏里南",
    "regionen": "Suriname",
    "domain": "SR",
    "telecode": "597"
  },
  {
    "regioncn": "斯瓦尔巴岛和扬马延岛",
    "regionen": "Svalbard and Jan Mayen",
    "domain": "SJ",
    "telecode": "47"
  },
  {
    "regioncn": "斯威士兰",
    "regionen": "Swaziland",
    "domain": "SZ",
    "telecode": "268"
  },
  {
    "regioncn": "瑞典",
    "regionen": "Sweden",
    "domain": "SE",
    "telecode": "46"
  },
  {
    "regioncn": "瑞士",
    "regionen": "Switzerland",
    "domain": "CH",
    "telecode": "41"
  },
  {
    "regioncn": "叙利亚",
    "regionen": "Syrian Arab Republic",
    "domain": "SY",
    "telecode": "963"
  },
  {
    "regioncn": "台湾",
    "regionen": "Taiwan,Province of China",
    "domain": "TW",
    "telecode": "886"
  },
  {
    "regioncn": "塔吉克斯坦",
    "regionen": "Tajikistan",
    "domain": "TJ",
    "telecode": "992"
  },
  {
    "regioncn": "坦桑尼亚",
    "regionen": "Tanzania,United Republic of",
    "domain": "TZ",
    "telecode": "255"
  },
  {
    "regioncn": "泰国",
    "regionen": "Thailand",
    "domain": "TH",
    "telecode": "66"
  },
  {
    "regioncn": "东帝汶",
    "regionen": "Timor-Leste",
    "domain": "TL",
    "telecode": "670"
  },
  {
    "regioncn": "多哥",
    "regionen": "Togo",
    "domain": "TG",
    "telecode": "228"
  },
  {
    "regioncn": "托克劳",
    "regionen": "Tokelau",
    "domain": "TK",
    "telecode": "690"
  },
  {
    "regioncn": "汤加",
    "regionen": "Tonga",
    "domain": "TO",
    "telecode": "676"
  },
  {
    "regioncn": "特立尼达和多巴哥",
    "regionen": "Trinidad and Tobago",
    "domain": "TT",
    "telecode": "1868"
  },
  {
    "regioncn": "突尼斯",
    "regionen": "Tunisia",
    "domain": "TN",
    "telecode": "216"
  },
  {
    "regioncn": "土耳其",
    "regionen": "Turkey",
    "domain": "TR",
    "telecode": ""
  },
  {
    "regioncn": "土库曼斯坦",
    "regionen": "Turkmenistan",
    "domain": "TM",
    "telecode": "993"
  },
  {
    "regioncn": "特克斯和凯科斯群岛",
    "regionen": "Turks and Caicos Islands",
    "domain": "TC",
    "telecode": "1649"
  },
  {
    "regioncn": "图瓦卢",
    "regionen": "Tuvalu",
    "domain": "TV",
    "telecode": "688"
  },
  {
    "regioncn": "乌干达",
    "regionen": "Uganda",
    "domain": "UG",
    "telecode": "256"
  },
  {
    "regioncn": "乌克兰",
    "regionen": "Ukraine",
    "domain": "UA",
    "telecode": "380"
  },
  {
    "regioncn": "阿联酋",
    "regionen": "United Arab Emirates",
    "domain": "AE",
    "telecode": "971"
  },
  {
    "regioncn": "英国",
    "regionen": "United Kingdom",
    "domain": "GB",
    "telecode": "44"
  },
  {
    "regioncn": "美国",
    "regionen": "United States",
    "domain": "US",
    "telecode": "1"
  },
  {
    "regioncn": "美国本土外小岛屿",
    "regionen": "United States Minor Outlying Islands",
    "domain": "UM",
    "telecode": "1808"
  },
  {
    "regioncn": "乌拉圭",
    "regionen": "Uruguay",
    "domain": "UY",
    "telecode": "598"
  },
  {
    "regioncn": "乌兹别克斯坦",
    "regionen": "Uzbekistan",
    "domain": "UZ",
    "telecode": "998"
  },
  {
    "regioncn": "瓦努阿图",
    "regionen": "Vanuatu",
    "domain": "VU",
    "telecode": "678"
  },
  {
    "regioncn": "委内瑞拉",
    "regionen": "Venezuela",
    "domain": "VE",
    "telecode": "58"
  },
  {
    "regioncn": "越南",
    "regionen": "Viet Nam",
    "domain": "VN",
    "telecode": "84"
  },
  {
    "regioncn": "英属维尔京群岛",
    "regionen": "Virgin Islands (British)",
    "domain": "VG",
    "telecode": "1284"
  },
  {
    "regioncn": "美属维尔京群岛",
    "regionen": "Virgin Islands (U.S.)",
    "domain": "VI",
    "telecode": "1340"
  },
  {
    "regioncn": "瓦利斯和富图纳",
    "regionen": "Wallis and Futuna",
    "domain": "WF",
    "telecode": "681"
  },
  {
    "regioncn": "西撒哈拉",
    "regionen": "Western Sahara",
    "domain": "EH",
    "telecode": "21228"
  },
  {
    "regioncn": "也门",
    "regionen": "Yemen",
    "domain": "YE",
    "telecode": "967"
  },
  {
    "regioncn": "赞比亚",
    "regionen": "Zambia",
    "domain": "ZM",
    "telecode": "260"
  },
  {
    "regioncn": "津巴布韦",
    "regionen": "Zimbabwe",
    "domain": "ZW",
    "telecode": "263"
  }
]`
