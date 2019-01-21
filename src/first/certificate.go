package main

import (
	"fmt"
	"bytes"
	"encoding/pem"
	"crypto/x509"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//获取当前操作用户名称
func getUserName(stub shim.ChaincodeStubInterface)(string,error){
	//获取当前操作用户证书信息
	cert, err := getCertificate(stub)
	if err != nil {
		fmt.Errorf("ParseCertificate failed")
	}
	//获取证书信息中的用户名
	userName:=cert.Subject.CommonName
	return userName,err
}

//获取操作用户所属组织名称
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

//获取当前操作用户证书信息
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
