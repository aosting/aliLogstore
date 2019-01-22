package cocolog

import (
	"log"
)

/**********************************************
** @Des: cocolog
** @Author: zhangxueyuan 
** @Date:   2019-01-22 11:50:15
** @Last Modified by:   zhangxueyuan 
** @Last Modified time: 2019-01-22 11:50:15
***********************************************/
var (
	RELEASE = true
)

func SetRelease(isRelease bool) {
	RELEASE = isRelease
}

func DEBUG(v ...interface{}) {
	if !RELEASE {
		log.Println("DEBUG", "COCOLOG", v)
	}
}
func INFO(v ...interface{}) {
	log.Println("INFO", "COCOLOG", v)

}
func WARN(v ...interface{}) {
	log.Println("WARN", "COCOLOG", v)
}


func ERROR(v ...interface{}) {
	log.Println("ERROR", "COCOLOG", v)
}

func SUCC(v ...interface{}) {
	log.Println("SUCC", "COCOLOG", v)
}

func FAIL(v ...interface{}) {
	log.Println("FAIL", "COCOLOG", v)
}


