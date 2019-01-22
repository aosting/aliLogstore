package cocolog

import (
	"github.com/aliyun/aliyun-log-go-sdk"
	"strings"
	"time"
	"errors"
)

/**********************************************
** @Des: cocoLogStoreConfig
** @Author: zhangxueyuan 
** @Date:   2019-01-22 11:17:49
** @Last Modified by:   zhangxueyuan 
** @Last Modified time: 2019-01-22 11:17:49
***********************************************/

type logStoreConfig struct {
	Project      *sls.LogProject //日志配置项
	CacheLength  int             //缓存大小
	logStorename string
	logstore     *sls.LogStore
}

func initlogStoreCofig(p *sls.LogProject, cacheSize int, name string) (logStoreConfig, error) {
	l := logStoreConfig{p, cacheSize, name, nil}
	var retry_times int
	var err error
	var logstore *sls.LogStore
	for retry_times = 0; ; retry_times++ {
		if retry_times > 3 {
			return l, errors.New("get LogStore error")
		}
		logstore, err = p.GetLogStore(name)
		if err != nil {
			WARN("GetLogStore fail, retry:%d, err:%v\n", retry_times, err)
			if strings.Contains(err.Error(), sls.PROJECT_NOT_EXIST) {
				return l, err
			} else if strings.Contains(err.Error(), sls.LOGSTORE_NOT_EXIST) {
				err = p.CreateLogStore(name, 7, 2, true, 64)
				if err != nil {
					WARN("CreateLogStore fail, err: ", err.Error())
				} else {
					INFO("CreateLogStore success")
				}
				logstore, err = p.GetLogStore(name)
				if err == nil {
					l.logstore = logstore
					break
				}
			}
		} else {
			l.logstore = logstore
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	return l, nil
}
