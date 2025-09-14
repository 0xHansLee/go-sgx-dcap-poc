package enclave

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	intelApiURL = "https://api.trustedservices.intel.com/sgx/certification/v3/tcb"
	//intelApiKey = "your_intel_api_key"
)

func FetchTcbInfo(fmspc string) (string, []byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?fmspc=%s", intelApiURL, fmspc), nil)
	if err != nil {
		return "", nil, err
	}
	//req.Header.Set("Ocp-Apim-Subscription-Key", intelApiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("intel PCCS API failed: %s", resp.Status)
	}
	var data struct {
		TcbInfo   map[string]interface{} `json:"tcbInfo"`
		Signature string                 `json:"signature"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", nil, err
	}
	tcbInfoBytes, err := json.Marshal(data.TcbInfo)
	if err != nil {
		return "", nil, err
	}
	sigBytes, err := base64.StdEncoding.DecodeString(data.Signature)
	if err != nil {
		return "", nil, err
	}

	return string(tcbInfoBytes), sigBytes, nil
}

func ExtractFmspc(quote []byte) string {
	if len(quote) < 546 {
		return ""
	}
	return hex.EncodeToString(quote[540:546])
}
