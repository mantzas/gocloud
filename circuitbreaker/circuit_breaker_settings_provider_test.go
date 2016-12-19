package circuitbreaker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalSettingsProvider_New(t *testing.T) {

	require := require.New(t)

	p := NewLocalSettingsProvider()

	require.NotNil(p)
}

func TestLocalSettingsProvider_Get(t *testing.T) {

	require := require.New(t)

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
			require.NotNil(err, "should not be nil")
			require.Nil(sett, "should be nil but was %v", sett)
		} else {
			require.Nil(err, "should be nil but was %v", err)
			require.NotNil(sett, "should not be nil")
		}
	}
}

func TestLocalSettingsProvider_GetKeys(t *testing.T) {

	require := require.New(t)

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "123"})
	pr.Save(Setting{Key: "456"})

	keys := pr.GetKeys()

	require.Contains(keys, "123")
	require.Contains(keys, "456")
}

func TestLocalSettingsProvider_Save(t *testing.T) {

	require := require.New(t)

	p := NewLocalSettingsProvider()
	p.Save(Setting{Key: "123"})

	require.NotNil(p.store["123"])
}
