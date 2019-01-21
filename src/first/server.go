package main

import (
	"bytes"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

/*
根据查询语句查询状态数据库信息
*/
func getListResultBuffer(stub shim.ChaincodeStubInterface,sql string) (bytes.Buffer,error){
	var buffer bytes.Buffer
	resultsIterator,error:= stub.GetQueryResult(sql)
	if error !=nil{
		return buffer,error
	}
	defer resultsIterator.Close()
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return buffer, err
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
	//return buffer.Bytes(), nil
	return buffer, nil
}
/*
根据查询语句查询历史账本信息
*/
func getHistoryResultBuffer(stub shim.ChaincodeStubInterface,key string) (bytes.Buffer,error){
	var buffer bytes.Buffer
	resultsIterator,error:= stub.GetHistoryForKey(key)
	if error !=nil{
		return buffer,error
	}
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return buffer, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		item,_:= json.Marshal( queryResponse)
		buffer.Write(item)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("queryResult:\n%s\n", buffer.String())
	//return buffer.Bytes(), nil
	return buffer, nil
}

/**
 根据map中条件，使用and逻辑查询列表
 */
func getListForAndMap(stub shim.ChaincodeStubInterface,mapValue map[string]string) ([]byte,error){
	var sql = buildAndSql(mapValue)
	value , err:= getListResultBuffer(stub,sql)
	return value.Bytes() , err
}
/**
 根据map中条件，使用or逻辑查询列表
 */
func getListForOrMap(stub shim.ChaincodeStubInterface,mapValue map[string]string) ([]byte,error){
	var sql = buildOrSql(mapValue)
	value , err:= getListResultBuffer(stub,sql)
	return value.Bytes() , err
}
/*
根据查询语句查询状态数据库信息并附带key值
*/
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

/**
 调用系统链码查询
 */
func invokeChainCode(stub shim.ChaincodeStubInterface,channelName, chaincodeName, fuc string , args []string ) (string){
	fmt.Printf("invokeChainCode channelName:%s\n", channelName)
	fmt.Printf("invokeChainCode chaincodeName:%s\n", chaincodeName)
	fmt.Printf("invokeChainCode fuc:%s\n", fuc)
	fmt.Printf("invokeChainCode args:%s\n", args)
	trans:=[][]byte{[]byte(fuc)}
	for _, v := range args {//range returns both the index and value
		trans = append(trans,[]byte(v))
	}
	response:= stub.InvokeChaincode(chaincodeName,trans,channelName)
	fmt.Printf("invokeChainCode Message:%s\n", response.Message)
	return string(response.Payload)
}

