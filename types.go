package devops

type BasicAuth struct {
	Username string
	Password string
}

// SshKeyFingerprint holds the key's fingerprint
type SshKeyFingerprint interface {
	// GetSHA256 returns a SHA256 fingerprint as a string
	// that looks like 'SHA256:aBcDeF ...'
	GetSHA256() string

	// GetMD5 returns a MD5 fingerprint (legacy) as a string
	// that looks like 'aa:bb:cc:dd ...'
	GetMD5() string
}
