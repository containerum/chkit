package cherry

// ErrLogger -- interface for logging origin and returned errors due to origin errors discarding
type ErrorLogger interface {
	Log(origin error, returning *Err)
}

// Log -- logs origin error for returning error using ErrLogger, Chainable.
func (err *Err) Log(origin error, logger ErrorLogger) *Err {
	logger.Log(origin, err)
	return err
}
