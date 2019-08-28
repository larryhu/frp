package hook

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatedier/frp/g"
	"github.com/fatedier/frp/utils/log"
)

var (
	errAuth = fmt.Errorf("authorization user failed")
)

// CheckUser check client user
func CheckUser(user string) error {
	if g.GlbServerCfg.HookCheckUser == "" {
		log.Info("skip hook_check_user")
		return nil
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	res, err := client.Get(fmt.Sprintf(g.GlbServerCfg.HookCheckUser, user))
	if err != nil {
		return errAuth
	}

	if res.StatusCode != 200 {
		return errAuth
	}
	return nil
}
