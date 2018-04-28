package main


import (
	"fmt"
	"lib1"
)

func main() {

	fmt.Println("Test Lib1 load......")
	lib1.Printf("Test Lib1 load......")

	mysql_config := lib1.Database_str_global("east_database")

	sql_str_value := make(map[string]string)

	//填充数据库字段和字段的值
	sql_str_value["assets_name"] = "assets 1"
	sql_str_value["assets_desc"] = "This is a test"

	//插入数据
	//生成 MySQL 插入语句
	sql_str := lib1.Assemble_insert(sql_str_value, "user_assets_type")
	lib1.Mysql_connect_query_insert_api_3(sql_str, mysql_config)

}
