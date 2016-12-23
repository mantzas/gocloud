package cloud

import "fmt"

// SettingsRetriever interface
type SettingsRetriever interface {
	Get(key string) (*Setting, error)
	GetKeys() []string
}

// SettingsSaver interface
type SettingsSaver interface {
	Save(sett Setting)
}

// LocalSettingsProvider definition
type LocalSettingsProvider struct {
	store map[string]Setting
}

// NewLocalSettingsProvider constructor
func NewLocalSettingsProvider() *LocalSettingsProvider {
	return &LocalSettingsProvider{make(map[string]Setting, 0)}
}

// Get the setting based on the key from the local store
func (lsp *LocalSettingsProvider) Get(key string) (*Setting, error) {
	sett, ok := lsp.store[key]
	if !ok {
		return nil, fmt.Errorf("%s not found", key)
	}

	return &sett, nil
}

// GetKeys returns the keys of the local store
func (lsp *LocalSettingsProvider) GetKeys() []string {

	var keys []string

	for key := range lsp.store {
		keys = append(keys, key)
	}

	return keys
}

// Save the setting to the local store
func (lsp *LocalSettingsProvider) Save(sett Setting) {
	lsp.store[sett.Key] = sett
}
