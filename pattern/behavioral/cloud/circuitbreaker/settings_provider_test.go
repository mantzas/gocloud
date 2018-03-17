package circuitbreaker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalSettingsProvider_New(t *testing.T) {

	assert := assert.New(t)

	p := NewLocalSettingsProvider()

	assert.NotNil(p)
}

func TestLocalSettingsProvider_Get(t *testing.T) {

	assert := assert.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "123"})

	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{"Success", "123", false},
		{"Failure", "456", true},
	}

	for _, tt := range tests {

		sett, err := pr.Get(tt.key)

		if tt.wantErr {
			assert.NotNil(err, "should not be nil")
			assert.Nil(sett, "should be nil but was %v", sett)
		} else {
			assert.Nil(err, "should be nil but was %v", err)
			assert.NotNil(sett, "should not be nil")
		}
	}
}

func TestLocalSettingsProvider_GetKeys(t *testing.T) {

	assert := assert.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "123"})
	pr.Save(Setting{Key: "456"})

	keys := pr.GetKeys()

	assert.Contains(keys, "123")
	assert.Contains(keys, "456")
}

func TestLocalSettingsProvider_Save(t *testing.T) {

	assert := assert.New(t)

	p := NewLocalSettingsProvider()
	p.Save(Setting{Key: "123"})

	assert.NotNil(p.store["123"])
}
