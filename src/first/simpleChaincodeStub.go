package main

import "github.com/hyperledger/fabric/core/chaincode/shim"

type SimpleChaincodeStub struct {
	stub shim.ChaincodeStubInterface
}

func (simpleChaincodeStub *SimpleChaincodeStub)getChaincodeStubInterface() (shim.ChaincodeStubInterface){
	//return stub.CreateCompositeKey("InvoiceKey",[]string{invoiceCode,invoiceNo})
	return simpleChaincodeStub.stub
}

func (simpleChaincodeStub *SimpleChaincodeStub)setChaincodeStubInterface(stub shim.ChaincodeStubInterface) {
	simpleChaincodeStub.stub = stub
}