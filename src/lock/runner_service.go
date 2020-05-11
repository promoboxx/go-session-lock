package lock

import (
	"sync"

	"github.com/divideandconquer/go-consul-client/src/config"
	"github.com/promoboxx/go-discovery/src/discovery"
)

// RunnerServer is the interface to allow operations on a RunnerServer
type RunnerServer interface {
	Run() error
	Stop() *sync.WaitGroup
}

type runner struct {
	runner *Runner
	state  string
}

func (r *runner) run() error {
	err := r.runner.Run()
	if err != nil {
		return err
	}
	r.state = "running"
	return nil
}

func (r *runner) stop() *sync.WaitGroup {
	if r.state == "running" {
		return r.runner.Stop()
	}
	return nil
}

type runnerServer struct {
	environment string
	serviceName string
	conf        config.Loader
	finder      discovery.Finder
	tracer      Tracer
	runners     []*runner
	// client manager field here?
}

// NewRunnerServer returns a RunnerServer
func NewRunnerServer(env, serviceName string, conf config.Loader, finder discovery.Finder, tracer Tracer) RunnerServer {
	ret := &runnerServer{environment: env, serviceName: serviceName, conf: conf, finder: finder, tracer: tracer}
	ret.init()
	return ret
}

func (s *runnerServer) Run() error {
	for _, runner := range s.runners {
		err := runner.run()
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
		wg := runner.stop()
		go func() {
			wg.Wait()
			ret.Done()
		}()
	}
	return &ret
}
