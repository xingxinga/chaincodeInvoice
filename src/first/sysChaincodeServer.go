package main

import (
	"strings"
	"fmt"
)
/**
获取账本列表中最新的账本信息
 */
func getNewestLedger(historyLedger []LedgerHistory) (LedgerHistory) {
	var newestLedger LedgerHistory
	var max int64;
	max = 0
	for _ , v := range historyLedger {//range returns both the index and value
		seconds:= v.Times.Seconds
		fmt.Printf("Ledger TxID data: %s\n",v.TxID)
		fmt.Printf("Ledger Value data: %s\n",v.Value)
		fmt.Printf("Ledger Seconds data: %d\n",seconds)
		if(max<seconds){
			max = seconds
			newestLedger = v
		}
	}
	return newestLedger
}

/*
获取交易中的交易数据
*/

func getTransactionValue(transaction string) (string) {
	left :=  strings.Index(transaction, "\"invoiceCode\"")
	right :=  strings.Index(transaction, "\"}")
	invoiceString := string(transaction[left -5:right+2+5])
	left =  strings.Index(invoiceString, "\"invoiceCode\"")
	right =  strings.Index(invoiceString, "\"}")
	invoiceString = string(invoiceString[left:right+2])
	invoiceString = "{"+invoiceString
	return invoiceString
}
