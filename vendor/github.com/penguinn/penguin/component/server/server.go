package server

import (
	"context"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/jinzhu/gorm"
	"github.com/penguinn/penguin/component/db"
	"github.com/penguinn/penguin/component/middleware"
	"github.com/penguinn/penguin/component/router"
	"github.com/penguinn/penguin/component/session"
	"github.com/penguinn/penguin/constants"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func Serve() {
	BeforeStart()
	Start()
}

func BeforeStart() {
	session.SetCodec(session.NewJsonCodec())
	if s, err := session.NewDBStore(
		func() *gorm.DB {
			c, _ := db.Write(constants.DefaultDBName)
			return c
		}); err != nil {
		panic(err)
	} else {
		session.SetStore(s)
	}

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.DebugMiddleware())
}

func Start() {

	svrs := []*http.Server{}

	svrs = append(svrs, buildAppSrv())

	if router.GlobalRouterConfig.Pprof {

		addrDiv := strings.Split(router.GlobalRouterConfig.Addr, ":")
		p, _ := strconv.Atoi(addrDiv[1])
		svrs = append(svrs, buildPPROFSrv(":"+(strconv.Itoa(p+1))))

	}

	if router.GlobalRouterConfig.Mode == "release" {

		gracehttp.Serve(svrs...)

	} else {
		for _, s := range svrs {
			go func(s *http.Server) {
				if err := s.ListenAndServe(); err != nil {
					panic(fmt.Sprintf("Server Listen : %s \n", err))
				}
			}(s)
		}

		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		fmt.Println("Shutdown Server.")

		for _, s := range svrs {
			func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := s.Shutdown(ctx); err != nil {
					panic(fmt.Sprintf("Server Shutdown error[%v] \n", err))
				}
			}()

		}
		fmt.Println("Exist Normally.")
	}
}

func buildPPROFSrv(addr string) *http.Server {
	h := http.NewServeMux()
	h.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	h.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	h.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	h.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	h.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	return &http.Server{Addr: addr, Handler: h}
}

func buildAppSrv() *http.Server {
	app := &http.Server{}
	app.Addr = router.GlobalRouterConfig.Addr
	app.Handler = router.GlobalRouter
	return app
}
