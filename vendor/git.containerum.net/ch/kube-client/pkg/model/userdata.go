package model

type UserHeaderData struct {
	ID     string `json:"id"`     // hosting-internal name
	Label  string `json:"label"`  // user-visible label for the object
	Access string `json:"access"` // one of: "owner", "read", "write", "read-delete", "none"
}

type User struct {
	Login     string   `json:"login"`
	Data      UserData `json:"data"`
	ID        string   `json:"id"`
	IsActive  bool     `json:"is_active"`
	CreatedAt string   `json:"created_at"`
}

type UserData struct {
	Email          string `json:"email"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	IsOrganization bool   `json:"is_organization"`
	TaxCode        string `json:"tax_code"`
	Company        string `json:"company"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CheckTokenResponse struct {
	Access struct {
		Namespace []Resource `json:"namespace"`
		Volume    []Resource `json:"volume"`
	} `json:"access"`
}

type Login struct {
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	Recaptcha *string `json:"recaptcha,omitempty"`
}
