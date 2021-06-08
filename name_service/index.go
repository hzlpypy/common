package name_service

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"runtime"
)

// GetNameService:获取对应组建的名字服务
func GetNameService() (map[string]string, error) {
	_, filename, _, _ := runtime.Caller(0) // get current filepath in runtime
	filepath := path.Dir(filename)
	by, err := ioutil.ReadFile(filepath + "/" + "name_service.json")
	if err != nil {
		return nil, err
	}
	res := make(map[string]string)
	err = json.Unmarshal(by, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
