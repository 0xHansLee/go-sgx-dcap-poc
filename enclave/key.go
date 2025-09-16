package enclave

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"

	"github.com/edgelesssys/ego/ecrypto"
)

const keyFile = "/home_mnt/priv_key.sealed"

type QEKey struct {
	Private ed25519.PrivateKey
}

func LoadOrCreateSealedQEKey() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	if _, err := os.Stat(keyFile); err == nil {
		data, err := os.ReadFile(keyFile)
		if err != nil {
			return nil, nil, fmt.Errorf("read sealed key file: %w", err)
		}

		unsealed, err := ecrypto.Unseal(data, nil)
		if err != nil {
			return nil, nil, fmt.Errorf("unseal key: %w", err)
		}

		var key QEKey
		if err := json.Unmarshal(unsealed, &key); err != nil {
			return nil, nil, fmt.Errorf("unmarshal key: %w", err)
		}

		return key.Private.Public().(ed25519.PublicKey), key.Private, nil
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("generate key: %w", err)
	}

	key := QEKey{Private: priv}
	raw, err := json.Marshal(key)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal key: %w", err)
	}

	sealed, err := ecrypto.SealWithUniqueKey(raw, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("seal key: %w", err)
	}

	if err := os.WriteFile(keyFile, sealed, 0600); err != nil {
		return nil, nil, fmt.Errorf("write sealed key: %w", err)
	}

	return pub, priv, nil
}
