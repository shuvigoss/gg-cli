package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ToResult(response *http.Response) (*Result, error) {
	var dataRaw json.RawMessage
	result := Result{
		Data: &dataRaw,
	}

	rsp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &result, err
	}

	err = json.Unmarshal(rsp, &result)
	if err != nil {
		return &result, err
	}

	return &result, nil
}

func ToData(result *Result, data interface{}) error {
	var dataRaw *json.RawMessage
	dataRaw = (result.Data).(*json.RawMessage)

	err := json.Unmarshal(*dataRaw, data)
	if err != nil {
		return err
	}
	return nil
}

func ToJsonString(data interface{}) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(data)
	if err != nil {
		return ""
	}
	return buffer.String()
}
