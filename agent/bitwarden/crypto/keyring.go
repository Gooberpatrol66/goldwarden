package crypto

import (
	"errors"
)

type Keyring struct {
	AccountKey               SymmetricEncryptionKey
	AsymmetricEncyryptionKey AsymmetricEncryptionKey
	IsMemguard               bool
	OrganizationKeys         map[string]string
}

func NewMemoryKeyring(accountKey *MemorySymmetricEncryptionKey) Keyring {
	return Keyring{
		AccountKey: accountKey,
	}
}

func NewMemguardKeyring(accountKey *MemguardSymmetricEncryptionKey) Keyring {
	return Keyring{
		AccountKey: accountKey,
	}
}

func (keyring Keyring) IsLocked() bool {
	return keyring.AccountKey == nil
}

func (keyring *Keyring) Lock() {
	keyring.AccountKey = nil
	keyring.AsymmetricEncyryptionKey = MemoryAsymmetricEncryptionKey{}
	keyring.OrganizationKeys = nil
}

func (keyring *Keyring) GetSymmetricKeyForOrganization(uuid string) (SymmetricEncryptionKey, error) {
	if key, ok := keyring.OrganizationKeys[uuid]; ok {
		decryptedOrgKey, err := DecryptWithAsymmetric([]byte(key), keyring.AsymmetricEncyryptionKey)
		if err != nil {
			return MemorySymmetricEncryptionKey{}, err
		}

		return MemorySymmetricEncryptionKeyFromBytes(decryptedOrgKey)
	}
	return MemorySymmetricEncryptionKey{}, errors.New("no key found for organization")
}
