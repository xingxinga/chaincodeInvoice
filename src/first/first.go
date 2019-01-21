package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	_"time"
	"fmt"
	"bytes"
	"encoding/pem"
	"crypto/x509"
	"strconv"
	"encoding/json"
)

type Invoice struct {
	InvoiceCode          string     `json:"invoiceCode"`
	InvoiceNo            string   	`json:"invoiceNo"`
	InvoiceCreatedate    string 	`json:"invoiceCreatedate"`
	InvoiceAmount        float64 	`json:"invoiceAmount"`    	     //总额
	InvoiceTaxtotal	     float64	`json:"invoiceTaxtotal"` 	     //税额
	InvoiceTotal	     float64	`json:"invoiceTotal"`    	     //含税金额
	InvoiceAttribution   string     `json:"invoiceAttribution"`     //发票归属方
	InvoiceBuyer   	     string     `json:"invoiceBuyer"`    	    //发票买方
	InvoiceSeller        string     `json:"invoiceSeller"`          //发票卖方
	InvoiceFinancingBank string 	`json:"invoiceFinancingBank"`  //融资银行
	CreateBy             string     `json:"createBy"`                //创建用户
}
type SimpleChaincode struct {
}
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface,) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 0 {
		return shim.Error("canshu shuliang cuowu")
	}
	return shim.Success([]byte("success"))
}

/**
  chaincode 执行交易逻辑列表
 */
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	firstServer := &FirstServer{}
	function, args := stub.GetFunctionAndParameters()

	switch function {

	  case "create":return firstServer.createInvoice(stub, args) //发票上传

	  case "getInvoiceInfo":return firstServer.getInvoice(stub, args) //获取发票信息

	  case "updateInvoiceAttribution":return firstServer.updateInvoiceAttribution(stub, args) //修改发票归属方

	  case "updateInvoiceFinancingBank":return firstServer.updateInvoiceFinancingBank(stub, args) //修改融资银行

	  case "getUserInvoiceList":return firstServer.getUserInvoiceList(stub)//获取用户发票集合

	  case "getBankInvoiceList":return firstServer.getBankInvoiceList(stub)//获取银行发票集合

	  case "getRelationInvoiceList":return firstServer.getRelationInvoiceList(stub)//获取买卖方的发票集合

	  case "getInvoiceHistory":return firstServer.getInvoiceHistory(stub,args)//获取发票历史

	  case "select":return t.selectInvoice(stub, args)

	  case "test":return firstServer.getNewInvoiceTest(stub, args)

	  case "testList":return t.testList(stub, args)


	}

	return shim.Error("canshu leixing cuowu")
}


func (t *SimpleChaincode) selectInvoice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("canshu shuliang cuowu")
	}

	results,error := getListResultBuffer(stub,args[0])
	if error !=nil{
		return shim.Error("zhuanhuan cuowu")
	}
	return shim.Success(results.Bytes())
}

func (t *SimpleChaincode) test(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	stub.GetHistoryForKey("")
	orgName,_:= getOrgName(stub)
	userName,_:=getUserName(stub)
	fmt.Printf("orgName is: %s", orgName)
	fmt.Printf("userName is: %s", userName)
	fmt.Printf("financingBankOrg is: %s", financingBankOrg)
	testCertificateIssuer(stub,args)
	testCertificateSubject(stub,args)
	return shim.Success([]byte("aaaaaaaa"))
}

func (t *SimpleChaincode) testList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	userName,err := getUserName(stub)
	mapValue := make(map[string]string)
	mapValue["createBy"] = userName
	value,err := getListForAndMap(stub,mapValue)
	if(err!=nil){
		return shim.Error(err.Error())
	}
	var list []Invoice
	json.Unmarshal(value, &list)
	fmt.Println("********** begin list *********")
	for _, v := range list {//range returns both the index and value
		fmt.Println("InvoiceCode : "+v.InvoiceCode)
		fmt.Println("InvoiceNo : "+v.InvoiceNo)
	}
	return shim.Success(value)
}

func testCertificateIssuer(stub shim.ChaincodeStubInterface, args []string)(string,string){
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

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Errorf("ParseCertificate failed")
	}
	var orgname string
	orgname = ""

	pkixName := cert.Issuer
	uname:=pkixName.CommonName
	org:=pkixName.Organization
	OrganizationalUnit:=pkixName.OrganizationalUnit
	Country:=pkixName.Country
	Locality:=pkixName.Locality
	SerialNumber:=pkixName.SerialNumber
	StreetAddress:=pkixName.StreetAddress
	fmt.Println("Subject.Organization.index:"+strconv.Itoa(len(org)))
	fmt.Println("Subject.OrganizationalUnit.index:"+strconv.Itoa(len(OrganizationalUnit)))
	fmt.Println("Subject.Country.index:"+strconv.Itoa(len(Country)))
	fmt.Println("Subject.Locality.index:"+strconv.Itoa(len(Locality)))
	fmt.Println("Subject.SerialNumber:"+SerialNumber)
	fmt.Println("Subject.StreetAddress.index:"+strconv.Itoa(len(StreetAddress)))
	for _, v := range org { // ignores index
		fmt.Println("****:"+v)
		orgname = orgname + v
	}
	for _, v := range OrganizationalUnit { // ignores index
		fmt.Println("OrganizationalUnit:"+v)
	}
	for _, v := range Country { // ignores index
		fmt.Println("Country:"+v)
	}
	for _, v := range Locality { // ignores index
		fmt.Println("Locality:"+v)
	}
	for _, v := range StreetAddress { // ignores index
		fmt.Println("StreetAddress:"+v)
	}
	fmt.Println("Subject.CommonName:"+uname)
	fmt.Println("Subject.Organization:"+orgname)
	return uname,orgname
}

func testCertificateSubject(stub shim.ChaincodeStubInterface, args []string)(string,string){
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

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Errorf("ParseCertificate failed")
	}
	var orgname string
	orgname = ""

	pkixName := cert.Issuer
	uname:=pkixName.CommonName
	org:=pkixName.Organization
	OrganizationalUnit:=pkixName.OrganizationalUnit
	Country:=pkixName.Country
	Locality:=pkixName.Locality
	SerialNumber:=pkixName.SerialNumber
	StreetAddress:=pkixName.StreetAddress
	fmt.Println("Subject.Organization.index:"+strconv.Itoa(len(org)))
	fmt.Println("Subject.OrganizationalUnit.index:"+strconv.Itoa(len(OrganizationalUnit)))
	fmt.Println("Subject.Country.index:"+strconv.Itoa(len(Country)))
	fmt.Println("Subject.Locality.index:"+strconv.Itoa(len(Locality)))
	fmt.Println("Subject.SerialNumber:"+SerialNumber)
	fmt.Println("Subject.StreetAddress.index:"+strconv.Itoa(len(StreetAddress)))
	for _, v := range org { // ignores index
		fmt.Println("****:"+v)
		orgname = orgname + v
	}
	for _, v := range OrganizationalUnit { // ignores index
		fmt.Println("OrganizationalUnit:"+v)
	}
	for _, v := range Country { // ignores index
		fmt.Println("Country:"+v)
	}
	for _, v := range Locality { // ignores index
		fmt.Println("Locality:"+v)
	}
	for _, v := range StreetAddress { // ignores index
		fmt.Println("StreetAddress:"+v)
	}
	fmt.Println("Subject.CommonName:"+uname)
	fmt.Println("Subject.Organization:"+orgname)
	return uname,orgname
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
