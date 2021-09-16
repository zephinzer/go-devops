package gpg

const (
	DefaultKeyType      = "RSA"
	DefaultKeyLength    = 4096
	DefaultSubkeyType   = "RSA"
	DefaultSubkeyLength = 4096
	DefaultName         = "undefined user"
	DefaultComment      = "(undefined comment)"
	DefaultEmail        = "undefined@email.com"
	DefaultExpireDate   = 365
	DefaultPassphrase   = "qweqwe123123!@#!@#"
)

var KeyTypes = []string{"RSA", "DSA"}

type TimeUnit string

const (
	Days   TimeUnit = "d"
	Weeks  TimeUnit = "w"
	Months TimeUnit = "m"
	Years  TimeUnit = "y"
)
