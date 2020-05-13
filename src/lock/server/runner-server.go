package server

import (
	"github.com/promoboxx/go-session-lock/src/lock"
	"sync"

	"github.com/divideandconquer/go-consul-client/src/config"
	"github.com/promoboxx/go-discovery/src/discovery"
)

// RunnerServer is the interface to allow operations on a RunnerServer
type RunnerServer interface {
	Run() error
	Stop() *sync.WaitGroup
}

type Runner interface {
	Run() error
	Stop() *sync.WaitGroup
}

type runnerServer struct {
	environment string
	serviceName string
	conf        config.Loader
	finder      discovery.Finder
	tracer      lock.Tracer
	runners     []Runner
}

// NewRunnerServer returns a RunnerServer
func NewRunnerServer(env, serviceName string, conf config.Loader, finder discovery.Finder, tracer lock.Tracer, runners []Runner) RunnerServer {
	ret := &runnerServer{environment: env, serviceName: serviceName, conf: conf, finder: finder, tracer: tracer, runners: runners}
	return ret
}

func (s *runnerServer) Run() error {
	for _, runner := range s.runners {
		err := runner.Run()
		if err != nil {
			s.Stop()
			return err
		}
	}
	return nil
}

func (s *runnerServer) Stop() *sync.WaitGroup {
	var ret sync.WaitGroup
	for _, runner := range s.runners {
		ret.Add(1)
		wg := runner.Stop()
		go func() {
			wg.Wait()
			ret.Done()
		}()
	}
	return &ret
}
