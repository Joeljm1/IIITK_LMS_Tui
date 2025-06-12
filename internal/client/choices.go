package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
)

const ChoiceFile = "courses.json"

type Choices struct {
	AllChoices   CourseList        `json:"allChoices"`
	AttendanceId map[string]string `json:"attendanceId"`
}

func writeToFile(filename string, val any) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, b, 0644)
	return err
}

func DeleteCoursesFile() error {
	err := os.Remove(ChoiceFile)
	return err
}

func (lmsCLient *LMSCLient) ChoicesInFile() (*Choices, error) {
	b, err := os.ReadFile(ChoiceFile)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return nil, errors.New("empty json")
	}
	br := bytes.NewReader(b)
	dec := json.NewDecoder(br)
	dec.DisallowUnknownFields()
	var choices *Choices
	err = dec.Decode(&choices)
	if choices == nil {
		return nil, errors.New("empty json")
	}
	if err != nil {
		return nil, err
	}
	return choices, nil
}
