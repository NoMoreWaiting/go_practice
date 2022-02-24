package gotemplate

import (
	"bytes"
	"strings"
	"testing"
)

func TestTranslateTemplate(t *testing.T) {
	type args struct {
		configPrefix string
	}
	tests := []struct {
		name   string
		args   args
		wantWr string
	}{
		{"case1", args{""}, "hello songyunxuan"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wr := &bytes.Buffer{}
			TranslateTemplate(wr, tt.args.configPrefix)
			if gotWr := wr.String(); !strings.Contains(gotWr , tt.wantWr) {
				t.Errorf("TranslateTemplate() = %v, want %v", gotWr, tt.wantWr)
			}
		})
	}
}
