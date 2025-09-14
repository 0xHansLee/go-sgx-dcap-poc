package enclave

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

func GenerateQEIdentity(pub ed25519.PublicKey, quote []byte) (string, []byte, error) {
	if len(quote) < OffsetAttributes+16 {
		return "", nil, fmt.Errorf("quote too short: %d bytes", len(quote))
	}

	// Extract fields from quote
	mrsigner := sha256.Sum256(pub)
	miscselect := quote[OffsetMiscSelect : OffsetMiscSelect+4]
	attributes := quote[OffsetAttributes : OffsetAttributes+16]

	identity := IdentityObj{
		ID:                      "QE",
		Version:                 1,
		IssueDateTimestamp:      uint64(time.Now().Unix()),
		NextUpdateTimestamp:     uint64(time.Now().AddDate(1, 0, 0).Unix()),
		TCBEvaluationDataNumber: 1,
		MiscSelect:              hex.EncodeToString(miscselect),
		MiscSelectMask:          "ffffffff",
		Attributes:              hex.EncodeToString(attributes),
		AttributesMask:          "ff000000000000000000000000000000",
		MRSIGNER:                hex.EncodeToString(mrsigner[:]),
		ISVPRODID:               0,
		TCB:                     []any{},
	}

	jsonBytes, err := json.Marshal(identity)
	if err != nil {
		return "", nil, err
	}
	return string(jsonBytes), mrsigner[:], nil
}
