package enclave

import "math/big"

// Quote layout offsets based on Intel DCAP v3 format
const (
	OffsetMiscSelect = 368
	OffsetAttributes = 372
)

// EnclaveId is an enum represented as string in JSON (e.g., "QE", "QVE", "TD_QE")
type IdentityObj struct {
	ID                      string `json:"id"` // Must match one of: "QE", "QVE", "TD_QE"
	Version                 uint32 `json:"version"`
	IssueDateTimestamp      uint64 `json:"issueDateTimestamp"`
	NextUpdateTimestamp     uint64 `json:"nextUpdateTimestamp"`
	TCBEvaluationDataNumber uint32 `json:"tcbEvaluationDataNumber"`
	MiscSelect              string `json:"miscselect"`     // 4-byte hex string
	MiscSelectMask          string `json:"miscselectMask"` // 4-byte hex string
	Attributes              string `json:"attributes"`     // 16-byte hex string
	AttributesMask          string `json:"attributesMask"` // 16-byte hex string
	MRSIGNER                string `json:"mrsigner"`       // 32-byte hex string
	ISVPRODID               uint16 `json:"isvprodid"`
	TCB                     []any  `json:"tcb"` // Optional for PoC
}

type QEIdentityInput struct {
	ID          *big.Int
	Version     *big.Int
	IdentityStr string
	Signature   []byte
}
