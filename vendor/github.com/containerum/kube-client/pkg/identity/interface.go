package identity

type Changer interface {
	SetToken(string)
	SetFingerprint(string)
}
