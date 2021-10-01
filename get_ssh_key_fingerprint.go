package devops

import (
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/ssh"
)

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
	// IsPrivateKey if set to true indicates that we are
	// targetting a private key. If both this and .IsPublicKey
	// are not set, .IsPublicKey will be set to true
	IsPrivateKey bool

	// IsPublicKey if set to true indicates that we are
	// targetting a public key. If both this and .IsPrivateKey
	// are not set, .IsPublicKey will be set to true
	IsPublicKey bool

	// Passphrase defines a passphrase for the private key if
	// applicable
	Passphrase string

	// Path defines the file directory path to the key file
	// of interest
	Path string
}

// SetDefaults sets defaults for this object instance
func (o *GetSshKeyFingerprintOpts) SetDefaults() {
	if !o.IsPrivateKey && !o.IsPublicKey {
		o.IsPublicKey = true
	}
}

// Validate checks if this object instance has sufficient parameters
// to be used by GetSshKeyFingerprint
func (o GetSshKeyFingerprintOpts) Validate() error {
	errors := []string{}

	if o.Path == "" {
		errors = append(errors, "missing path to key")
	}

	if o.IsPrivateKey && o.IsPublicKey {
		errors = append(errors, "cannot be both public and private key")
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to validate options: ['%s']", strings.Join(errors, "', '"))
	}
	return nil
}

// GetSshKeyFingerprint returns the fingerprints for the
// provided key as specified in the GetSshKeyFingerprintOpts
// parameter
func GetSshKeyFingerprint(opts GetSshKeyFingerprintOpts) (SshKeyFingerprint, error) {
	opts.SetDefaults()
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("failed to get ssh fingerprint: %s", err)
	}
	keyPath := opts.Path
	/* #nosec - this is needed to read the file */
	keyContent, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file at '%s': %s", keyPath, err)
	}

	var publicKey ssh.PublicKey
	if opts.IsPrivateKey {
		privateKey, err := ssh.ParsePrivateKey(keyContent)
		if _, ok := err.(*ssh.PassphraseMissingError); ok {
			if opts.Passphrase == "" {
				return nil, fmt.Errorf("failed to provide a required passphrase")
			}
			privateKey, err = ssh.ParsePrivateKeyWithPassphrase(keyContent, []byte(opts.Passphrase))
			if err != nil {
				return nil, fmt.Errorf("failed to parse private key using provided passphrase: %s", err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %s", err)
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
