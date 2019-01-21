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

//账本历史数据
type LedgerHistory struct {
	TxID           string     `json:"tx_id"`
	Value          string     `json:"value"`
	Times          Timestamp  `json:"timestamp"`
}
//账本历史数据中时间信息
type Timestamp struct {
	Seconds         int64     `json:"seconds"`
	Nanos           int64     `json:"nanos"`
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
  获取发票历史交易的基本信息
 */
func (F *FirstServer)getInvoiceHistory(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) != 2 {
		return shim.Error("canshu shuliang cuowu")
	}
	invoiceCode:= args[0]
	invoiceNo:= args[1]
	key,err:=creatKey(invoiceCode,invoiceNo)
	if err !=nil{
		return shim.Error(err.Error())
	}
	//获取发票历史
	value,error := getHistoryResultBuffer(stub,key)
	if(error!=nil){
		return shim.Error(error.Error())
	}
	return shim.Success(value.Bytes())
}

/**
  银行用户获取发票列表，条件是融资银行为当前用户所属银行组织的发票
 */
func (F *FirstServer)getBankInvoiceList(stub shim.ChaincodeStubInterface) pb.Response{
	//获取当前用户所属组织
	orgName,err:= getOrgName(stub)
	if(err!=nil){
		return shim.Error(err.Error())
	}
	//用户组织不为银行组织
	if(orgName!=financingBankOrg){
		return shim.Error("Insufficient authority")
	}
	mapValue := make(map[string]string)
	mapValue["invoiceFinancingBank"] = orgName
	//获取根据融资银行信息查询发票列表
	value,err := getListForAndMap(stub,mapValue)
	if(err!=nil){
		return shim.Error(err.Error())
	}
	return shim.Success(value)
}

/**
  获取买卖方相关的发票列表
 */
func (F *FirstServer)getRelationInvoiceList(stub shim.ChaincodeStubInterface) pb.Response{
	//获取操作用户名称
	userName,err := getUserName(stub)
	mapValue := make(map[string]string)
	mapValue["invoiceSeller"] = userName
	mapValue["invoiceBuyer"] = userName
	//查询买卖方为当前用户名称的发票列表
	value,err := getListForOrMap(stub,mapValue)
	if(err!=nil){
		return shim.Error(err.Error())
	}
	return shim.Success(value)
}
/**
  获取当前用户创建的发票
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
	//获取账本中当前发票数据
	value,_ := getNewestInvoice(stub,args[0],args[1])
	//当发票已存在返回错误
	if(value.InvoiceCode!=""&&value.InvoiceNo!=""){
		return shim.Error("invoice already exist")
	}else{
		//构造发票信息
		invoiceCode := args[0]
		invoiceNo := args[1]
		invoiceCreatedate := args[2]
		invoiceAmount,_ := strconv.ParseFloat(args[3], 64)
		invoiceTaxtotal,_  := strconv.ParseFloat(args[4], 64)
		invoiceTotal,_  := strconv.ParseFloat(args[5], 64)
		//invoiceAttribution := args[6]
		invoiceBuyer := args[6]
		invoiceSeller := args[7]
		//发票归属方为fabric用户
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
   修改发票融资银行信息
 */
func (F *FirstServer)updateInvoiceFinancingBank(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) != 3 {
		return shim.Error("canshu shuliang cuowu")
	}
	invoice,_ := getNewestInvoice(stub,args[0],args[1])
	if(invoice.InvoiceCode==""||invoice.InvoiceNo==""){
		return shim.Error("invoice not exist")
	}
	//获取fabric用户名称
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
   修改发票归属方信息
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
	invoice,_ := getNewestInvoice(stub,args[0],args[1])
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

/**
获取发票信息byte方式
 */
func getInvoiceByte(stub shim.ChaincodeStubInterface,invoiceCode string, invoiceNo string) ([]byte,error){
	var err error
	key,error:=creatKey(invoiceCode,invoiceNo)
	if error !=nil{
		return nil,error
	}
	value,error := stub.GetState(key)
	return value,err
}

/**
获取状态数据库的发票信息
 */
func getStateDBInvoice(stub shim.ChaincodeStubInterface,invoiceCode string, invoiceNo string) (Invoice,error){
	var invoice Invoice
	value,error:=getInvoiceByte(stub,invoiceCode,invoiceNo)
	if error !=nil{
		return invoice,error
	}
	json.Unmarshal(value, &invoice)
	return invoice,error
}


/**
  获取发票账本历史数据字符串
 */
func getInvoiceHistoryLedgerString(stub shim.ChaincodeStubInterface,invoiceCode string, invoiceNo string) (string,error){
	result := ""
	key,err:=creatKey(invoiceCode,invoiceNo)
	if err !=nil{
		return result,err
	}
	//查询发票历史
	value,error := getHistoryResultBuffer(stub,key)
	if error !=nil{
		return result,error
	}
	return value.String(),error
}

/**
获取当前发票最新的数据
 */
func getNewestInvoice(stub shim.ChaincodeStubInterface,invoiceCode string, invoiceNo string) (Invoice,error){
	var invoice Invoice
	var historyLedger[]LedgerHistory
	//获取发票历史账本信息
	historyLedgerString,err :=  getInvoiceHistoryLedgerString(stub,invoiceCode,invoiceNo)
	if (err !=nil || historyLedgerString == "[]") {
		return invoice,err
	}
	//转化历史账本信息
	json.Unmarshal([]byte(historyLedgerString), &historyLedger)
	//获取账本列表中最新的账本信息
	newestLedger := getNewestLedger(historyLedger);
	//获取当前通道ID
	channelName:= stub.GetChannelID()
	//初始化调用系统链码参数
	chaincodeArgs :=[]string{channelName,newestLedger.TxID}
	//调用系统链码查询最新发票的交易信息
	stringTransaction := invokeChainCode(stub,channelName,sysChaincodeQsccName,sysChaincodeGetTransactionByID,chaincodeArgs)
	//提取交易中发票具体数据
	stringInvoice:=getTransactionValue(stringTransaction)
	json.Unmarshal([]byte(stringInvoice), &invoice)
	return invoice,err
}

/*
保存覆盖发票信息
*/
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

/**
  获取发票信息
 */
func (F *FirstServer)checkInvoice(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	fmt.Printf("kaisih ceshi lalala :\n")
	var historyLedger []LedgerHistory
	historyLedgerString,err :=  getInvoiceHistoryLedgerString(stub,args[0],args[1])
	fmt.Printf("historyLedgerString:%s\n", historyLedgerString)
	if err !=nil{
		return shim.Error("canshu shuliang cuowuaaaaaaaa")
	}
	fmt.Printf("begin zhuanhua!!!!!!!!!\n")
	json.Unmarshal([]byte(historyLedgerString), &historyLedger)
	fmt.Printf("jieguo changdu :%d\n", len(historyLedger))
	fmt.Printf("zhuanhua chenggong !!!!!!!!!\n")
	newestLedger := getNewestLedger(historyLedger);
	fmt.Printf("newestLedger TxID data: %s\n",newestLedger.TxID)
	channelName:= stub.GetChannelID()
	fmt.Printf("channelName: %s\n",channelName)
	chaincodeArgs :=[]string{channelName,newestLedger.TxID}
	fmt.Printf("begin diaoyong syschaincode!!!!!!!!!\n")
	stringTransaction := invokeChainCode(stub,channelName,sysChaincodeQsccName,sysChaincodeGetTransactionByID,chaincodeArgs)
	fmt.Printf("syschaincode jieguo :%s\n",stringTransaction)
	fmt.Printf("begin syschaincode jieguo Invoice!!!!!!!!!\n")
	stringInvoice:=getTransactionValue(stringTransaction)
	return shim.Success([]byte(stringInvoice))
}

func creatKeyByStub(stub shim.ChaincodeStubInterface,invoiceCode string ,invoiceNo string) (string, error){
	return stub.CreateCompositeKey("InvoiceKey",[]string{invoiceCode,invoiceNo})
}

func creatKey(invoiceCode string ,invoiceNo string) (string, error){
	return invoiceCode+"_"+invoiceNo,nil
}