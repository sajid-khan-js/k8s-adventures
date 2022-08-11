package k8sclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNamespaces(t *testing.T) {

	clientSet, err := InitClient()
	if err != nil {
		t.Error(err)
	}

	namespaces, err := GetNamespaces(*clientSet)
	if err != nil {
		t.Error(err)
	}
	assert.Contains(t, namespaces, "kube-system", "default")
}

// TODO write more tests for 100% coverage and use mocks
