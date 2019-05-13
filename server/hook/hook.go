package hook

import (
	"fmt"
	"time"

	"github.com/fatedier/frp/g"
	"github.com/imroc/req"
)

// CheckUser check client user
func CheckUser(user string) error {
	c := req.New()
	c.SetTimeout(time.Second * 3)

	res, err := c.Get(fmt.Sprintf(g.GlbServerCfg.HookCheckUser, user))
	if err != nil {
		return fmt.Errorf("authorization user failed")
	}

	if res.Response().StatusCode != 200 {
		return fmt.Errorf("authorization user failed")
	}
	return nil
}
