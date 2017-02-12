package main

import (
	"reflect"
	"testing"
)

const testGolangSrc = `package somepkg

import "fmt"

type T int

func (r T) String() string {
	return fmt.Sprintf("%v", r)
}

func fn() string {
	return "Hello"
}
`

func Test_getGolangFuncs(t *testing.T) {
	tests := []struct {
		name        string
		fileContent []byte
		wantResult  []Func
		wantErr     bool
	}{
		{
			name:        "without error",
			fileContent: []byte(testGolangSrc),
			wantResult: []Func{
				{Name: "String", Begin: 44, End: 103},
				{Name: "fn", Begin: 105, End: 141},
			},
			wantErr: false,
		},
		{
			name:        "with error",
			fileContent: []byte("..."),
			wantResult:  nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := getGolangFuncs(tt.fileContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("getGolangFuncs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("getGolangFuncs() = %#v, want %#v", gotResult, tt.wantResult)
			}
		})
	}
}
