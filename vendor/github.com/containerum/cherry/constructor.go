package cherry

type ErrConstruct func(...func(*Err)) *Err

func (constr ErrConstruct) Error() string {
	return constr().Error()
}

func (constr ErrConstruct) AddDetails(details ...string) ErrConstruct {
	return func(options ...func(*Err)) *Err {
		err := constr().AddDetails(details...)
		for _, option := range options {
			option(err)
		}
		return err
	}
}

func (constr ErrConstruct) AddDetailsErr(details ...error) ErrConstruct {
	return func(options ...func(*Err)) *Err {
		err := constr().AddDetailsErr(details...)
		for _, option := range options {
			option(err)
		}
		return err
	}
}

func (constr ErrConstruct) AddDetailF(f string, vals ...interface{}) ErrConstruct {
	return func(options ...func(*Err)) *Err {
		err := constr().AddDetailF(f, vals...)
		for _, option := range options {
			option(err)
		}
		return err
	}
}

func(constr ErrConstruct) WithField(key, value string) ErrConstruct {
	return func(options ...func(*Err)) *Err {
		err := constr().WithField(key, value)
		for _, option := range options {
			option(err)
		}
		return err
	}
}



func(constr ErrConstruct) WithFields(fields Fields) ErrConstruct {
	return func(options ...func(*Err)) *Err {
		err := constr().WithFields(fields)
		for _, option := range options {
			option(err)
		}
		return err
	}
}

