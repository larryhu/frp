package sub

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/fatedier/frp/utils/util"
	ini "github.com/vaughan0/go-ini"
)

// reqAuthByKey 请求服务端配置
func reqAuthByKey() ([]byte, string, error) {
	bytesURL, err := base64.RawURLEncoding.DecodeString(authKey)
	if err != nil {
		return nil, "", err
	}
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	strURL := string(bytesURL)

	u, err := url.Parse(strURL)
	if err != nil {
		return nil, "", err
	}

	res, err := client.Get(strURL)
	if err != nil {
		return nil, "", err
	}
	confData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}

	return confData, filepath.Base(u.Path), nil
}

// fetchRemoteCfg 根据服务端配置启动服务
func fetchRemoteCfg() (string, error) {
	confData, key, err := reqAuthByKey()
	if err != nil {
		return "", err
	}

	confData, err = util.AESCFBDecrypter(key, confData)
	if err != nil {
		return "", err
	}

	if _, err := ini.Load(bytes.NewBuffer(confData)); err != nil {
		return "", err
	}

	fw, err := ioutil.TempFile("", "config")
	if err != nil {
		return "", err
	}
	defer fw.Close()

	fw.Write(confData)

	return fw.Name(), nil
}
