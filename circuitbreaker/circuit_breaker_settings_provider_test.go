package circuitbreaker

import (
	"reflect"
	"testing"
)

func TestNewLocalSettingsProvider(t *testing.T) {
	tests := []struct {
		name string
		want *LocalSettingsProvider
	}{
		{"Constructor", NewLocalSettingsProvider()},
	}
	for _, tt := range tests {
		if got := NewLocalSettingsProvider(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewLocalSettingsProvider() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestLocalSettingsProvider_Get(t *testing.T) {

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "123"})

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		lsp     *LocalSettingsProvider
		args    args
		want    *Setting
		wantErr bool
	}{
		{"GetKey Success", pr, args{"123"}, &Setting{Key: "123"}, false},
		{"GetKey Failure", pr, args{"234"}, nil, true},
	}
	for _, tt := range tests {
		got, err := tt.lsp.Get(tt.args.key)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. LocalSettingsProvider.Get() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. LocalSettingsProvider.Get() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestLocalSettingsProvider_GetKeys(t *testing.T) {

	pr := NewLocalSettingsProvider()
	pr.Save(Setting{Key: "123"})
	pr.Save(Setting{Key: "456"})

	tests := []struct {
		name string
		lsp  *LocalSettingsProvider
		want []string
	}{
		{"GetKeys", pr, []string{"123", "456"}},
	}
	for _, tt := range tests {
		if got := tt.lsp.GetKeys(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. LocalSettingsProvider.GetKeys() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestLocalSettingsProvider_Save(t *testing.T) {
	type args struct {
		sett Setting
	}
	tests := []struct {
		name string
		lsp  *LocalSettingsProvider
		args args
	}{
		{"Save", NewLocalSettingsProvider(), args{Setting{}}},
	}
	for _, tt := range tests {
		tt.lsp.Save(tt.args.sett)
	}
}
