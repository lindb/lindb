package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode_Indicator(t *testing.T) {
	node := &Node{IP: "1.1.1.1", Port: 19000}
	indicator := node.Indicator()
	assert.Equal(t, "1.1.1.1:19000", indicator)
	node2, err := ParseNode(indicator)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, *node, *node2)
}

func TestNode_ParseNode(t *testing.T) {
	if _, err := ParseNode("xxx:123"); err == nil {
		t.Fatal("should be error")
	}

	if _, err := ParseNode("1.1.1.1123"); err == nil {
		t.Fatal("should be error")
	}

	if _, err := ParseNode("1.1.1.1:-1"); err == nil {
		t.Fatal("should be error")
	}

	if _, err := ParseNode("1.1.1.1:65536"); err == nil {
		t.Fatal("should be error")
	}

	node, err := ParseNode("1.1.1.1:65535")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, node.IP, "1.1.1.1")
	assert.Equal(t, node.Port, uint16(65535))

	if _, err = ParseNode(":123"); err == nil {
		t.Fatal(err)
	}
}
