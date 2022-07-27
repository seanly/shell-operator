package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type MutatingResponse struct {
	Allowed  bool     `json:"allowed"`
	Message  string   `json:"message,omitempty"`
	Patch    []byte   `json:"patch,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

func MutatingResponseFromFile(filePath string) (*MutatingResponse, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %s", filePath, err)
	}

	if len(data) == 0 {
		return nil, nil
	}
	return MutatingResponseFromBytes(data)
}

func MutatingResponseFromBytes(data []byte) (*MutatingResponse, error) {
	return MutatingResponseFromReader(bytes.NewReader(data))
}

func MutatingResponseFromReader(r io.Reader) (*MutatingResponse, error) {
	response := new(MutatingResponse)

	dec := json.NewDecoder(r)

	err := dec.Decode(response)

	if err != nil {
		return nil, err
	}

	return response, nil
}


func (r *MutatingResponse) Dump() string {
	b := new(strings.Builder)
	b.WriteString("MutatingResponse(allowed=")
	b.WriteString(strconv.FormatBool(r.Allowed))
	if r.Message != "" {
		b.WriteString(",")
		b.WriteString(r.Message)
	}
	for _, warning := range r.Warnings {
		b.WriteString(",")
		b.WriteString(warning)
	}
	b.WriteString(")")
	return b.String()
}