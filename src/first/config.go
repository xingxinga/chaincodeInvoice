package main

/*
 chaincode中的静态变量
*/

//银行组织名称
var financingBankOrg string = "org3.example.com"

//系统负责查询的链码名称
var sysChaincodeQsccName string = "qscc"

//系统查询链码中根据交易ID查询交易详细信息的方法名称
var sysChaincodeGetTransactionByID string = "GetTransactionByID"

