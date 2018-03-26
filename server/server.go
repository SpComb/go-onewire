package server

import (
	"github.com/SpComb/go-onewire/api"
	"github.com/SpComb/go-onewire/netlink/connector/w1"

	"fmt"
	"sync"
	"time"
)

const RefreshInterval = 1 * time.Second

func NewServer() (*Server, error) {
	var server = makeServer()

	if queryConn, err := w1.Dial(); err != nil {
		return nil, fmt.Errorf("w1.Dial: %v", err)
	} else {
		server.queryConn = queryConn
	}

	if eventConn, err := w1.Dial(); err != nil {
		return nil, fmt.Errorf("w1.Dial: %v", err)
	} else {
		server.eventConn = eventConn
	}

	return &server, nil
}

func makeServer() Server {
	return Server{
		refreshInterval: RefreshInterval,

		eventChan: make(chan w1.Event),
		sensors:   make(map[api.SensorID]*Sensor),
	}
}

type Server struct {
	refreshInterval time.Duration

	queryConn *w1.Conn
	eventConn *w1.Conn

	eventChan chan w1.Event

	sensors      map[api.SensorID]*Sensor
	sensorsMutex sync.RWMutex
}

func (s *Server) listenEvents() {
	defer close(s.eventChan)

	err := s.eventConn.Listen(func(event w1.Event) error {
		s.eventChan <- event
		return nil
	})

	log.Errorf("Listen stopped: %v", err)
}

func (s *Server) scan() error {
	if masterList, err := s.queryConn.ListMasters(); err != nil {
		return fmt.Errorf("w1 ListMasters: %v", err)
	} else {
		for _, masterID := range masterList {
			if slaveList, err := s.queryConn.ListSlaves(masterID); err != nil {
				return fmt.Errorf("w1 ListSlaves %v: %v", masterID, err)
			} else {
				for _, slaveID := range slaveList {
					s.initSlave(slaveID)
				}
			}
		}

		return nil
	}
}

func (s *Server) run() error {
	var refreshChan = time.Tick(s.refreshInterval)

	for {
		select {
		case <-refreshChan:
			s.refresh()

		case event, ok := <-s.eventChan:
			if !ok {
				return nil
			}

			log.Debugf("event: %#v", event)

			switch event.Type {
			case w1.MsgTypeSlaveAdd:
				s.addSlave(event.SlaveID())
			case w1.MsgTypeSlaveRemove:
				s.removeSlave(event.SlaveID())
			}
		}
	}
}

func (s *Server) initSlave(id w1.SlaveID) {
	var sensor = newSensor(s, id)

	log.Infof("init slave %v", id)

	s.SetSensor(sensor)
}

func (s *Server) addSlave(id w1.SlaveID) {
	var sensor = newSensor(s, id)

	log.Infof("add slave %v", id)

	s.SetSensor(sensor)
}

func (s *Server) removeSlave(id w1.SlaveID) {
	if sensor := s.GetSensor(api.SensorID(id.String())); sensor != nil {
		log.Infof("remove slave %v", id)

		s.DeleteSensor(sensor)
	}
}

func (s *Server) refresh() {
	log.Debugf("refresh")

	s.WalkSensors(func(sensor *Sensor) {
		sensor.refresh()
	})
}

func (s *Server) GetSensor(id api.SensorID) *Sensor {
	s.sensorsMutex.RLock()
	defer s.sensorsMutex.RUnlock()

	return s.sensors[id]
}

func (s *Server) SetSensor(sensor *Sensor) {
	s.sensorsMutex.Lock()
	defer s.sensorsMutex.Unlock()

	s.sensors[sensor.id] = sensor
}

func (s *Server) DeleteSensor(sensor *Sensor) {
	s.sensorsMutex.Lock()
	defer s.sensorsMutex.Unlock()

	delete(s.sensors, sensor.id)
}

func (s *Server) WalkSensors(f func(sensor *Sensor)) {
	s.sensorsMutex.RLock()
	defer s.sensorsMutex.RUnlock()

	for _, sensor := range s.sensors {
		f(sensor)
	}
}

func (s *Server) Run() error {
	go s.listenEvents()

	if err := s.scan(); err != nil {
		return err
	}

	return s.run()
}

func (s *Server) MakeAPISensors() []api.Sensor {
	var apiList []api.Sensor

	s.WalkSensors(func(sensor *Sensor) {
		apiList = append(apiList, sensor.MakeAPI())
	})

	return apiList
}
