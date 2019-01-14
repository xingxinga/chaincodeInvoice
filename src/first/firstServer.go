package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"strconv"
	"fmt"
)
type FirstServer struct {
}

/**
  获取发票信息
 */
func (F *FirstServer)getInvoice(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) != 2 {
		return shim.Error("canshu shuliang cuowu")
	}
	value,err := getInvoiceByte(stub,args[0],args[1])
	if(err!=nil){
		return shim.Error(err.Error())
	}
	return shim.Success(value)
}

/**
  获取发票信息
 */
func (F *FirstServer)getBankInvoiceList(stub shim.ChaincodeStubInterface) pb.Response{
	orgName,err:= getOrgName(stub)
	if(err!=nil){
		return shim.Error(err.Error())
	}
	if(orgName!=financingBankOrg){
		return shim.Error("Insufficient authority")
	}
	mapValue := make(map[string]string)
	mapValue["invoiceFinancingBank"] = orgName
	value,err := getListForAndMap(stub,mapValue)
	if(err!=nil){
		return shim.Error(err.Error())
	}
	return shim.Success(value)
}

/**
  获取发票信息
 */
func (F *FirstServer)getRelationInvoiceList(stub shim.ChaincodeStubInterface) pb.Response{
	userName,err := getUserName(stub)
	mapValue := make(map[string]string)
	mapValue["invoiceSeller"] = userName
	mapValue["invoiceBuyer"] = userName
	value,err := getListForOrMap(stub,mapValue)
	if(err!=nil){
		return shim.Error(err.Error())
	}
	return shim.Success(value)
}
/**
  获取发票信息
 */
func (F *FirstServer)getUserInvoiceList(stub shim.ChaincodeStubInterface) pb.Response{
	userName,err := getUserName(stub)
	mapValue := make(map[string]string)
	mapValue["createBy"] = userName
	value,err := getListForAndMap(stub,mapValue)
	if(err!=nil){
		return shim.Error(err.Error())
	}
	return shim.Success(value)
}

/**
  创建发票信息
 */
func (F *FirstServer)createInvoice(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) != 8 {
		return shim.Error("canshu shuliang cuowu")
	}
	value,_ := getInvoice(stub,args[0],args[1])
	if(value.InvoiceCode!=""&&value.InvoiceNo!=""){
		return shim.Error("invoice already exist")
	}else{
		invoiceCode := args[0]
		invoiceNo := args[1]
		invoiceCreatedate := args[2]
		invoiceAmount,_ := strconv.ParseFloat(args[3], 64)
		invoiceTaxtotal,_  := strconv.ParseFloat(args[4], 64)
		invoiceTotal,_  := strconv.ParseFloat(args[5], 64)
		//invoiceAttribution := args[6]
		invoiceBuyer := args[6]
		invoiceSeller := args[7]
		invoiceAttribution,err := getUserName(stub)
		invoiceFinancingBank := ""
		createBy := invoiceAttribution
		if(err!=nil){
			return shim.Error(err.Error())
		}
		invoice := Invoice{invoiceCode,invoiceNo,invoiceCreatedate,invoiceAmount,invoiceTaxtotal,invoiceTotal,invoiceAttribution,invoiceBuyer,invoiceSeller,invoiceFinancingBank,createBy}
		value,err := saveInvoice(stub,invoice)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(value)
	}
}
/**
   修改发票信息
 */
func (F *FirstServer)updateInvoiceFinancingBank(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) != 3 {
		return shim.Error("canshu shuliang cuowu")
	}
	invoice,_ := getInvoice(stub,args[0],args[1])
	if(invoice.InvoiceCode==""||invoice.InvoiceNo==""){
		return shim.Error("invoice not exist")
	}
	userName,err := getUserName(stub)
	if(err!=nil){
		return shim.Error(err.Error())
	}
	if(invoice.InvoiceAttribution!=userName){
		return shim.Error("invoice Insufficient authority")
	}
	invoice.InvoiceFinancingBank = args[2]
	result,err := saveInvoice(stub,invoice)
	if(string(result)!="success"||err!=nil){
		return shim.Error("Error")
	}else{
		return shim.Success([]byte("InvoiceFinancingBank update success"))
	}
}

/**
   修改发票信息
 */
func (F *FirstServer)updateInvoiceAttribution(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	orgName,err:= getOrgName(stub)
	if(err!=nil){
		return shim.Error(err.Error())
	}
	if(orgName!=financingBankOrg){
		return shim.Error("Insufficient authority")
	}
	if len(args) != 2 {
		return shim.Error("canshu shuliang cuowu")
	}
	invoice,_ := getInvoice(stub,args[0],args[1])
	if(invoice.InvoiceFinancingBank!=orgName){
		return shim.Error("FinancingBank Inconformity")
	}
	invoiceAttribution,err := getUserName(stub)
	invoice.InvoiceAttribution = invoiceAttribution
	result,err := saveInvoice(stub,invoice)
	if(string(result)!="success"||err!=nil){
		return shim.Error("Error")
	}else{
		return shim.Success([]byte("invoiceAttribution update success"))
	}
}

func (F *FirstServer)selectInvoice(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) != 2 {
		return shim.Error("canshu shuliang cuowu")
	}
	value,_ := getInvoiceByte(stub,args[0],args[1])
	return shim.Success(value)
}


func getInvoiceByte(stub shim.ChaincodeStubInterface,invoiceCode string, invoiceNo string) ([]byte,error){
	var err error
	key,error:=creatKey(invoiceCode,invoiceNo)
	if error !=nil{
		return nil,error
	}
	value,error := stub.GetState(key)
	return value,err
}

func getInvoice(stub shim.ChaincodeStubInterface,invoiceCode string, invoiceNo string) (Invoice,error){
	var invoice Invoice
	value,error:=getInvoiceByte(stub,invoiceCode,invoiceNo)
	if error !=nil{
		return invoice,error
	}
	json.Unmarshal(value, &invoice)
	return invoice,error
}

func saveInvoice(stub shim.ChaincodeStubInterface,invoice Invoice) ([]byte,error){
	key,err:=creatKey(invoice.InvoiceCode,invoice.InvoiceNo)
	if err != nil {
		return nil,err
	}

	_json, err := json.Marshal(invoice)
	fmt.Printf("saveInvoice:\n%s\n", _json)
	err = stub.PutState(key, _json)
	if err != nil {
		return nil,err
	}else{
		return []byte("success"),err
	}
}


func creatKeyByStub(stub shim.ChaincodeStubInterface,invoiceCode string ,invoiceNo string) (string, error){
	return stub.CreateCompositeKey("InvoiceKey",[]string{invoiceCode,invoiceNo})
}

func creatKey(invoiceCode string ,invoiceNo string) (string, error){
	return invoiceCode+"_"+invoiceNo,nil
}