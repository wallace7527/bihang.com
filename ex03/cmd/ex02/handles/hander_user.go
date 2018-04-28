/*
hander_user.go
wanglei.ok@foxmail.com

1.0
版本时间：2018年4月13日18:32:12


*/

package handles

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"bihang.com/ex03/router"
	"time"
	"log"
	"strconv"
	"bihang.com/ex03/cmd/ex02/modules"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"database/sql"
)


func init() {
	//router.Add("getuser","/user/:uid","GET", getuser)
	//router.Add("adduser","/user/","POST",	adduser)
	//router.Add("deleteuser","/user/:uid","DELETE",deleteuser)
	//router.Add("modifyuser","/user/:uid","PUT",modifyuser)
	//router.Add("getalluser","/user","GET",getalluser)
	//curl -i -H Accept:application/json -X GET http://localhost:8080/v1/user/getverifycode?tel=<phonenumber>
	router.Add("getVerifyCode","/v1/user/getverifycode","GET",getVerifyCode)
	//curl -i -H Accept:application/json -X GET http://localhost:8080/v1/user/register?verifycode=<vcode>&tel=<phonenumber>&un=<username>&pwd=<password>&nick=<nickname>&gender=<1 or 2>&id_number=<idnumber>
	router.Add("register","/v1/user/register","GET",register)
	//curl -i -H Accept:application/json -X GET http://localhost:8080/v1/user/login?un=<un or tel>&pwd=<password>
	router.Add("login","/v1/user/login","GET",login)

	//修改当前用户资料
	router.Add("changeinfo","/v1/user/changeinfo","GET",changeinfo)

	//修改当前用户密码
	//http://localhost:8080/v1/user/changepassword?token=<token>&id=<uid>&password=<pwd>&newpassword=<newpwd>
	router.Add("changepassword","/v1/user/changepassword","GET",changepassword)

	//找回密码
	router.Add("recoverpassword","/v1/user/recoverpassword","GET",recoverpassword)

	//下线注销
	//curl -i -H Accept:application/json -X GET http://localhost:8080/v1/user/logoff?token=<token>&id=<uid>
	router.Add("logoff","/v1/user/logoff","GET",logoff)

}


func getalluser(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "you are get all user.")
}

func getuser(ctx *fasthttp.RequestCtx) {
	uid := ctx.UserValue("uid")
	fmt.Fprintf(ctx, "you are get user %s", uid)
}

func modifyuser(ctx *fasthttp.RequestCtx) {
	uid := ctx.UserValue("uid")
	fmt.Fprintf(ctx, "you are modify user %s", uid)
}

func deleteuser(ctx *fasthttp.RequestCtx) {
	uid := ctx.UserValue("uid")
	fmt.Fprintf(ctx, "you are delete user %s", uid)
}
func adduser(ctx *fasthttp.RequestCtx) {
	uid := ctx.FormValue("uid")
	fmt.Fprintf(ctx, "you are add user %s", string(uid))
}

func getVerifyCode(ctx *fasthttp.RequestCtx) {

	tel := string(ctx.FormValue("tel"))
	if len(tel) == 0 {
		//无效的电话号码直接返回
		JsonErrorResult(ctx,ERROR_PARAMS_INVALID,"The tel is invalid")
		return
	}

	code := GenVerifyCode(tel)
	//todo 发送SMS

	JsonMsgResult(ctx, VerifyCode{code/*, JsonTime(time.Now())*/})
	log.Printf("tel = %s, vcode = %s\n",tel, code)
}

func register(ctx *fasthttp.RequestCtx) {

	//检查验证码
	tel := string(ctx.FormValue("tel"))
	verifyCode := string(ctx.FormValue("verifycode"))
	if !CheckVerifyCode(tel,verifyCode) {
		JsonErrorResult(ctx, ERROR_PARAMS_INVALID,"Error VerifyCode")
		return
	}

	//获取其他参数
	param := make(map[string]interface{})
	param["un"] = string(ctx.FormValue("un"))
	pwd := string(ctx.FormValue("pwd"))
	param["pwd"] = password(pwd)
	intGender,err:=strconv.Atoi(string(ctx.FormValue("gender")))
	if err != nil {
		intGender = int(modules.GENDER_NULL)
	}else if intGender < int(modules.GENDER_NULL) || intGender > int(modules.GENDER_FEMALE) {
		intGender = int(modules.GENDER_NULL)
	}

	param["gender"] = intGender
	param["nick"] = string(ctx.FormValue("nick"))
	param["id_number"] = string(ctx.FormValue("id_number"))
	strBirthday := string(ctx.FormValue("birthday"))
	if len(strBirthday) == 0 {
		param["birthday"] = sql.NullString{}
	}else {
		param["birthday"] = strBirthday
	}

	param["last_access"] = time.Now()
	param["tel"] = tel

	//参数检查
	err = checkUserInfoParam(param)
	if  err != nil {
		JsonErrorResult(ctx, ERROR_PARAMS_INVALID,err.Error())
		return
	}

	//获取以太坊地址
	param["eth_account"], err = getEthAccount()
	if err != nil {
		JsonErrorResult(ctx, ERROR_ETHEREUM,err.Error())
		return
	}

	userModule := modules.NewModel("user_info")
	_, err = userModule.Insert(param)
	if err != nil {
		JsonErrorResult(ctx, ERROR_SQL,err.Error())
		return
	}

	JsonSuccResult(ctx)
}

func password(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)) )
}

func checkUserInfoParam(param map[string]interface{}) error {
	if len(param["tel"].(string)) == 0 {
		return fmt.Errorf("The tel can't null")
	}

	userModule := modules.NewModel("user_info")

	where := "tel='"+param["tel"].(string)+"'"
	count := userModule.Where(where).Count()
	if count > 0 {
		return fmt.Errorf("The phone number already exists")
	}

	where = "un='"+param["un"].(string)+"'"
	count = userModule.Where(where).Count()
	if count > 0 {
		return fmt.Errorf("The username already exists")
	}

	return nil
}

func getEthAccount() (string, error) {
	//todo 获取以太坊地址
	return "12345678901234567890", nil
}

func login(ctx *fasthttp.RequestCtx) {
	un := string(ctx.FormValue("un"))
	pwd := string(ctx.FormValue("pwd"))
	//防注入
	if strings.Contains(un, "'") {
		JsonErrorResult(ctx, ERROR_PARAMS_INVALID,"The param is invalid.")
		return
	}

	if len(un) == 0 || len(pwd) == 0 {
		JsonErrorResult(ctx, ERROR_PARAMS_INVALID,"The param is invalid.")
		return
	}

	//查找账户
	userModule := modules.NewModel("user_info")
	where := fmt.Sprintf("un='%s' or tel='%s' and pwd = '%s'", un, un, password(pwd))
	user := userModule.Fileds("id").Where(where).Limit(1).Get()
	if len(user) < 1 {
		JsonErrorResult(ctx, ERROR_PASSWORD,"username or password is incorrect")
		return
	}
	id, err := strconv.ParseInt(user[1]["id"], 10, 64)
	if err != nil {
		JsonErrorResult(ctx, ERROR_SQL,"Error query result.")
		return
	}

	accessModule := modules.NewModel("user_access")
	if accessModule.Where("user_id="+user[1]["id"]+" and online_state=1").Count() > 0 {
		JsonErrorResult(ctx, ERROR_ALREADY_ONLINE,"Already online.")
		return
	}

	var values = make(map[string]interface{})
	//登录成功刷新登录时间
	values["last_access"] = time.Now()
	where = fmt.Sprintf("id='%s'", user[1]["id"])
	userModule.Where(where).Update(values)


	//更新token
	token := getToken(id)

	//更新登录表
	accessModule = modules.NewModel("user_access")
	values = make(map[string]interface{})
	values["user_id"] = id
	values["ip"] = ctx.RemoteIP().String()
	//todo ip位置查询
	values["loc"] = "todo ip位置查询"
	values["token"] = token
	values["last_access"] = time.Now()
	values["online_state"] = int(modules.ONLINESTATE_ON)
	if accessModule.Where("user_id="+user[1]["id"]).Count() > 0 {
		//delete(values, "user_id")
		accessModule.Update(values)
	}else {
		accessModule.Insert(values)
	}

	JsonMsgResult(ctx, LoginSucc{id, token})
}

func getToken( id int64 ) string {
	m := md5.Sum([]byte(fmt.Sprintf("%d%s", id, vcode())))
	return hex.EncodeToString( m[:] )
}

func recoverpassword(ctx *fasthttp.RequestCtx) {
	//tel	required	用户手机号
	//password string	required 用户新密码
 	//verifyCode	string	required	验证码

	//检查验证码
	tel := string(ctx.FormValue("tel"))
	verifyCode := string(ctx.FormValue("verifycode"))
	if !CheckVerifyCode(tel,verifyCode) {
		JsonErrorResult(ctx, ERROR_PARAMS_INVALID,"Error VerifyCode")
		return
	}

	//获取其他参数
	values := make(map[string]interface{})
	pwd := string(ctx.FormValue("pwd"))
	values["pwd"] = password(pwd)
	values["last_access"] = time.Now()

	userModule := modules.NewModel("user_info")
	_, err := userModule.Where("tel='"+tel+"'").Update(values)
	if err != nil {
		JsonErrorResult(ctx, ERROR_SQL,err.Error())
		return
	}

	JsonSuccResult(ctx)
}


func changeinfo(ctx *fasthttp.RequestCtx) {
	//token	required    用户token
	//id	required	用户id

	//gender option  性别
	//nick   option  昵称
	//id_number option 身份证
	//birthday   option  生日

	intGender,err:=strconv.Atoi(string(ctx.FormValue("gender")))
	if err != nil {
		intGender = int(modules.GENDER_NULL)
	}else if intGender < int(modules.GENDER_NULL) || intGender > int(modules.GENDER_FEMALE) {
		intGender = int(modules.GENDER_NULL)
	}

	values := make(map[string]interface{})
	values["gender"] = intGender
	values["nick"] = string(ctx.FormValue("nick"))
	values["id_number"] = string(ctx.FormValue("id_number"))
	strBirthday := string(ctx.FormValue("birthday"))
	if len(strBirthday) == 0 {
		values["birthday"] = sql.NullString{}
	}else {
		values["birthday"] = strBirthday
	}
	if len(values) > 0 {

	}
}

func changepassword(ctx *fasthttp.RequestCtx) {
	//token	required    用户token
	//id	required	用户id
	//password string	required 用户密码
	//newPassword	string	required	用户新密码

	token := string(ctx.FormValue("token"))
	id := string(ctx.FormValue("id"))
	passwd := string(ctx.FormValue("password"))
	newpassword := string(ctx.FormValue("newpassword"))

	if len(token) == 0 || len(id) == 0 || len(passwd) == 0 || len(newpassword) == 0{
		JsonErrorResult(ctx, ERROR_PARAMS_INVALID,"Parameter token and id are required")
		return
	}

	accessModule := modules.NewModel("user_access")

	if accessModule.Where("token='"+token+"' and user_id="+id+" and online_state = 1").Count() < 1 {
		JsonErrorResult(ctx, ERROR_NOT_ONLINE,"User is not online")
		return
	}

	userModel := modules.NewModel("user_info")
	if userModel.Where("id="+id+" and pwd='"+password(passwd)+"'").Count() < 1 {
		JsonErrorResult(ctx, ERROR_PASSWORD,"username or password is incorrect")
		return
	}

	values := make(map[string]interface{})
	values["pwd"] = password(newpassword)
	userModel.Update(values)
	JsonSuccResult(ctx)
}

func logoff(ctx *fasthttp.RequestCtx) {
	//token	required    用户token
	//id	required	用户id
	token := string(ctx.FormValue("token"))
	id := string(ctx.FormValue("id"))

	if len(token) == 0 || len(id) == 0 {
		JsonErrorResult(ctx, ERROR_PARAMS_INVALID,"Parameter token and id are required")
		return
	}

	accessModule := modules.NewModel("user_access")

	if accessModule.Where("token='"+token+"' and user_id="+id+" and online_state = 1").Count() < 1 {
		JsonErrorResult(ctx, ERROR_NOT_ONLINE,"User is not online")
		return
	}

	values := make(map[string]interface{})
	values["online_state"] = int(modules.ONLINESTATE_OFF)
	accessModule.Update(values)

	JsonSuccResult(ctx)
}
