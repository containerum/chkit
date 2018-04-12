package activekit

type Action interface {
	Run() (bool, error)
}

type ActionSimple func()

func (action ActionSimple) Run() (bool, error) {
	action()
	return true, nil
}

type ActionWithErr func() error

func (action ActionWithErr) Run() (bool, error) {
	return true, action()
}

type ActionFull func() (bool, error)

func (action ActionFull) Run() (bool, error) {
	ok, err := action()
	return ok, err
}
