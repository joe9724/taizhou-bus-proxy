package utils

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
)

func Response200(code int64,msg string) (string) {
	/*fmt.Println(code)
	fmt.Println(msg)
	responseStr := `{
    "code": 201,
    "msg": "ok"
    }`
	return responseStr*/
	var response Respoonse
	response.Code = code
	response.Msg = msg

	out, _ := json.Marshal(response)
	fmt.Println("ous is",out)
	return(string(out))

}

type Respoonse struct{
	Code int64 `json:"code,omitempty"`
	// cover
	Msg string `json:"msg,omitempty"`
}

func GetResourceDomain(filetype string) (string){
	var val string
	if(filetype == "icon") {
		val = "http://tingting-resource.bitekun.xin/resource/image/icon/"
	}else if(filetype == "cover"){
		val = "http://tingting-resource.bitekun.xin/resource/image/cover/"
	}else if(filetype == "m4a"){
		val = "http://tingting-resource.bitekun.xin/resource/mp3/"
	}else if(filetype == "amr"){
		val = "http://tingting-resource.bitekun.xin/resource/amr/"
	}else if(filetype == "other"){
		val = "http://tingting-resource.bitekun.xin/resource/other/"
	}
	return val
}

func OpenConnection() (db *gorm.DB, err error) {

	db, err = gorm.Open("mysql", "root:Tingtingyuedu654321!!!@tcp(47.104.131.147:3306)/tingting?charset=utf8&parseTime=True&loc=Local")
	return db,err

}