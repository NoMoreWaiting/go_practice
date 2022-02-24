package gormdemo

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Test_operationSqlite3(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"case1"},
		{"case2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operationSqlite3()
		})
	}
}
