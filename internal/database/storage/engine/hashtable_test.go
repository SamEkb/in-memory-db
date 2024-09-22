package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHashtable(t *testing.T) {
	table := NewHashtable()

	require.NotNil(t, table)
	assert.NotNil(t, table.data)

}

func TestHashtable_Get(t *testing.T) {
	table := &Hashtable{
		data: map[string]string{
			"1": "1",
			"2": "2",
			"3": "3",
		},
	}

	tests := map[string]struct {
		key           string
		expectedValue string
		exists        bool
	}{
		"get by existing key": {
			key:           "1",
			expectedValue: "1",
			exists:        true,
		},
		"get by non existing key": {
			key:           "4",
			expectedValue: "",
			exists:        false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value, ok := table.Get(test.key)
			assert.Equal(t, value, test.expectedValue)
			assert.Equal(t, ok, test.exists)
		})
	}
}

func TestHashtable_Insert(t *testing.T) {
	table := &Hashtable{
		data: map[string]string{},
	}

	tests := map[string]struct {
		key   string
		value string
	}{
		"insert key and value": {
			key:   "1",
			value: "1",
		},
		"insert existing key": {
			key:   "2",
			value: "0",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			key := test.key
			table.Insert(key, test.value)
			value, ok := table.Get(key)
			assert.Equal(t, test.value, value)
			assert.True(t, ok, value)
		})
	}
}

func TestHashtable_Del(t *testing.T) {
	table := &Hashtable{
		data: map[string]string{
			"1": "1",
			"2": "2",
			"3": "3",
		},
	}

	tests := map[string]struct {
		key      string
		initial  string
		afterDel string
		exists   bool
	}{
		"delete existing key": {
			key:      "1",
			initial:  "1",
			afterDel: "",
			exists:   true,
		},
		"delete non-existing key": {
			key:      "4",
			initial:  "",
			afterDel: "",
			exists:   false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value, ok := table.Get(test.key)
			assert.Equal(t, test.initial, value)
			assert.Equal(t, test.exists, ok)

			table.Del(test.key)

			value, ok = table.Get(test.key)
			assert.Equal(t, test.afterDel, value)
			assert.False(t, ok)
		})
	}
}
