package cocolog

import (
	"time"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/golang/protobuf/proto"
	"os"
	"strings"
	"fmt"
	"errors"
	"sync"
)

/**********************************************
** @Des: cocoService
** @Author: zhangxueyuan 
** @Date:   2019-01-22 12:19:46
** @Last Modified by:   zhangxueyuan 
** @Last Modified time: 2019-01-22 12:19:46
***********************************************/
const PUSH_TRY_NUMS = 10

type LogService struct {
	rw              sync.RWMutex
	lsc             *logStoreConfig
	wlog            *wrap
	wrapinitcapcity int
	cache           int
	update          int64
}

func InitlogStore(p *sls.LogProject, cacheSize int, name string) (*LogService, error) {

	ilog, err := initlogStoreCofig(p, cacheSize, name)
	if err != nil {
		return nil, errors.New("initlogStoreCofig error!")
	}
	return &LogService{
		rw:              sync.RWMutex{},
		lsc:             &ilog,
		wlog:            initwrap(cacheSize),
		wrapinitcapcity: cacheSize,
		cache:           cacheSize,
		update:          time.Now().Unix(),
	}, nil
}

func (logService *LogService) Push(param map[string]string) {
	if logService == nil {
		WARN("logService is  error,do nothing!")
		return
	}

	logService.rw.RLock()
	wlog := logService.wlog
	cachelimit := logService.cache
	logService.rw.RUnlock()

	content := make([]*sls.LogContent, 0, len(param))
	for key := range param {
		content = append(content, &sls.LogContent{
			Key:   proto.String(key),
			Value: proto.String(param[key]),
		})
	}
	log := &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: content,
	}
	DEBUG(" Push map!")
	wlog.append(log)

	if wlog.size() >= cachelimit {
		DEBUG(" arrary cachelimit!")
		logService.rw.Lock()
		logService.update = time.Now().Unix()
		pushloghub := wlog.copy(logService.wrapinitcapcity)
		go logService.pushLogStore(pushloghub)

		logService.wlog = wlog
		logService.rw.Unlock()
	}
}

//日志上传
func (logService *LogService) pushLogStore(w *wrap) {
	if logService == nil {
		WARN("logService is  error,do nothing!")
		return
	}
	host, _ := os.Hostname()
	loggroup := &sls.LogGroup{
		Topic:  proto.String(""),
		Source: proto.String(host),
		Logs:   w.put,
	}

	var retry_times int
	// PostLogStoreLogs API Ref: https://intl.aliyun.com/help/doc-detail/29026.htm
	for retry_times = 0; retry_times < PUSH_TRY_NUMS; retry_times++ {
		err := logService.lsc.logstore.PutLogs(loggroup)
		if err == nil {
			DEBUG("PutLogs success, retry: " + fmt.Sprint(retry_times))
			break
		} else {
			msg := err.Error()
			//handle exception here, you can add retryable erorrCode, set appropriate put_retry
			if strings.Contains(msg, sls.POST_BODY_TOO_LARGE) {
				hw1, hw2 := w.sub()
				logService.pushLogStore(&hw1)
				logService.pushLogStore(&hw2)
				WARN("sub logs")
				break
			} else if strings.Contains(msg, sls.WRITE_QUOTA_EXCEED) || strings.Contains(msg, sls.PROJECT_QUOTA_EXCEED) || strings.Contains(err.Error(), sls.SHARD_WRITE_QUOTA_EXCEED) {
				//mayby you should split shard
				time.Sleep(500 * time.Millisecond)
			} else if strings.Contains(msg, sls.INTERNAL_SERVER_ERROR) || strings.Contains(msg, sls.SERVER_BUSY) {
				time.Sleep(200 * time.Millisecond)
			}
			if retry_times == (PUSH_TRY_NUMS - 1) {
				WARN("PutLogs fail, retry: "+fmt.Sprint(retry_times)+" , err:", msg)
				logService.pushLogStore(w)
			}
		}
	}
}

func (logService *LogService) Clear() {
	logService.rw.RLock()
	st := logService.update
	wlog := logService.wlog
	logService.rw.RUnlock()
	if time.Now().Unix()-st > 10 && wlog.size() > 0 {
		logService.rw.Lock()
		logService.update=time.Now().Unix()
		pushloghub := wlog.copy(logService.wrapinitcapcity)
		go logService.pushLogStore(pushloghub)
		logService.wlog = wlog
		logService.rw.Unlock()
	}
}
