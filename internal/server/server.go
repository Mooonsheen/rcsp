package server

import (
	"encoding/json"
	"log"
	cahce "rcsp/internal/cahce"
	"rcsp/internal/cahce/redis"
	"rcsp/internal/database"
	"rcsp/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

type Server struct {
	cache  cahce.Cache
	db     *database.Database
	config *config
	router *chi.Mux
	sc     stan.Conn
	sub    stan.Subscription
}

func setConfigs(path string) (*database.Database, *config, error) {
	db, err := database.SetConfig(path)
	if err != nil {
		return nil, nil, err
	}
	config, err := newConfig(path)
	if err != nil {
		return nil, nil, err
	}
	return db, config, nil
}

func NewServer(path string) (*Server, error) {
	db, config, err := setConfigs(path)
	if err != nil {
		return nil, err
	}
	return &Server{
		db:     db,
		cache:  redis.NewRedisClient(),
		config: config,
		router: chi.NewRouter(),
	}, nil
}

func (s *Server) connectToStream() error {
	sc, err := stan.Connect("test-cluster", "subscriber", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		return err
	}
	sub, err := sc.Subscribe(s.config.SubscribeSubject, s.handleRequest)
	if err != nil {
		return err
	}
	s.sc, s.sub = sc, sub
	return nil
}

func (s *Server) Up() error {
	s.db.Open()
	logrus.Info("Database is up")
	if err := s.setCache(); err != nil {
		log.Printf("can't set cache in server up, err: %e", err)
		return err
	}
	if err := s.connectToStream(); err != nil {
		log.Printf("can't connect to stream in server up, err: %e", err)
		return err
	}
	s.startRouting()
	return nil
}

func (s *Server) Down() {
	logrus.Info("Server is down")
	s.db.Close()
	s.sub.Unsubscribe()
	s.sc.Close()
}

func (s *Server) handleRequest(m *stan.Msg) {
	data := model.Order{}
	err := json.Unmarshal(m.Data, &data)
	if err != nil {
		return
	}
	if ok := s.addToCache(data); ok {
		logrus.Info("Add to cache")
		if err := s.db.AddOrder(data); err != nil {
			s.deleteFromCache(data)
			logrus.Info("Failed add db\n", err)
		}
	}
}

func (s *Server) setCache() error {
	orders := make([]model.Order, 0)
	err := s.db.DB.Model(&orders).Select()
	if err != nil {
		return err
	}
	for _, order := range orders {
		s.cache.AddOrder(order)
	}
	return nil
}

func (s *Server) addToCache(data model.Order) bool {
	err := s.cache.AddOrder(data)
	if err != nil {
		log.Printf("can't add order to cache, err: %e", err)
		return false
	}

	return true
}

func (s *Server) deleteFromCache(data model.Order) bool {
	err := s.cache.DeleteOrder(data.OrderUid)
	if err != nil {
		log.Printf("can't delete order from cache, err: %e", err)
		return false
	}
	return true
}
