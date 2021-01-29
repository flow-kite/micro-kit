package test

import "testing"

func TestGet(t *testing.T) {
	instance := Get()
	for key, value := range instance.DBs {
		t.Logf("key = %v and value = %v", key, value)
	}
}
