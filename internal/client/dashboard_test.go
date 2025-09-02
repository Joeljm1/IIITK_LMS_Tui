package client

import "testing"

func TestDashBoard(t *testing.T) {
	lms, err := NewClient("2023BCS0061", "RM76YXMQ")
	if err != nil {
		t.Fatal("Could not make client")
	}
	dash, err := lms.GetDashBoard()
	if err != nil {
		t.Fatalf("Could not get dashboard, err:%v", err.Error())
	}
	if dash[0].Error {
		t.Fatal("Could not get dashboard2")
	}
	t.Log(dash[0].Data.Events[0].Name, "\n")
	t.Log(dash[0].Data.Events[0].Course, "\n")
}
