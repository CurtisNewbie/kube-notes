package main

import (
	"os"

	"github.com/curtisnewbie/gocommon/client"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
)

func main() {

	common.ScheduleCron("0/15 * * * * *", func() {
		ec := common.EmptyExecContext()

		client := client.NewDefaultTClient(ec.Ctx, "http://empty-mind:8081/ping")
		r := client.Get(map[string][]string{})
		if r.Err != nil {
			ec.Log.Errorf("Failed to ping empty-mind, %v", r.Err)
			return
		}
		rs, e := r.ReadStr()
		if e != nil {
			ec.Log.Errorf("Failed to read resp from empty-mind, %v", e)
			return
		}

		ec.Log.Infof("Ping(ed) empty-mind, r: %d, %s", r.Resp.StatusCode, rs)
	})

	server.DefaultBootstrapServer(os.Args)
}
