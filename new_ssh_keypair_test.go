package devops

import (
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/ssh"
)

type NewSSHKeypairTest struct {
	suite.Suite
}

func TestNewSSHKeypair(t *testing.T) {
	suite.Run(t, &NewSSHKeypairTest{})
}

func (s NewSSHKeypairTest) Test_NewSSHKeypairOpts_SetDefaults() {
	opts := NewSSHKeypairOpts{}
	opts.SetDefaults()
	s.Equal(DefaultSSHKeyLength, opts.Bytes)
}

func (s NewSSHKeypairTest) Test_NewSSHKeypairOpts_Validate() {
	opts := NewSSHKeypairOpts{}
	err := opts.Validate()
	s.NotNil(err)
	s.Contains(err.Error(), "missing key length")
}

func (s NewSSHKeypairTest) Test_NewSSHKeypair() {
	keypair, err := NewSSHKeypair(NewSSHKeypairOpts{Bytes: 1024})
	s.Nil(err)

	privatePEM, _ := pem.Decode(keypair.Private)
	privateKey, err := x509.ParsePKCS1PrivateKey(privatePEM.Bytes)
	s.Nil(err)
	publicKey, _, _, _, err := ssh.ParseAuthorizedKey(keypair.Public)
	s.Nil(err)

	s.EqualValues(privateKey.Public(), publicKey)
}
