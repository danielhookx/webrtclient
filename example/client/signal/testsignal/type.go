package testsignal

import "encoding/json"

type Msg struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

type RegisterMsg struct {
	Name string `json:"name"`
}

func Parse(data []byte) (*Msg, error) {
	msg := Msg{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, err
}

func ParseRegisterMsg(data []byte) (*RegisterMsg, error) {
	msg := RegisterMsg{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, err
}
