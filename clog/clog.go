package clog

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

//cfg 配置文件
type Cfg struct {
	Host       string
	IndexName  string
	CfgFiles   []*CfgFile // access:xxx1.txt error:xxx2.txt
	EsAdders   []string   //ES addr
	EsUser     string     //ES user
	EsPassword string     //ES password
}

type CfgFile struct {
	Name string
	File *os.File
}

// access error info
type LogCategory struct {
	log        *logrus.Logger
	LogNameMap map[string]*os.File
}

//Init 初始化日志配置信息logrus
func Init(c *Cfg) (*LogCategory, error) {
	//if !utils.RegexpIPV4(d.Host) {
	//	return errors.New("Fail: Please check if the host is correct")
	//}
	res := &LogCategory{LogNameMap: make(map[string]*os.File)}
	l := logrus.New()
	validateLogName := map[string]bool{"access": true, "error": true, "info": true}
	for _, cfgFile := range c.CfgFiles {
		if _, ok := validateLogName[cfgFile.Name]; !ok {
			return nil, errors.New("log name is invalid")
		}
		res.LogNameMap[cfgFile.Name] = cfgFile.File
	}
	err := setupLogrus(l, c)
	if err != nil {
		return nil, err
	}
	res.log = l
	return res, nil
}

func (l *LogCategory) Access() *logrus.Logger {
	l.log.SetOutput(l.LogNameMap["access"])
	logLvl, _ := logrus.ParseLevel("info")
	l.log.SetLevel(logLvl)
	return l.log
}

func (l *LogCategory) Error() *logrus.Logger {
	l.log.SetOutput(l.LogNameMap["error"])
	logLvl, _ := logrus.ParseLevel("error")
	l.log.SetLevel(logLvl)
	return l.log
}

func (l *LogCategory) Info() *logrus.Logger {
	l.log.SetOutput(l.LogNameMap["info"])
	logLvl, _ := logrus.ParseLevel("info")
	l.log.SetLevel(logLvl)
	return l.log
}

//setupLogrus 初始化logrus 同时把logrus的logger var 引用到这个common.Logger
func setupLogrus(l *logrus.Logger, c *Cfg) error {
	l.SetReportCaller(true)
	//l.SetOutput(o)
	//开启 logrus ES hook
	// esh := newEsHook(cc)
	// logrus.AddHook(esh)
	// fmt.Printf(">= Error level, check the logrus* index in the log %s", SetInit.Host)

	return nil
}

//esHook 自定义的ES hook
type esHook struct {
	cmd       string // 记录启动的命令
	client    *elastic.Client
	indexName string
}

//newEsHook 初始化
func newEsHook(cc *Cfg) (*esHook, error) {
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
	return &esHook{client: es, cmd: strings.Join(os.Args, " "), indexName: cc.IndexName}, nil
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
	_, err := hook.client.Index().Index(doc.indexName(hook.indexName)).Type("_doc").BodyJson(doc).Do(context.Background())
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
func (m *appLogDocModel) indexName(indexName string) string {
	return indexName + "-" + time.Now().Local().Format("2006-01-02")
}
