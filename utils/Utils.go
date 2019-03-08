package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func RequireU64(name string, r *http.Request, w http.ResponseWriter) (uint64, bool) {
	if r.URL.Query().Get(name) == "" {
		return 0, false
	} else {
		val, err := strconv.ParseInt(r.URL.Query().Get(name), 10, 64)
		if err != nil {
		}
		return uint64(val), err == nil
	}
}

func SendResponse(resp interface{}, w http.ResponseWriter) {
	res, er := json.Marshal(resp)
	if er != nil {
		panic("Error marshalling response: " + er.Error())
	}
	_, _ = fmt.Fprint(w, string(res))
}

func ParseJsonConfig(file string) (map[string]interface{}, error) {
	// Open config
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	return result, err
}

func MustParseJsonConfig(file string) map[string]interface{} {
	// Open config
	jsonFile, err := os.Open("configs/" + file)
	if err != nil {
		panic("[ConfigParser] Error loading config: " + err.Error())
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		panic("[ConfigParser] Error loading config: " + err.Error())
	}

	return result
}
