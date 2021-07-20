package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/Gunyoung-Kim/blockchain/utils"
)

const (
	walletFileName string = "coin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(walletFileName)
	return !os.IsNotExist(err)
}

func createPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	return privateKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleError(err)
	utils.HandleError(os.WriteFile(walletFileName, bytes, 0644)) // read and write
}

func restoreKey() (key *ecdsa.PrivateKey) {
	keyAsBytes, err := os.ReadFile(walletFileName)
	utils.HandleError(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleError(err)
	return
}

func encodeBigInts(a, b *big.Int) string {
	c := append(a.Bytes(), b.Bytes()...)
	return fmt.Sprintf("%x", c)
}

func AddressFromKey(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X, key.Y)
}

func Sign(payload string, w *wallet) string {
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleError(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsBytes)
	utils.HandleError(err)
	return encodeBigInts(r, s)
}

func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]
	bigFirst, bigSecond := big.Int{}, big.Int{}
	bigFirst.SetBytes(firstHalfBytes)
	bigSecond.SetBytes(secondHalfBytes)
	return &bigFirst, &bigSecond, nil
}

func verify(signature, payload, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleError(err)
	x, y, err := restoreBigInts(address)
	utils.HandleError(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleError(err)
	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			w.privateKey = restoreKey()
		} else {
			key := createPrivateKey()
			persistKey(key)
			w.privateKey = key
		}
		w.Address = AddressFromKey(w.privateKey)
	}

	return w
}
