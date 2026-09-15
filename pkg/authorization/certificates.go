package authorization

import (
	"crypto/rsa"
	"os"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

// variabes para las firmas
var (
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
	once      sync.Once
)

// cargar una sola vez el archivo
func LoadFiles(privateFile, publicFile string) error {
	var err error
	once.Do(func() {
		err = loadFiles(privateFile, publicFile)
	})
	return err
}

// leer el contenido de los archivos de las claves en sus variables
func loadFiles(privateFile, publicFile string) error {
	privateBytes, err := os.ReadFile(privateFile)
	if err != nil {
		return err
	}

	publicBytes, err := os.ReadFile(publicFile)
	if err != nil {
		return err
	}

	return parseRSA(privateBytes, publicBytes)
}

// analizar los bytes y cargar las variables signKey y verifyKey
func parseRSA(privateBytes, publicBytes []byte) error {
	var err error
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		return err
	}

	return nil
}
