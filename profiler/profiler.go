package profiler

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()
}
