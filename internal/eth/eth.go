package eth

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func NewClient(ctx context.Context, rpcUrl string) (*ethclient.Client, error) {
	return ethclient.DialContext(ctx, rpcUrl)
}

func ParsePrivateKey(rawPrivateKey string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(rawPrivateKey)
}
