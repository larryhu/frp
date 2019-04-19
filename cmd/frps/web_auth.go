package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatedier/frp/utils/util"
	"github.com/gin-gonic/gin"
)

func authKey(c *gin.Context) {
	key := c.Param("key")
	if len(key) != 32 {
		c.Status(400)
		return
	}

	fr, err := os.Open(filepath.Join("keys", key))
	if err != nil {
		c.Status(400)
		return
	}

	configData, err := ioutil.ReadAll(fr)
	if err != nil {
		c.Status(400)
		return
	}
	c.Data(200, "application/octet-stream", util.AESCFBEncrypter(key, configData))
}

// 启动验证服务
func startAuthServer() {
	w := gin.Default()
	w.GET("/auth/key/:key", authKey)

	if err := w.Run(authAddr); err != nil {
		panic(err)
	}
}
