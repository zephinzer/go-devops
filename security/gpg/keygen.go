package gpg

import (
	"bytes"
	"fmt"
)

type Keygen interface {
	Bytes() []byte
	String() string
}

type keygen struct {
	KeyType         string
	KeyLength       uint
	SubkeyType      string
	SubkeyLength    uint
	Name            string
	Comment         string
	Email           string
	ExpireDateCount string
	ExpireDateUnit  TimeUnit
}

func (k *keygen) Bytes() []byte {
	var data bytes.Buffer
	data.WriteString("%echo generating an OpenPGP key...\n")
	data.WriteString(fmt.Sprintf("Key-Type: %s\n", k.KeyType))
	data.WriteString(fmt.Sprintf("Key-Length: %v\n", k.KeyLength))
	data.WriteString(fmt.Sprintf("Subkey-Type: %s\n", k.SubkeyType))
	data.WriteString(fmt.Sprintf("Subkey-Length: %v\n", k.SubkeyLength))
	data.WriteString(fmt.Sprintf("Name-Real: %s\n", k.Name))
	data.WriteString(fmt.Sprintf("Name-Comment: %s\n", k.Comment))
	data.WriteString(fmt.Sprintf("Name-Email: %s\n", k.Email))
	data.WriteString(fmt.Sprintf("Expire-Date: %v%s\n", k.ExpireDateCount, k.ExpireDateUnit))
	data.WriteString("%commit")
	data.WriteString("%echo done")
	return data.Bytes()
}

func (k *keygen) String() string {
	return string(k.Bytes())
}
