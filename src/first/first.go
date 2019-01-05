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

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	firstServer := &FirstServer{}
	function, args := stub.GetFunctionAndParameters()

	switch function {

	  case "create":return firstServer.createInvoice(stub, args) //已完成

	  case "getInvoiceInfo":return firstServer.getInvoice(stub, args) //已完成

	  case "updateInvoiceAttribution":return firstServer.updateInvoiceAttribution(stub, args) //已完成

	  case "updateInvoiceFinancingBank":return firstServer.updateInvoiceFinancingBank(stub, args) //已完成

	  case "getUserInvoiceList":return firstServer.getUserInvoiceList(stub)//已完成

	  case "getBankInvoiceList":return firstServer.getBankInvoiceList(stub)//已完成

	  case "select":return t.selectInvoice(stub, args)

	  case "test":return t.test(stub, args)


	}

	return shim.Error("canshu leixing cuowu")
}


func (t *SimpleChaincode) selectInvoice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("canshu shuliang cuowu")
	}

	results,error := getListResult(stub,args[0])
	if error !=nil{
		return shim.Error("zhuanhuan cuowu")
	}
	return shim.Success(results)
}

func (t *SimpleChaincode) test(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orgName,_:= getOrgName(stub)
	userName,_:=getUserName(stub)
	fmt.Printf("orgName is: %s", orgName)
	fmt.Printf("userName is: %s", userName)
	fmt.Printf("financingBankOrg is: %s", financingBankOrg)
	testCertificateIssuer(stub,args)
	testCertificateSubject(stub,args)
	return shim.Success([]byte("aaaaaaaa"))
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