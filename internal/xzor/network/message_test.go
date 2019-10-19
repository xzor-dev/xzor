package network_test

import (
	"testing"

	"github.com/xzor-dev/xzor/internal/xzor/network"
)

func TestMessages(t *testing.T) {
	type testDataA struct {
		Foo string
	}

	type testDataB struct {
		Bar string
	}

	data := &testDataA{
		Foo: "bar",
	}
	msg, err := network.NewMessage(data)
	if err != nil {
		t.Fatalf("%v", err)
	}
	encoded, err := msg.Encode()
	if err != nil {
		t.Fatalf("%v", err)
	}
	dataB := &testDataB{}
	_, err = encoded.Decode(dataB)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if dataB.Bar != "" {
		t.Fatalf("expected message data to be empty")
	}

	dataA := &testDataA{}
	_, err = encoded.Decode(dataA)
	if dataA.Foo != data.Foo {
		t.Fatalf("expected dataA to have Foo value of %s, got %s", data.Foo, dataA.Foo)
	}
}
