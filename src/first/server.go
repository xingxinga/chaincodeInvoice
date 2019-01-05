package main

import (
	"bytes"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)



func getListResultAppendKey(stub shim.ChaincodeStubInterface,sql string) ([]byte,error){
	resultsIterator,error:= stub.GetQueryResult(sql)
	if error !=nil{
		return nil,error
	}
	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}

func getListResult(stub shim.ChaincodeStubInterface,sql string) ([]byte,error){
	resultsIterator,error:= stub.GetQueryResult(sql)
	if error !=nil{
		return nil,error
	}
	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}

func getListForMap(stub shim.ChaincodeStubInterface,mapValue map[string]string) ([]byte,error){
	var sql = buildEqualsSql(mapValue)
	value , err:= getListResult(stub,sql)
	return value , err
}