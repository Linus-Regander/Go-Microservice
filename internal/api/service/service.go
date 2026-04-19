package service

import "log"

// Service holds information about an API service.
type Service struct {
	log *log.Logger
}

// New returns a new Service.
func New(logger *log.Logger) *Service {
	return &Service{
		log: logger,
	}
}
