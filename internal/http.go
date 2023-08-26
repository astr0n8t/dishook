package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
)

func (w *WebhookSlashCommand) getPayloadWithArguments(data string) ([]byte, error) {
	args := make(map[string]interface{})

	for _, arg := range w.Arguments {
		for _, calledArg := range w.CalledOptions {
			if arg.Name == calledArg.Name {
				args[arg.Name] = calledArg.Value
			}
		}
	}

	// Try to template our data
	tmpl, err := template.New(w.Name).Parse(data)
	if err != nil {
		return nil, fmt.Errorf("creating payload template: %v", err)
	}

	var tmplBuf bytes.Buffer
	err = tmpl.Execute(&tmplBuf, args)
	if err != nil {
		return nil, fmt.Errorf("executing template payload: %v", err)
	}

	return tmplBuf.Bytes(), nil
}

func (w *WebhookSlashCommand) getJSONPaylod() ([]byte, error) {
	// No payload
	if w.Data == nil {
		return nil, nil
	}

	jsonBytes, err := json.Marshal(w.Data)
	// If we cannot marshal the json return an error
	if err != nil {
		return nil, fmt.Errorf("marshalling JSON payload: %v", err)
	}
	if len(w.Arguments) > 0 {
		jsonBytes, err = w.getPayloadWithArguments(string(jsonBytes))
		if err != nil {
			return nil, fmt.Errorf("templating JSON payload: %v", err)
		}
	}

	return jsonBytes, nil
}

func (w *WebhookSlashCommand) request() error {
	payload, err := w.getJSONPaylod()
	// If we have an issue getting the payload just return the error
	if err != nil {
		return fmt.Errorf("constructing JSON payload: %v", err)
	}
	client := &http.Client{}

	method := w.Method
	if method == "" {
		method = "POST"
	}

	req, err := http.NewRequest(method, w.URL, bytes.NewBuffer(payload))

	for _, header := range w.Headers {
		req.Header.Add(header.Name, header.Value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling http client request: %v", err)
	}

	if w.RespCode != 0 && resp.StatusCode != w.RespCode {
		return fmt.Errorf("requested URL returned response %d, expected %d", resp.StatusCode, w.RespCode)
	} else if w.RespCode == 0 && resp.StatusCode != 200 {
		return fmt.Errorf("requested URL returned response %d, expected %d", resp.StatusCode, 200)
	}

	return nil
}
