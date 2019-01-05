package main

import "github.com/hyperledger/fabric/core/chaincode/shim"

type ChaincodeStubClient interface {
	getChaincodeStubInterface() shim.ChaincodeStubInterface
}
