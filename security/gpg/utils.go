package gpg

// CreateKey creates a new GPG key given the manifest for a key
func CreateKey(keyGenerator string) {

}

type CreateKeysConfigurer func(defaultKey *keygen) Keygen

func CreateKeygen(configure CreateKeysConfigurer) string {
	keygenInstance := configure(newKey())
	return keygenInstance.String()
}

func newKey() *keygen {
	return &keygen{
		KeyType: DefaultKeyType,
	}
}
