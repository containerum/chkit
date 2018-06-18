package model

// Resources -- represents namespace limits and user resources.
//
// swagger:model
type Resources struct {
	// Hard resource limits
	//
	// required: true
	Hard Resource  `json:"hard"`
	Used *Resource `json:"used,omitempty"`
}

// Resource -- represents namespace CPU and RAM
//
// swagger:model
type Resource struct {
	// CPU in m
	//
	// required: true
	CPU uint `json:"cpu"`
	// RAM in Mi
	//
	// required: true
	Memory uint `json:"memory"`
}

// UpdateNamespaceName -- contains new namespace label
//
// swagger:model
type UpdateNamespaceName struct {
	// required: true
	Label string `json:"label"`
}

// NamespacesList -- model for namespaces list
//
// swagger:model
type NamespacesList struct {
	Namespaces []Namespace `json:"namespaces"`
}

// Namespace -- namespace representation
//
// swagger:model
type Namespace struct {
	ID string `json:"id,omitempty"`
	//creation date in RFC3339 format
	CreatedAt  *string `json:"created_at,omitempty"`
	Owner      string  `json:"owner,omitempty"`
	OwnerLogin string  `json:"owner_login,omitempty"`
	// user-visible label for the namespace
	Label         string      `json:"label,omitempty"`
	Access        AccessLevel `json:"access,omitempty"`
	TariffID      string      `json:"tariff_id",omitempty`
	MaxExtService uint        `json:"max_ext_service,omitempty"`
	MaxIntService uint        `json:"max_int_service,omitempty"`
	MaxTraffic    uint        `json:"max_traffic,omitempty"`
	// required: true
	Resources Resources    `json:"resources,omitempty"`
	Users     []UserAccess `json:"users,omitempty"`
}

// Mask removes information not interesting for users
func (ns *Namespace) Mask() {
	ns.Owner = ""
}
