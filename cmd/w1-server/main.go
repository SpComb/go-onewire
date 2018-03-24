package main

import (
	"github.com/SpComb/go-onewire/netlink/connector/w1"
	"github.com/SpComb/go-onewire/server"
	"github.com/qmsk/go-logging"
	"github.com/qmsk/go-web"

	"flag"
	"os"
)

var log logging.Logging

var options struct {
	Log       logging.Options
	LogW1     logging.Options
	LogServer logging.Options
	LogWeb    logging.Options

	Web web.Options
}

func init() {
	options.LogW1.Module = "w1"
	options.LogW1.Defaults = &options.Log
	options.LogServer.Module = "server"
	options.LogServer.Defaults = &options.Log
	options.LogWeb.Module = "web"
	options.LogWeb.Defaults = &options.Log

	options.Log.InitFlags()
	options.LogW1.InitFlags()
	options.LogServer.InitFlags()
	options.LogWeb.InitFlags()

	flag.StringVar(&options.Web.Listen, "http-listen", ":8286", "HTTP server listen: [HOST]:PORT")
	flag.StringVar(&options.Web.Static, "http-static", "", "HTTP sever /static path: PATH")

}

func webServer(server *server.Server) {
	err := options.Web.Server(
		options.Web.RouteAPI("/api/", server.WebAPI()),
		options.Web.RouteStatic("/"),
	)

	if err != nil {
		log.Errorf("web Server: %v", err)
		os.Exit(1)
	}
}

func run(server *server.Server) error {
	go webServer(server)

	return server.Run()
}

func main() {
	flag.Parse()

	log = options.Log.MakeLogging()

	w1.SetLogging(options.LogW1.MakeLogging())
	server.SetLogging(options.LogServer.MakeLogging())
	web.SetLogging(options.LogWeb.MakeLogging())

	if server, err := server.NewServer(); err != nil {
		log.Errorf("server.NewServer: %v", err)
		os.Exit(1)
	} else if err := run(server); err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}
}
