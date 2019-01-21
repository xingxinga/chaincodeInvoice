package main

var appendAndSql = "\"_id\": { \"$gt\": null }"

var appendOrSql = "{ \"_id\": { \"$eq\": null } }"

var begin string = "{ \"selector\" : {"

var end string = "} } \" "

/**
 根据Map构建and逻辑的查询语句
 */
func buildAndSql(mapValue map[string]string) (string){
	var mqpSql string = ""
	for k, v := range mapValue {//range returns both the index and value
		mqpSql = mqpSql + "\"" + k + "\":" + "\"" + v + "\"" + ","
	}
	mqpSql = mqpSql+appendAndSql
	return getrResultSql(mqpSql)
}

/**
 根据Map构建or逻辑的查询语句
 */
func buildOrSql(mapValue map[string]string) (string){
	var mqpSql string = " \"$or\": [ "
	for k, v := range mapValue {//range returns both the index and value
		mqpSql = mqpSql + " {\"" + k + "\":" + "\"" + v + "\" }" + ","
	}
	mqpSql = mqpSql+appendOrSql + "]"
	return getrResultSql(mqpSql)
}

/**
 拼接查询语句
 */
func getrResultSql(sql string) (string){
	var resultSql string
	resultSql = begin + sql + end
	return resultSql
}
