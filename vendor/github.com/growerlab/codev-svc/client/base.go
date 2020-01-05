package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

func Post(client *http.Client, apiURL string, bodyMap map[string]interface{}) (*Result, error) {
	bodyRaw, _ := json.Marshal(bodyMap)

	resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(bodyRaw))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("status: %v  body: %v", resp.StatusCode, string(respBody))
	}
	result, err := BuildResult(respBody)
	return result, errors.WithStack(err)
}

func BuildResult(respBody []byte) (*Result, error) {
	result := new(Result)
	err := json.Unmarshal(respBody, result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(result.Errors) > 0 {
		return nil, errors.New(result.Errors[0].Message)
	}
	result.DataPath = gjson.ParseBytes(result.Data)
	return result, nil
}
