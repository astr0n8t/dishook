package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
)

// Processes arguments and templates the payload into a byte array
func (w *WebhookSlashCommand) getPayloadWithArguments(data string) ([]byte, error) {
	// Create a hashmap to store our arguments
	args := make(map[string]interface{})

	// Iterate over the provided arguments for the webhook
	for _, arg := range w.Arguments {
		// Check if this argument is provided through the discord information
		if arg.DiscordInfo {
			// Case where its the username
			// additional arguments can be added here
			if arg.Name == "discord_user_name" {
				// Get the caller's username from the struct
				args[arg.Name] = w.CalledUser.User.Username
			}
			// Otherwise arguments are provided in the command
		} else {
			// Iterate over the called arguments
			for _, calledArg := range w.CalledOptions {
				// Match the correct called argument with the webhook arguments
				if arg.Name == calledArg.Name {
					args[arg.Name] = calledArg.Value
				}
			}
			// Check if this argument was found in the called options
			_, valExists := args[arg.Name]
			// If this argument wasn't in the called options,
			// and its not a required argument, and a default
			// argument exists for this argument
			if !valExists && !arg.Req && arg.Default != nil {
				// Set the argument to the default value
				args[arg.Name] = arg.Default
			}
		}
	}

	// Try to create a template for the data
	tmpl, err := template.New(w.Name).Parse(data)
	if err != nil {
		return nil, fmt.Errorf("creating payload template: %v", err)
	}

	// Try to actually execute the template
	var tmplBuf bytes.Buffer
	err = tmpl.Execute(&tmplBuf, args)
	if err != nil {
		return nil, fmt.Errorf("executing template payload: %v", err)
	}

	// Return the templated data as a byte array
	return tmplBuf.Bytes(), nil
}

// Marshall the webhook data as a JSON payload into a byte array
// This will also get the arguments and template them
func (w *WebhookSlashCommand) getJSONPaylod() ([]byte, error) {
	// If there is no payload return nil
	if w.Data == nil {
		return nil, nil
	}

	// Try to marshall the data for the webhook
	jsonBytes, err := json.Marshal(w.Data)
	// If we cannot marshal the json return an error
	if err != nil {
		return nil, fmt.Errorf("marshalling JSON payload: %v", err)
	}
	// Check if there are arguments for this webhook
	if len(w.Arguments) > 0 {
		// If there are, template the payload with those arguments
		jsonBytes, err = w.getPayloadWithArguments(string(jsonBytes))
		if err != nil {
			return nil, fmt.Errorf("templating JSON payload: %v", err)
		}
	}

	// Return the JSON payload as a byte array
	return jsonBytes, nil
}

// Calls the given webhook
func (w *WebhookSlashCommand) request() error {
	// Get the JSON payload for this webhook
	payload, err := w.getJSONPaylod()
	// If we have an issue getting the payload just return the error
	if err != nil {
		return fmt.Errorf("constructing JSON payload: %v", err)
	}
	// Create a new HTTP client
	client := &http.Client{}

	// Try to set the HTTP method
	method := w.Method
	// If there is no method, assume POST
	if method == "" {
		method = "POST"
	}

	// Create our new request using our method, url, and payload
	req, err := http.NewRequest(method, w.URL, bytes.NewBuffer(payload))

	// Add any headers for the webhook
	for _, header := range w.Headers {
		req.Header.Add(header.Name, header.Value)
	}

	// Actually do the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling http client request: %v", err)
	}

	// Check if the reponse code is what we expected
	if w.RespCode != 0 && resp.StatusCode != w.RespCode {
		return fmt.Errorf("requested URL returned response %d, expected %d", resp.StatusCode, w.RespCode)
		// Default we assume a response code of 200
	} else if w.RespCode == 0 && resp.StatusCode != 200 {
		return fmt.Errorf("requested URL returned response %d, expected %d", resp.StatusCode, 200)
	}

	// If everything worked as intended return no error
	return nil
}
