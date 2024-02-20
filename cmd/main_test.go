package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartServer(t *testing.T) {
	server := newServer()

	assert.NotNil(t, server)
}
