package devops

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

const (
	DefaultSSHKeyLength = 8192
)

type SSHKeypair struct {
	Private []byte
	Public  []byte
}

type NewSSHKeypairOpts struct {
	Bytes    int
	Password string
}

func (o *NewSSHKeypairOpts) SetDefaults() {
	if o.Bytes == 0 {
		o.Bytes = DefaultSSHKeyLength
	}
}

func (o NewSSHKeypairOpts) Validate() error {
	errors := []string{}

	if o.Bytes == 0 {
		errors = append(errors, "missing key length")
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to validate options: ['%s']", strings.Join(errors, "', '"))
	}
	return nil
}

func NewSSHKeypair(opts NewSSHKeypairOpts) (*SSHKeypair, error) {
	opts.SetDefaults()
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("failed to create new ssh key pair: %s", err)
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, opts.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate a private key: %s", err)
	}
	privateKeyPEM := &pem.Block{
		Type:  "RSA_PRIVATE_KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if opts.Password != "" {
		privateKeyPEM, err = x509.EncryptPEMBlock(rand.Reader, privateKeyPEM.Type, privateKeyPEM.Bytes, []byte(opts.Password), x509.PEMCipherAES256)
		if err != nil {
			return nil, fmt.Errorf("failed to protect private key with a password: %s", err)
		}
	}
	privateKeyData := pem.EncodeToMemory(privateKeyPEM)
	if privateKeyData == nil {
		return nil, fmt.Errorf("failed to encode private key: %s", err)
	}
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate the public key: %s", err)
	}
	publicKeyData := ssh.MarshalAuthorizedKey(publicKey)
	return &SSHKeypair{
		Private: privateKeyData,
		Public:  publicKeyData,
	}, nil
}
