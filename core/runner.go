package core

type Runner interface {
	ChStdOut() chan []byte
	ChStdErr() chan []byte
	Run(Task) error
	Close()
}

func GetRunner(host ConfigServer) (Runner, error) {
	var runner Runner
	var err error
	if host.RunLocally {
		runner, err = NewLocalRunner()
		if err != nil {
			return nil, err
		}
	} else {
		runner, err = NewRemoteRunner(host)
		if err != nil {
			return nil, err
		}
	}

	return runner, nil
}
