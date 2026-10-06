package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	labRootCertFile = "certs/lab-root-ca.crt"
	labRootKeyFile  = "certs/lab-root-ca.key"
	labRootDerFile  = "certs/lab-root-ca.der"
	serverCertFile  = "certs/server.crt"
	serverKeyFile   = "certs/server.key"
)

func randomSerial() (*big.Int, error) {
	return rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
}

func writePEM(path, typ string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer f.Close()
	return pem.Encode(f, &pem.Block{Type: typ, Bytes: data})
}

func GenerateCertificate() error {
	if err := os.MkdirAll("certs", 0755); err != nil {
		return err
	}

	// Local CA used only by the emulator lab.
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	caSerial, err := randomSerial()
	if err != nil {
		return err
	}
	caTemplate := &x509.Certificate{
		SerialNumber: caSerial,
		Subject: pkix.Name{CommonName: "Nextendo MHRise Lab Root CA"},
		NotBefore: time.Now().Add(-time.Minute),
		NotAfter:  time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA: true,
		SubjectKeyId: []byte{0x4E, 0x58, 0x4C, 0x41, 0x42},
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return err
	}
	if err := writePEM(labRootCertFile, "CERTIFICATE", caDER, 0644); err != nil {
		return err
	}
	caKeyDER, err := x509.MarshalPKCS8PrivateKey(caKey)
	if err != nil {
		return err
	}
	if err := writePEM(labRootKeyFile, "PRIVATE KEY", caKeyDER, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(labRootDerFile, caDER, 0644); err != nil {
		return err
	}

	// Server certificate signed by the local CA.
	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	serverSerial, err := randomSerial()
	if err != nil {
		return err
	}
	serverTemplate := &x509.Certificate{
		SerialNumber: serverSerial,
		Subject: pkix.Name{CommonName: "t-e1c218b5-lp1.lp1.t.npln.srv.nintendo.net"},
		NotBefore: time.Now().Add(-time.Minute),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames: []string{
			"t-e1c218b5-lp1.lp1.t.npln.srv.nintendo.net",
			"localhost",
		},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
	}
	serverDER, err := x509.CreateCertificate(rand.Reader, serverTemplate, caTemplate, &serverKey.PublicKey, caKey)
	if err != nil {
		return err
	}
	if err := writePEM(serverCertFile, "CERTIFICATE", serverDER, 0644); err != nil {
		return err
	}
	serverKeyDER, err := x509.MarshalPKCS8PrivateKey(serverKey)
	if err != nil {
		return err
	}
	if err := writePEM(serverKeyFile, "PRIVATE KEY", serverKeyDER, 0600); err != nil {
		return err
	}

	fmt.Println("[+] Created lab-root-ca.crt/.key/.der and server.crt/.key")
	return nil
}
