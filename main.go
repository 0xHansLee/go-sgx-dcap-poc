package main

import (
	"fmt"
	egoenclave "github.com/edgelesssys/ego/enclave"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hanslee/go-sgx-dcap-poc/enclave"
	"github.com/hanslee/go-sgx-dcap-poc/eth"
	"log"
)

func main() {
	// 1. Load or generate QE key pair
	pub, _, err := enclave.LoadOrCreateSealedQEKey()
	if err != nil {
		log.Fatalf("failed to generate keypair: %v", err)
	}
	fmt.Println("Public key: ", hexutil.Encode(pub))

	// 2. Create a new Eth Client
	ethClient, err := eth.NewEthClient()
	if err != nil {
		log.Fatalf("failed to create eth client: %v", err)
	}

	// 3. Get real quote from enclave
	report, err := egoenclave.GetRemoteReport(pub)
	if err != nil {
		log.Fatalf("failed to get quote: %v", err)
	}

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

	// 10. Submit Quote for On-Chain Remote Attestation
	parsedReport, _ := egoenclave.VerifyRemoteReport(report)
	fmt.Println("report: ", parsedReport)

	txQuote, err := ethClient.VerifyAndAttestOnChain(report)
	if err != nil {
		log.Fatalf("failed to submit quote: %v", err)
	}
	fmt.Println("quote verification tx submitted:", txQuote.Hash().Hex())
}
