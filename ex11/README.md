
编译
go build -o go-execsql.exe


配置文件说明

;配置文件与可执行程序同名、同文件夹
;例如: go-execsql.exe 的配置文件为go-exesql.ini

[Mysql]
;数据源名称
DSN="wanglei:123123@tcp(192.168.1.192:3306)/btc_dealing?charset=utf8"

[Action]
;开启事务模式 允许的值 0, 1 , yes, no, on, off
; 0, no, off 关闭
; 1, yes, on 开启
;useTrans=on

;执行SQL，可以多行
;多行SQL按顺序依次执行
SQL=update address_log set state='1' where address='0xc45A06040061E006bFFCD7edfc124613a9F79825'
SQL=update address_log set receive=receive+1, update_number=update_number+1 where address='0xc45A06040061E006bFFCD7edfc124613a9F79825'
SQL=update address_log set last_block=987654 where address='0xc45A06040061E006bFFCD7edfc124613a9F79825'
SQL=insert address_log set address='0x11111111111111111111111111111111'
SQL=delete from address_log where address='0x11111111111111111111111111111111'

;错误出现时
;如果事务模式，将产生回滚，所有SQL执行结果被取消
;非事务模式，错误被忽略继续执行
;SQL=delete from address_log2 where address='0x11111111111111111111111111111111' ;table address_log2 doesn't exist


wanglei.ok@foxmail.com
