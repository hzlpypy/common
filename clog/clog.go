package clog

import (
	"context"
	"errors"
	"fmt"
	"github.com/hzlpypy/common/utils"
	"log"
	"os"
	"strings"
	"time"

	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

//cfg 配置文件
type cfg struct {
	LogLvl     string   // 日志级别
	EsAdders   []string //ES addr
	EsUser     string   //ES user
	EsPassword string   //ES password
}

//SetDefault 默认初始化参数
type SetDefault struct {
	Host      string
	IndexName string
	LogLvl    string
}

//SetInit 默认初始化参数设置
var SetInit SetDefault

//Init 初始化日志配置信息logrus
//	param host http://10.0.1.77:9021/
func Init(d *SetDefault) error {
	if !utils.RegexpIPV4(d.Host) {
		return errors.New("Fail: Please check if the host is correct")
	}

	SetInit.Host = d.Host
	SetInit.IndexName = d.IndexName
	var level = "error"
	if d.LogLvl != "" {
		level = d.LogLvl
	}
	cc := cfg{
		LogLvl:     level,
		EsAdders:   []string{SetInit.Host},
		EsUser:     "",
		EsPassword: "",
	}
	err := setupLogrus(cc)
	if err != nil {
		return err
	}
	return nil
}

//setupLogrus 初始化logrus 同时把logrus的logger var 引用到这个common.Logger
func setupLogrus(cc cfg) error {
	logLvl, err := logrus.ParseLevel(cc.LogLvl)
	if err != nil {
		return err
	}
	logrus.SetLevel(logLvl)
	logrus.SetReportCaller(true)
	//开启 logrus ES hook
	// esh := newEsHook(cc)
	// logrus.AddHook(esh)
	// fmt.Printf(">= Error level, check the logrus* index in the log %s", SetInit.Host)

	return nil
}

//esHook 自定义的ES hook
type esHook struct {
	cmd    string // 记录启动的命令
	client *elastic.Client
}

//newEsHook 初始化
func newEsHook(cc cfg) (*esHook, error) {
	es, err := elastic.NewClient(
		elastic.SetURL(cc.EsAdders...),
		elastic.SetBasicAuth(cc.EsUser, cc.EsPassword),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(15*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "ES:", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "ES:", log.LstdFlags)),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create Elastic V6 Client: %v", err)
	}
	return &esHook{client: es, cmd: strings.Join(os.Args, " ")}, nil
}

//Fire logrus hook interface 方法
func (hook *esHook) Fire(entry *logrus.Entry) error {
	doc := newEsLog(entry)
	doc["cmd"] = hook.cmd
	// go hook.sendEs(doc)
	return nil
}

//Levels logrus hook interface 方法
func (hook *esHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
}

//sendEs 异步发送日志到es
func (hook *esHook) sendEs(doc appLogDocModel) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("send entry to es failed: ", r)
		}
	}()
	_, err := hook.client.Index().Index(doc.indexName()).Type("_doc").BodyJson(doc).Do(context.Background())
	if err != nil {
		log.Println(err)
	}

}

//appLogDocModel es model
type appLogDocModel map[string]interface{}

func newEsLog(e *logrus.Entry) appLogDocModel {
	ins := map[string]interface{}{}
	for kk, vv := range e.Data {
		ins[kk] = vv
	}
	ins["time"] = time.Now().Local()
	ins["lvl"] = e.Level
	ins["message"] = e.Message
	ins["caller"] = fmt.Sprintf("%s:%d  %#v", e.Caller.File, e.Caller.Line, e.Caller.Func)
	return ins
}

// indexName es index name 时间分割
func (m *appLogDocModel) indexName() string {
	return SetInit.IndexName + "-" + time.Now().Local().Format("2006-01-02")
}
