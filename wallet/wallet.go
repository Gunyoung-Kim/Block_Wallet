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

// hasWalletFile return whether there is walletFile or not
func hasWalletFile() bool {
	_, err := os.Stat(walletFileName)
	return !os.IsNotExist(err)
}

// createPrivateKey create new {@code ecdsa.PrivateKey} then return it
func createPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	return privateKey
}

// persistKey persist privateKey into file
// turn {@code ecdsa.PrivateKey} into slice of bytes
func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleError(err)
	utils.HandleError(os.WriteFile(walletFileName, bytes, 0644)) // read and write
}

// restoreKey restore privateKey from walletFile
// turn slice of bytes into {@code *ecdsa.PrivateKey}
func restoreKey() (key *ecdsa.PrivateKey) {
	keyAsBytes, err := os.ReadFile(walletFileName)
	utils.HandleError(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleError(err)
	return
}

// encodeBigInts return hexa-decimal string made from two big int
func encodeBigInts(a, b *big.Int) string {
	c := append(a.Bytes(), b.Bytes()...)
	return fmt.Sprintf("%x", c)
}

//AddressFromKey return address(hex) made by public key(from input private key)
func AddressFromKey(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X, key.Y)
}

//Sign return hexa-decimal string made by r, s
// r, s made by input privatekey(from wallet) and payload(transaction id)
func Sign(payload string, w *wallet) string {
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleError(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsBytes)
	utils.HandleError(err)
	return encodeBigInts(r, s)
}

// restoreBigInts turn hexa-decimal string into two {@code big.int}
// this method is used for turning signature into r, s and turning address into x,y(used for making public key)
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

//Verify input signature is correct with transaction id and public key.
//public key is made by input address
func Verify(signature, payload, address string) bool {
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

// Wallet return pointer of wallet which is made by singleton pattern
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
