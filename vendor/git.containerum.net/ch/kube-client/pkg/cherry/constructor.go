package cherry

type ErrConstruct func(...func(*Err)) *Err

func (constr ErrConstruct) Error() string {
	return constr().Error()
}
