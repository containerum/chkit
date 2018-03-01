package cherry

import (
	"bytes"
	"fmt"
	"net/http"
)

// ErrSID -- represents service ID of error
//go:generate stringer -type=ErrSID
type ErrSID uint64

// ErrKind -- represents kind of error
type ErrKind uint64

// ErrID -- represnsents unique error ID
type ErrID struct {
	SID  ErrSID  `json:"sid"`
	Kind ErrKind `json:"kind"`
}

func (errID *ErrID) String() string {
	return fmt.Sprintf("%v-%v", errID.SID, errID.Kind)
}

// Err -- standart serializable API error
// Message -- constant error message:
//		+ "invalid username"
//		+ "quota exceeded"
//		+ "validation error"
//		...etc...
// ID -- unique error identification code
// Details -- optional context error messages kinda
// 		+ "field 'Replicas' must be non-zero value"
//		+ "not enough tights to feed gopher"
//		+ "resource 'God' does't exist"
type Err struct {
	Message    string   `json:"message"`
	StatusHTTP int      `json:"status_http"`
	ID         ErrID    `json:"id"`
	Details    []string `json:"details,omitempty"`
}

// NewErr -- constructs Err struct with provided message and ID
func NewErr(msg string, status int, ID ErrID) *Err {
	return &Err{
		Message:    msg,
		StatusHTTP: status,
		ID:         ID,
	}
}

// BuildErr -- produces Err constructor with custom
// ID prefix
// Example:
// 	MyErr := BuildErr("serivice_id")
//  ErrNotEnoughCheese = MyErr("not enough cheese", "666")
//  	--> "not enough cheese [service_id666]"
func BuildErr(SID ErrSID) func(string, int, ErrKind) *Err {
	return func(msg string, status int, kind ErrKind) *Err {
		return NewErr(msg, status, ErrID{SID: SID, Kind: kind})
	}
}

// Returns text representation kinda
// "unable to parse quota []"
func (err *Err) Error() string {
	buf := bytes.NewBufferString(" [" + err.ID.String() + "] " +
		http.StatusText(err.StatusHTTP) + " " +
		err.Message)
	detailsLen := len(err.Details)
	if detailsLen > 0 {
		buf.WriteString(": ")
	}
	for i, msg := range err.Details {
		if i+1 == detailsLen {
			buf.WriteString(msg)
		} else {
			buf.WriteString(msg + "; ")
		}
	}
	return buf.String()
}

// AddDetails -- adds detail messages to Err, chainable
func (err *Err) AddDetails(details ...string) *Err {
	err.Details = append(err.Details, details...)
	return err
}

// AddDetailF --adds formatted message to Err, chainable
func (err *Err) AddDetailF(formatS string, args ...interface{}) *Err {
	return err.AddDetails(fmt.Sprintf(formatS, args...))
}

// AddDetailsErr -- adds errors as detail messages to Err, chainable
func (err *Err) AddDetailsErr(details ...error) *Err {
	for _, detail := range details {
		err.AddDetails(detail.Error())
	}
	return err
}
