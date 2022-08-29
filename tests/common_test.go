package main

import (
	"net/url"
	"testing"
)

func TestGethttp(t *testing.T) {
	remoteUrl := "http://192.168.1.236:9000/install/do_conf"
	queryValues := make(url.Values)
	queryValues["address"] = []string{"0x01842a2CF56400A245a56955dc407c2C4137321e  "}
	msg := `'{"type":"exchange_auth","version":1,"exchanger_owner":"0x01842a2CF56400A245a56955dc407c2C4137321e","to":"0x7fBC8ad616177c6519228FCa4a7D9EC7d1804900","block_number":1000,"sig":"0xf3e46258b50ea68e112e120e72f09202111adbd4aed1318745c7a87559af1ce948f6fc25ca2d510b6c684f51f44078e92938674fe03a53815233c63b5129157b1c"}'`
	queryValues["params"] = []string{msg}

	_, _ = HttpGetParamsSendRev(remoteUrl, queryValues)
}
