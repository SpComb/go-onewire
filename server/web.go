package server

import (
	"github.com/SpComb/iot-poc/api"
	"github.com/qmsk/go-web"
)

func (s *Server) WebAPI() web.API {
	return web.MakeAPI(&serverView{s})
}

type serverView struct {
	server *Server
}

func (view *serverView) Index(name string) (web.Resource, error) {
	switch name {
	case "":
		return view, nil
	default:
		return nil, nil
	}
}

func (view *serverView) MakeAPI() api.Index {
	return api.Index{
		Sensors: view.server.MakeAPISensors(),
	}
}

func (view *serverView) GetREST() (web.Resource, error) {
	return view.MakeAPI(), nil
}
