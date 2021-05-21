package utils

import (
	"archive/zip"
	"bou.ke/monkey"
	"errors"
	"io"
	"reflect"
	"sync"
	"testing"
)

func TestZipFilesWithBytes(t *testing.T) {
	type args struct {
		needZipPkgName string
		file           *File
	}
	tests := []struct {
		name             string
		args             args
		want             []byte
		wantErr          bool
		wantZipCreateErr bool
		wantZipWriteErr  bool
		wantZipCloseErr  bool
	}{
		{
			name: "zipWriter.Create error",
			args: args{
				needZipPkgName: "",
				file: &File{
					Mu: &sync.Mutex{},
					FileInfos: []*FileInfo{
						&FileInfo{
							Filename: "test",
							FileByte: []byte("test"),
						},
					},
				},
			},
			want:             nil,
			wantErr:          true,
			wantZipCreateErr: true,
		},
		{
			name: "zipWriter.Write error",
			args: args{
				needZipPkgName: "",
				file: &File{
					Mu: &sync.Mutex{},
					FileInfos: []*FileInfo{
						&FileInfo{
							Filename: "test",
							FileByte: []byte("test"),
						},
					},
				},
			},
			want:            nil,
			wantErr:         true,
			wantZipWriteErr: true,
		},
		{
			name: "zipWriter.Close error",
			args: args{
				needZipPkgName: "",
				file: &File{
					Mu: &sync.Mutex{},
					FileInfos: []*FileInfo{
						&FileInfo{
							Filename: "test",
							FileByte: []byte("test"),
						},
					},
				},
			},
			want:            nil,
			wantErr:         true,
			wantZipCloseErr: true,
		},
		{
			name: "all is ok",
			args: args{
				needZipPkgName: "",
				file: &File{
					Mu: &sync.Mutex{},
					FileInfos: []*FileInfo{
						&FileInfo{
							Filename: "test",
							FileByte: []byte("test"),
						},
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&zip.Writer{}), "Create", func(_ *zip.Writer, name string) (io.Writer, error) {
				if tt.wantZipCreateErr {
					return nil, errors.New("test")
				}
				if tt.wantZipWriteErr {
					return &MockIoWriterErr{}, nil
				}
				return &MockIoWriter{}, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&zip.Writer{}), "Close", func(_ *zip.Writer) error {
				if tt.wantZipCloseErr {
					return errors.New("test")
				}
				return nil
			})
			got, err := ZipFilesWithBytes(tt.args.needZipPkgName, tt.args.file)
			monkey.UnpatchAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("ZipFilesWithBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZipFilesWithBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type MockIoWriter struct {
}
type MockIoWriterErr struct {
}

func (m *MockIoWriter) Write(p []byte) (int, error) {
	return 1, nil
}

func (m *MockIoWriterErr) Write(p []byte) (int, error) {
	return 0, errors.New("test")
}
