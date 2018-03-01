package cherry

const (
	// UninitializedSID -- default SID
	UninitializedSID ErrSID = iota
	Auth                    // 1
	KubeAPI                 // 2
	ResourceService         // 3
	UserManager             // 4
	Billing                 // 5
	Gluster                 // 6
	MailTemplaiter          // 7
	Gateway                 // 8
	Archive                 // 9
	Cache                   // 10
)
