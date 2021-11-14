package profiler

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"techtrainingcamp-group3/config"
)

func init() {
	if config.Env.Profiler == "true" {
		go func() {
			log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}()
	}
}
