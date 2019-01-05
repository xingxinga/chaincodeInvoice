package main

import (
	"fmt"
	"bytes"
	"encoding/pem"
	"crypto/x509"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func getUserName(stub shim.ChaincodeStubInterface)(string,error){
	cert, err := getCertificate(stub)
	if err != nil {
		fmt.Errorf("ParseCertificate failed")
	}
	userName:=cert.Subject.CommonName
	return userName,err
}

func getOrgName(stub shim.ChaincodeStubInterface)(string,error){
	cert, err := getCertificate(stub)

	if err != nil {
		fmt.Errorf("ParseCertificate failed")
	}
	org:=cert.Issuer.Organization
	orgName:= org[0]
	//orgName:=cert.Subject.CommonName
	return orgName,err
}

func getCertificate(stub shim.ChaincodeStubInterface)(*x509.Certificate, error) {
	creatorByte,_:= stub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		fmt.Errorf("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		fmt.Errorf("Could not decode the PEM structure")
	}
	return x509.ParseCertificate(bl.Bytes)
}
