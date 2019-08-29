package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fatedier/frp/utils/log"
	"github.com/fatedier/frp/utils/util"
	"github.com/gorilla/mux"
)

var (
	authAddr string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&authAddr, "auth_addr", "", "", "bind auth address :10080")
}

func authKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars == nil {
		w.WriteHeader(400)
		return
	}

	key := vars["key"]

	fr, err := os.Open(filepath.Join("keys", key))
	if err != nil {
		w.WriteHeader(404)
		return
	}
	defer fr.Close()

	configData, err := ioutil.ReadAll(fr)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	encryptData, err := util.AESCFBEncrypter(key, configData)
	if err != nil {
		println(err.Error())
		w.WriteHeader(417)
		return
	}

	w.WriteHeader(200)

	w.Header().Add("Content-Type", "application/octet-stream")
	w.Write(encryptData)
}

// 启动验证服务
func startAuthServer() {
	if authAddr == "" {
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/auth/key/{key}", authKey).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         authAddr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Info("web auth server at %s", authAddr)

	panic(srv.ListenAndServe())
}
