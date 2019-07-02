package onewallet

import (
	"os"
	"os/signal"
	"syscall"

	rabbit "github.com/djansyle/go-rabbit"
)

// Service holds the information of the service
type Service struct {
	name    string
	version uint16

	server rabbit.Server
}

// NewServiceOpt serves as the options available for instantiating service
type NewServiceOpt struct {
	URL     string
	Name    string
	Version uint16
	Root    interface{}
	Command interface{}
	Query   interface{}
}

// Name exposes the name of the service
func (s *Service) Name() string {
	return s.name
}

// Version exposos the version of the service
func (s *Service) Version() uint16 {
	return s.version
}

// Start the service and start accepting requests
func (s *Service) Start() {
	go func() {
		graceful := make(chan os.Signal)

		signal.Notify(graceful, syscall.SIGTERM)
		signal.Notify(graceful, syscall.SIGINT)

		<-graceful

		s.server.Close()
		os.Exit(0)
	}()

	s.server.Serve()
}

// NewService creates a new instance of the service
func NewService(opt NewServiceOpt) (*Service, error) {
	server, err := rabbit.CreateServer(opt.URL, opt.Name)
	if err != nil {
		return nil, err
	}

	if opt.Root != nil {
		server.RegisterName("", opt.Root)
	}

	if opt.Command != nil {
		server.RegisterName("Command", opt.Command)
	}

	if opt.Query != nil {
		server.RegisterName("Query", opt.Query)
	}

	return &Service{name: opt.Name, version: opt.Version, server: server}, nil
}
