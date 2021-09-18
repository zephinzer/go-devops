package devops

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// SshKeyFingerprint holds the key's fingerprint
type SshKeyFingerprint interface {
	// GetSHA256 returns a SHA256 fingerprint as a string
	// that looks like 'SHA256:aBcDeF ...'
	GetSHA256() string

	// GetMD5 returns a MD5 fingerprint (legacy) as a string
	// that looks like 'aa:bb:cc:dd ...'
	GetMD5() string
}

type sshKeyFingerprint struct {
	sha256 string
	md5    string
}

func (o sshKeyFingerprint) GetSHA256() string {
	return o.sha256
}

func (o sshKeyFingerprint) GetMD5() string {
	return o.md5
}

// GetSshKeyFingerprintOpts provides the configuration values
// for identifying the key whose fingerprint we want
type GetSshKeyFingerprintOpts struct {
	IsPrivateKey bool
	IsPublicKey  bool
	Passphrase   string
	Path         string
}

// GetSshKeyFingerprint returns the fingerprints for the
// provided key as specified in the GetSshKeyFingerprintOpts
// parameter
func GetSshKeyFingerprint(opts GetSshKeyFingerprintOpts) (SshKeyFingerprint, error) {
	keyPath := opts.Path
	keyContent, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file at '%s': %s", keyPath, err)
	}

	var publicKey ssh.PublicKey
	if opts.IsPrivateKey {
		privateKey, err := ssh.ParsePrivateKey(keyContent)
		if err != nil {
			if _, ok := err.(*ssh.PassphraseMissingError); ok {
				if opts.Passphrase == "" {
					return nil, fmt.Errorf("failed to provide a required passphrase")
				} else {
					privateKey, err = ssh.ParsePrivateKeyWithPassphrase(keyContent, []byte(opts.Passphrase))
					if err != nil {
						return nil, fmt.Errorf("failed to parse private key using provided passphrase: %s", err)
					}
				}
			} else {
				return nil, fmt.Errorf("failed to parse private key: %s", err)
			}
		}
		publicKey = privateKey.PublicKey()
	} else if opts.IsPublicKey {
		publicKey, _, _, _, err = ssh.ParseAuthorizedKey(keyContent)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %s", err)
		}
	}

	fingerprint := sshKeyFingerprint{
		md5:    ssh.FingerprintLegacyMD5(publicKey),
		sha256: ssh.FingerprintSHA256(publicKey),
	}
	return fingerprint, nil
}
