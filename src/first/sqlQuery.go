package main

var appendSql = "\"_id\": { \"$gt\": null }"

var begin string = "{ \"selector\" : {"

var end string = "} } \" "

func buildSql(sql string) (string){
	var resultSql string
	resultSql = begin + sql + end
	return resultSql
}

func buildEqualsSql(mapValue map[string]string) (string){
	var resultSql string
	var mqpSql string = ""
	for k, v := range mapValue {//range returns both the index and value
		mqpSql = mqpSql + "\"" + k + "\":" + "\"" + v + "\"" + ","
	}
	mqpSql = mqpSql+appendSql
	resultSql = begin + mqpSql + end
	return resultSql
}
