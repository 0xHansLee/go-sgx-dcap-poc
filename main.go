package main

import (
	"encoding/hex"
	"fmt"
	"github.com/edgelesssys/ego/enclave"
	"log"
)

func main() {
	// 1. Load or generate QE key pair
	//pub, _, err := enclave.LoadOrCreateSealedQEKey()
	//if err != nil {
	//	log.Fatalf("failed to generate keypair: %v", err)
	//}
	//fmt.Println("Public key: ", hexutil.Encode(pub))

	// 2. Create a new Eth Client
	//ethClient, err := eth.NewEthClient()
	//if err != nil {
	//	log.Fatalf("failed to create eth client: %v", err)
	//}

	// 3. Get real quote from enclave
	//fmt.Println("pub key as a report data", hexutil.Encode(pub))
	//report, err := egoenclave.GetRemoteReport(pub)
	//if err != nil {
	//	log.Fatalf("failed to get quote: %v", err)
	//}

	// 4. Extract FMSPC from quote
	//fmspc := enclave.ExtractFmspc(report)
	//if fmspc == "" {
	//	log.Fatal("failed to extract fmspc from quote")
	//}
	//fmt.Println("Extracted fmspc:", fmspc)

	// 5. Fetch TCB info from Intel
	//tcbInfoStr, sigBytes, err := enclave.FetchTcbInfo(fmspc)
	//if err != nil {
	//	log.Fatal("failed to fetch TCBInfo:", err)
	//}

	// 6. Upsert the TCBInfo.json on-chain
	//txFmspcTcb, err := ethClient.SubmitFmspcTcb(tcbInfoStr, sigBytes)
	//if err != nil {
	//	log.Fatal("failed to submit FMSPC TCBInfo:", err)
	//}
	//fmt.Println("upsert FMSPC TCB tx submitted:", txFmspcTcb.Hash().Hex())

	// 7. Generate QEIdentity
	//identityStr, _, err := enclave.GenerateQEIdentity(pub, report)
	//if err != nil {
	//	log.Fatalf("failed to generate QE identity: %v", err)
	//}

	// 8. Sign the JSON
	//signature := ed25519.Sign(priv, []byte(identityStr))

	// 9. Submit a transaction for registration of QE identity
	//tx, err := ethClient.SubmitQEIdentity(identityStr, signature)
	//if err != nil {
	//	log.Fatalf("contract call failed: %v", err)
	//}
	//fmt.Println("upsert QE identity tx submitted:", tx.Hash().Hex())

	//fmspc, err := extractFMSPC(report)
	//if err != nil {
	//	fmt.Println("err in extracting fmspc: ", err)
	//} else {
	//	fmt.Println("fmspc: ", fmspc)
	//}

	// 10. Submit Quote for On-Chain Remote Attestation
	//parsedReport, err := egoenclave.VerifyRemoteReport(report)
	//if err != nil {
	//	fmt.Println("error: ", err)
	//} else {
	//	fmt.Println("generated report")
	//	fmt.Println("data (pub key): ", hexutil.Encode(parsedReport.Data))
	//	fmt.Println("security version: ", parsedReport.SecurityVersion)
	//	fmt.Println("unique id: ", hexutil.Encode(parsedReport.UniqueID))
	//	fmt.Println("signer id: ", hexutil.Encode(parsedReport.SignerID))
	//	fmt.Println("product id: ", hexutil.Encode(parsedReport.ProductID))
	//	fmt.Println("tcb status: ", parsedReport.TCBStatus)
	//	fmt.Println("tcb advisories: ", parsedReport.TCBAdvisories)
	//	fmt.Println("tcb advisories err: ", parsedReport.TCBAdvisoriesErr)
	//}

	//txQuote, err := ethClient.VerifyAndAttestOnChain(report)
	//if err != nil {
	//	log.Fatalf("failed to submit quote: %v", err)
	//}
	//fmt.Println("quote verification tx submitted:", txQuote.Hash().Hex())

	expected := make([]byte, 32)
	for i := 0; i < 32; i++ {
		expected[i] = byte(i)
	}

	log.Printf("Expected report data = %s \n", hex.EncodeToString(expected))

	var reportData [64]byte
	copy(reportData[:32], expected)

	log.Printf("Local reportData[0:32] = %s \n", hex.EncodeToString(reportData[:32]))
	log.Printf("Local reportData[32:64] = %s \n", hex.EncodeToString(reportData[32:]))

	quote, err := enclave.GetRemoteReport(reportData[:])
	if err != nil {
		log.Fatalf("failed to generate remote report: %v", err)
	}

	log.Printf("Quote length = %d", len(quote))
	if len(quote) >= 432 {
		reportDataFromQuote := quote[368:432]
		log.Printf("Extracted report_data = %s \n", hex.EncodeToString(reportDataFromQuote))
	}
}

func extractFMSPC(quote []byte) (string, error) {
	const offset = 464
	if len(quote) < offset+6 {
		return "", fmt.Errorf("quote too short")
	}
	fmspc := quote[offset : offset+6]

	// reverse to big endian
	for i, j := 0, len(fmspc)-1; i < j; i, j = i+1, j-1 {
		fmspc[i], fmspc[j] = fmspc[j], fmspc[i]
	}

	return hex.EncodeToString(fmspc), nil
}
