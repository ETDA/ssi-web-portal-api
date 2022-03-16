package consts

type KeyType string

const (
	KeyTypeECDSA KeyType = "ECDSA"
	KeyTypeRSA   KeyType = "RSA"

	KeyTypeSecp256r12019 = "EcdsaSecp256r1VerificationKey2019"
	KeyTypeRSA2018       = "RsaVerificationKey2018"
)
