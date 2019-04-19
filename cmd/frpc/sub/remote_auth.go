package sub

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/fatedier/frp/utils/util"

	"gopkg.in/ini.v1"
)

// reqAuthByKey 请求服务端配置
func reqAuthByKey() ([]byte, string) {
	bytesURL, err := base64.RawURLEncoding.DecodeString(authKey)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	strURL := string(bytesURL)

	u, err := url.Parse(strURL)
	if err != nil {
		panic(err)
	}

	res, err := client.Get(strURL)
	if err != nil {
		panic(err)
	}
	confData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return confData, filepath.Base(u.Path)
}

// startByRemoteCfg 根据服务端配置启动服务
func startByRemoteCfg() {
	confData, key := reqAuthByKey()
	confData = util.AESCFBDecrypter(key, confData)
	f, err := ini.Load(confData)
	if err != nil {
		panic(err)
	}
	for _, s := range f.Sections() {
		if s.Name() == "common" {
			if !s.HasKey("user") {
				s.NewKey("user", key)
			}
		}
	}
	fw, err := ioutil.TempFile("", "config")
	if err != nil {
		panic(err)
	}
	f.WriteTo(fw)
	fw.Close()
	log.Printf("write temp config file %s", fw.Name())
	if err := runClient(fw.Name()); err != nil {
		panic(err)
	}
}
