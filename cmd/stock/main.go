package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/sirupsen/logrus"

	"github.com/vikashvverma/stock-backend/auth"
	"github.com/vikashvverma/stock-backend/config"
	"github.com/vikashvverma/stock-backend/factory"
	"github.com/vikashvverma/stock-backend/log"
	"github.com/vikashvverma/stock-backend/router"
)

var (
	version      string
	publicRoutes = []string{"/healthcheck", "/version"}
)

func main() {
	var c *config.Config
	var err error

	useFlags := !(len(os.Args) > 2 && os.Args[1] == "-config")
	if useFlags {
		c, err = config.FromFlags(os.Args)
		if err != nil {
			logrus.Fatalln(err)
		}
	} else {
		c, err = config.FromFile(os.Args[2])
		if err != nil {
			logrus.Fatalln(err)
		}
	}
	l := logrus.New()
	f := factory.NewFactory(c, l)
	muxRouter := router.Router(f, c, l)

	n := negroni.New()
	n.Use(log.New(l))
	n.Use(auth.New(l, c.APIKey, publicRoutes))
	n.UseHandler(muxRouter)
	n.Run(fmt.Sprintf(":%d", c.AppPort()))
}
