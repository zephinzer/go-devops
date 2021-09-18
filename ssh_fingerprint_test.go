package devops

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestSshKeysFingerprint defines the md5/sha256 fingerprint hashes
// for the files in ./tests/sshkeys
var TestSshKeysFingerprint = map[string]sshKeyFingerprint{
	"id_rsa_1024": {
		md5:    "b8:c4:84:10:fb:33:62:93:9a:2a:5f:48:93:81:f0:bb",
		sha256: "SHA256:umIxkNk9Qt3BxCYqT8ZrUnMo3bEDpudnPIs0DgPLbqc",
	},
	"id_rsa_1024-w-password": {
		md5:    "32:74:7a:86:e3:54:8b:18:42:b7:88:91:bb:01:07:f8",
		sha256: "SHA256:gFJqayhy+zqU5axeguZXTgUY0DOCcxpd4legvOSd+tk",
	},
}

type SshFingerprintTests struct {
	suite.Suite
}

func TestSshFingerprint(t *testing.T) {
	suite.Run(t, &SshFingerprintTests{})
}

func (s SshFingerprintTests) Test_GetSshKeyFingerprint_privateKey() {
	keyPath := "./tests/sshkeys/id_rsa_1024"
	fingerprint, err := GetSshKeyFingerprint(GetSshKeyFingerprintOpts{
		IsPrivateKey: true,
		Path:         keyPath,
	})
	s.Nil(err)
	s.Equal(TestSshKeysFingerprint["id_rsa_1024"].md5, fingerprint.GetMD5(),
		"fingerprint should match the md5 generated from 'ssh-keygen -E md5 -lf ...'")
	s.Equal(TestSshKeysFingerprint["id_rsa_1024"].sha256, fingerprint.GetSHA256(),
		"fingerprint should match the sha256 generated from 'ssh-keygen -lf ...'")
}

func (s SshFingerprintTests) Test_GetSshKeyFingerprint_privateKeyWithPassphrase() {
	keyPath := "./tests/sshkeys/id_rsa_1024-w-password"
	options := GetSshKeyFingerprintOpts{
		IsPrivateKey: true,
		Path:         keyPath,
	}

	fingerprint, err := GetSshKeyFingerprint(options)
	s.NotNil(err)
	s.Contains(err.Error(), "failed to provide a required passphrase")

	options.Passphrase = "password"
	fingerprint, err = GetSshKeyFingerprint(options)
	s.Nil(err)
	s.Equal(TestSshKeysFingerprint["id_rsa_1024-w-password"].md5, fingerprint.GetMD5(),
		"fingerprint should match the md5 generated from 'ssh-keygen -E md5 -lf ...'")
	s.Equal(TestSshKeysFingerprint["id_rsa_1024-w-password"].sha256, fingerprint.GetSHA256(),
		"fingerprint should match the sha256 generated from 'ssh-keygen -lf ...'")
}

func (s SshFingerprintTests) Test_GetSshKeyFingerprint_publicKey() {
	keyPath := "./tests/sshkeys/id_rsa_1024.pub"
	fingerprint, err := GetSshKeyFingerprint(GetSshKeyFingerprintOpts{
		IsPublicKey: true,
		Path:        keyPath,
	})
	s.Nil(err)
	s.Equal(TestSshKeysFingerprint["id_rsa_1024"].md5, fingerprint.GetMD5(),
		"fingerprint should match the md5 generated from 'ssh-keygen -E md5 -lf ...'")
	s.Equal(TestSshKeysFingerprint["id_rsa_1024"].sha256, fingerprint.GetSHA256(),
		"fingerprint should match the sha256 generated from 'ssh-keygen -lf ...'")
}

func (s SshFingerprintTests) Test_GetSshKeyFingerprint_publicKeyWithPassphrasedPrivateKey() {
	keyPath := "./tests/sshkeys/id_rsa_1024-w-password.pub"
	fingerprint, err := GetSshKeyFingerprint(GetSshKeyFingerprintOpts{
		IsPublicKey: true,
		Path:        keyPath,
	})
	s.Nil(err)
	s.Equal(TestSshKeysFingerprint["id_rsa_1024-w-password"].md5, fingerprint.GetMD5(),
		"fingerprint should match the md5 generated from 'ssh-keygen -E md5 -lf ...'")
	s.Equal(TestSshKeysFingerprint["id_rsa_1024-w-password"].sha256, fingerprint.GetSHA256(),
		"fingerprint should match the sha256 generated from 'ssh-keygen -lf ...'")
}
