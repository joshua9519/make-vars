package main

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

func Test_findTfFiles(t *testing.T) {
	type args struct {
		root string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				root: "test",
			},
			want: []string{
				"test/dir1/file2.tf",
				"test/file1.tf",
				"test/main.tf",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findTfFiles(tt.args.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("findTfFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findTfFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeDups(t *testing.T) {
	type args struct {
		vars []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Success",
			args: args{
				vars: []string{
					"var1",
					"var2",
					"var3",
					"var1",
					"var2",
					"var4",
				},
			},
			want: []string{
				"var1",
				"var2",
				"var3",
				"var4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDups(tt.args.vars); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeDups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readTF(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				file: "test/dir1/file2.tf",
			},
			want: []string{
				"var1",
				"var2",
				"var3",
			},
			wantErr: false,
		},
		{
			name: "Failure - no file",
			args: args{
				file: "test/dir1/blah.tf",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Success - no matches",
			args: args{
				file: "test/file1.tf",
			},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readTF(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("readTF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readTF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeVarsFile(t *testing.T) {
	type args struct {
		vars []string
	}
	tests := []struct {
		name    string
		args    args
		wantF   string
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				vars: []string{
					"var1",
					"var2",
					"var3",
				},
			},
			wantF: `variable "var1" {}
variable "var2" {}
variable "var3" {}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &bytes.Buffer{}
			if err := makeVarsFile(tt.args.vars, f); (err != nil) != tt.wantErr {
				t.Errorf("makeVarsFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotF := f.String(); gotF != tt.wantF {
				t.Errorf("makeVarsFile() = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}

func Test_main(t *testing.T) {
	t.Run("Test main", func(t *testing.T) {
		originalGetwd := getwd
		getwd = func() (string, error) {
			return "test", nil
		}

		defer func() {
			getwd = originalGetwd
		}()

		main()
		os.Remove("test/vars.tf")
	})
}
