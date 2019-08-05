package svc_test

import (
	"fmt"
	"testing"

	"github.com/xzor-dev/xzor/svc"
)

func TestModuleRouter(t *testing.T) {
	mName := "my_module"
	m := NewTestModule(mName)
	r := &svc.ModuleRouter{}
	r.Register(m)

	testData := "Hello Module!"
	testMsg := fmt.Sprintf("%s%s%s", mName, svc.ModuleNameDelimiter, testData)
	err := r.Process([]byte(testMsg))
	if err != nil {
		t.Fatalf("failed to process message: %v", err)
	}

	go func() {
		data := m.Read()
		if string(data) != testData {
			t.Fatalf("unexpected data received from module: wanted %s, got %s", testData, data)
		}
	}()
}

func NewTestModule(name string) *testModule {
	return &testModule{
		name: svc.ModuleName(name),
	}
}

var _ svc.Module = &testModule{}

type testModule struct {
	data chan []byte
	name svc.ModuleName
}

func (m *testModule) Name() svc.ModuleName {
	return m.name
}

func (m *testModule) Process(data []byte) error {
	if m.data == nil {
		m.data = make(chan []byte, 1)
	}
	m.data <- data
	return nil
}

func (m *testModule) Read() []byte {
	return <-m.data
}
