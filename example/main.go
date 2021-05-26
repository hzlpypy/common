package main

import (
	"github.com/hzlpypy/common/clog"
	"log"
	"os"
)

func main() {
	logPathMap := map[string]string{"access": "./log_access.txt", "error": "./log_error.txt"}
	logCfg := &clog.Cfg{}
	for name, path := range logPathMap {
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			file, _ = os.Create(path)
		}
		defer file.Close()
		logCfg.CfgFiles = append(logCfg.CfgFiles, &clog.CfgFile{
			Name: name,
			File: file,
		})
	}

	l, err := clog.Init(logCfg)
	if err != nil {
		log.Fatal(err)
	}
	l.Access().WithField("test", "info").Info("info")
	l.Error().WithField("test", "err123").Error("errwocuola")

}
