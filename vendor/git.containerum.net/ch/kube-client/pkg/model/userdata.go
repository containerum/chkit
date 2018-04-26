package model

// represents header data for X-User-Namespace and X-User-Volume headers (encoded in base64)
//
//swagger:model
type UserHeaderData struct {
	// hosting-internal name
	// required: true
	ID string `json:"id"`
	// user-visible label for the object
	// required: true
	Label string `json:"label"`
	// one of: "owner", "read", "write", "read-delete", "none"
	// required: true
	Access string `json:"access"`
}

// User --
//swagger:ignore
type User struct {
	Login     string   `json:"login"`
	Data      UserData `json:"data"`
	ID        string   `json:"id"`
	IsActive  bool     `json:"is_active"`
	CreatedAt string   `json:"created_at"`
}

// UserData --
//swagger:ignore
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

// Tokens --
//swagger:ignore
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// CheckTokenResponse --
//swagger:ignore
type CheckTokenResponse struct {
	Access struct {
		Namespace []Resource `json:"namespace"`
		Volume    []Resource `json:"volume"`
	} `json:"access"`
}

// Login --
//swagger:ignore
type Login struct {
	Login     string  `json:"login"`
	Password  string  `json:"password"`
	Recaptcha *string `json:"recaptcha,omitempty"`
}
