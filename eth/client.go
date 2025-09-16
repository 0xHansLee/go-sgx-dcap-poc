package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
	"strings"
)

const (
	pccsAbiPath                = "/home_mnt/abi/automata_pccs.abi.json"
	fmspcAbiPath               = "/home_mnt/abi/automata_fmspc.abi.json"
	attestationAbiPath         = "/home_mnt/abi/automata_attestation.abi.json"
	pccsContractAddress        = "0xe20C4d54afBbea5123728d5b7dAcD9CB3c65C39a"
	fmspcContractAddress       = "0x63eF330eAaadA189861144FCbc9176dae41A5BAf"
	attestationContractAddress = "0x95175096a9B74165BE0ac84260cc14Fc1c0EF5FF"
	rpcURL                     = "https://0xrpc.io/hoodi"
	chainID                    = 560048
)

type EthClient struct {
	Client              *ethclient.Client
	Auth                *bind.TransactOpts
	PCCSABI             abi.ABI
	FMSPCABI            abi.ABI
	AttestationABI      abi.ABI
	PCCSContract        *bind.BoundContract
	FMSPCContract       *bind.BoundContract
	AttestationContract *bind.BoundContract
}

func NewEthClient() (*EthClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		return nil, fmt.Errorf("PRIVATE_KEY environment variable not set")
	}

	privKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(chainID))
	if err != nil {
		return nil, err
	}

	pccsAbiData, err := os.ReadFile(pccsAbiPath)
	if err != nil {
		return nil, err
	}

	parsedPCCSABI, err := abi.JSON(strings.NewReader(string(pccsAbiData)))
	if err != nil {
		return nil, err
	}

	pccsBound := bind.NewBoundContract(common.HexToAddress(pccsContractAddress), parsedPCCSABI, client, client, client)

	fmspcAbiData, err := os.ReadFile(fmspcAbiPath)
	if err != nil {
		return nil, err
	}

	parsedFMSPCABI, err := abi.JSON(strings.NewReader(string(fmspcAbiData)))
	if err != nil {
		return nil, err
	}

	fmspcBound := bind.NewBoundContract(common.HexToAddress(fmspcContractAddress), parsedFMSPCABI, client, client, client)

	attestationAbiData, err := os.ReadFile(attestationAbiPath)
	if err != nil {
		return nil, err
	}

	parsedAttestationABI, err := abi.JSON(strings.NewReader(string(attestationAbiData)))
	if err != nil {
		return nil, err
	}

	attestationBound := bind.NewBoundContract(common.HexToAddress(attestationContractAddress), parsedAttestationABI, client, client, client)

	return &EthClient{
		Client:              client,
		Auth:                auth,
		PCCSABI:             parsedPCCSABI,
		FMSPCABI:            parsedFMSPCABI,
		AttestationABI:      parsedAttestationABI,
		PCCSContract:        pccsBound,
		FMSPCContract:       fmspcBound,
		AttestationContract: attestationBound,
	}, nil
}

func (ec *EthClient) SubmitQEIdentity(identityStr string, signature []byte) (*types.Transaction, error) {
	return ec.PCCSContract.Transact(ec.Auth, "upsertEnclaveIdentity", big.NewInt(0), big.NewInt(1), map[string]interface{}{
		"identityStr": identityStr,
		"signature":   signature,
	})
}

func (ec *EthClient) SubmitFmspcTcb(tcbInfoStr string, signature []byte) (*types.Transaction, error) {
	return ec.FMSPCContract.Transact(ec.Auth, "upsertFmspcTcb", map[string]interface{}{
		"tcbInfoStr": tcbInfoStr,
		"signature":  signature,
	})
}

func (ec *EthClient) VerifyAndAttestOnChain(rawQuote []byte) (*types.Transaction, error) {
	ec.Auth.Value = new(big.Int).SetUint64(10000000000000000)

	return ec.AttestationContract.Transact(ec.Auth, "verifyAndAttestOnChain", rawQuote)
}
