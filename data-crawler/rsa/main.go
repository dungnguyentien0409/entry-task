package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func main() {
	// Generate the keys
	priv, pub := generateRsaKeyPair(512)
	// Format the keys into PEM
	privPem, pubPem, _ := exportKeytoPEM(priv, pub)
	// Store the key to file
	writeFile([]byte(privPem), "private.pem")
	writeFile([]byte(pubPem), "public.pem")

	writeFilePrivate([]byte(privPem), "jwtRS256.key")
	writeFilePublic([]byte(pubPem), "jwtRS256.key.pub")
}

func writeFile(content []byte, filename string) (err error) {
	absPath, _ := filepath.Abs("")
	filepath := fmt.Sprintf(absPath+"/data-crawler/rsa/"+"%s", filename)
	err = ioutil.WriteFile(filepath, content, 0644)
	if err != nil {
		return
	}
	return
}

func writeFilePrivate(content []byte, filename string) (err error) {
	absPath, _ := filepath.Abs("")
	filepath := fmt.Sprintf(absPath+"/tcp-server/"+"%s", filename)
	err = ioutil.WriteFile(filepath, content, 0644)
	if err != nil {
		return
	}
	return
}

func writeFilePublic(content []byte, filename string) (err error) {
	absPath, _ := filepath.Abs("")
	filepath := fmt.Sprintf(absPath+"/http-server/"+"%s", filename)
	err = ioutil.WriteFile(filepath, content, 0644)
	if err != nil {
		return
	}

	filepath = fmt.Sprintf(absPath+"/tcp-server/"+"%s", filename)
	err = ioutil.WriteFile(filepath, content, 0644)
	if err != nil {
		return
	}
	return
}

func exportKeytoPEM(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) (
	privPEM, pubPEM string, err error) {

	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privPEMBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privKeyBytes,
		},
	) // step 3
	privPEM = string(privPEMBytes)
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return
	}
	pubPEMBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubKeyBytes,
		},
	) // step 4
	pubPEM = string(pubPEMBytes)
	return
}

func generateRsaKeyPair(size int) (privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) {
	privKey, _ = rsa.GenerateKey(rand.Reader, size) // step 1
	pubKey = &privKey.PublicKey                     // step 2
	return
}
