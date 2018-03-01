```go
type Err struct {
    Message    string   `json:"message"`
    StatusHTTP int      `json:"-"`
    ID         string   `json:"id"`
    Details    []string `json:"details,omitempty"`
}
```
    Err -- standart serializable API error Message -- constant error
    message:

	+ "invalid username"
	+ "quota exceeded"
	+ "validation error"
	...etc...

    ID -- unique error identification code Details -- optional context error
    messages kinda

	+ "field 'Replicas' must be non-zero value"
	+ "not enough tights to feed gopher"
	+ "resource 'God' does't exist"

```go
func BuildErr(prefix string) func(string, int, string) *Err
```
BuildErr -- produces Err constructor with custom ID prefix Example:


```go
	MyErr := BuildErr("42")
	ErrNotEnoughCheese = MyErr("not enough cheese", 404, "666")
```
     	--> "[42666] HTTP 400 not enough cheese "

```go
func NewErr(msg string, status int, ID string) *Err
``` 
NewErr -- constructs Err struct with provided message and ID

```go
func (err *Err) AddDetailF(formatS string, args ...interface{}) *Err
```    
AddDetailF --adds formatted message to Err, chainable

```go
func (err *Err) AddDetails(details ...string) *Err
```
AddDetails -- adds detail messages to Err, chainable

```go
func (err *Err) AddDetailsErr(details ...error) *Err
```
AddDetailsErr -- adds errors as detail messages to Err,chainable

```go
func (err *Err) Error() string
```
Returns text representation kinda "unable to parse quota []"

```go
func (err *Err) Gonic(ctx *gin.Context)
```
Gonic -- aborts gin HTTP request with StatusHTTP and provides json representation of error


