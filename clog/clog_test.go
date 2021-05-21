package clog

import (
	"bou.ke/monkey"
	"context"
	"errors"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	type args struct {
		d *SetDefault
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "len(d.Host) == 0",
			args: args{
				d: &SetDefault{},
			},
			wantErr: true,
		},
		{
			name: "setupLogrus error",
			args: args{
				d: &SetDefault{
					Host:   "127.0.0.1",
					LogLvl: "???",
				},
			},
			wantErr: true,
		},
		{
			name: "all is ok",
			args: args{
				d: &SetDefault{
					Host:   "127.0.0.1",
					LogLvl: "error",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Init(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("Fire() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_appLogDocModel_indexName(t *testing.T) {
	tests := []struct {
		name string
		m    appLogDocModel
		want string
	}{
		{
			name: "all is ok",
			m:    nil,
			want: "-1617849131",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti := time.Time{}
			monkey.PatchInstanceMethod(reflect.TypeOf(ti), "Format", func(_ time.Time, layout string) string {
				return "1617849131"
			})
			if got := tt.m.indexName(); got != tt.want {
				t.Errorf("indexName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_esHook_Fire(t *testing.T) {
	type fields struct {
		cmd    string
		client *elastic.Client
	}
	type args struct {
		entry *logrus.Entry
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "all is ok",
			fields: fields{
				cmd:    "",
				client: &elastic.Client{},
			},
			args: args{
				entry: &logrus.Entry{
					Data: map[string]interface{}{
						"test": "test",
					},
					Caller: &runtime.Frame{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := &esHook{
				cmd:    tt.fields.cmd,
				client: tt.fields.client,
			}
			if err := hook.Fire(tt.args.entry); (err != nil) != tt.wantErr {
				t.Errorf("Fire() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_esHook_Levels(t *testing.T) {
	type fields struct {
		cmd    string
		client *elastic.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   []logrus.Level
	}{
		{
			name:   "all is ok",
			fields: fields{},
			want: []logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := &esHook{
				cmd:    tt.fields.cmd,
				client: tt.fields.client,
			}
			if got := hook.Levels(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Levels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_esHook_sendEs(t *testing.T) {
	type fields struct {
		cmd    string
		client *elastic.Client
	}
	type args struct {
		doc appLogDocModel
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "panic",
			fields: fields{},
			args:   args{},
		},
		{
			name:    "err",
			fields:  fields{},
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := &esHook{
				cmd:    tt.fields.cmd,
				client: tt.fields.client,
			}
			if tt.wantErr {
				monkey.PatchInstanceMethod(reflect.TypeOf(&elastic.IndexService{}), "Do", func(_ *elastic.IndexService, ctx context.Context) (*elastic.IndexResponse, error) {
					return nil, errors.New("test")
				})
			}
			hook.sendEs(tt.args.doc)
			monkey.UnpatchAll()
		})
	}
}

func Test_newEsHook(t *testing.T) {
	type args struct {
		cc cfg
	}
	tests := []struct {
		name    string
		args    args
		want    *esHook
		wantErr bool
	}{
		{
			name: "NewClient error",
			args: args{
				cc: cfg{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "all is ok",
			args: args{
				cc: cfg{},
			},
			want:    &esHook{client: nil, cmd: strings.Join(os.Args, " ")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				monkey.Patch(elastic.NewClient, mockNewClinetErr)
			} else {
				monkey.Patch(elastic.NewClient, mockNewClinetOk)
			}
			got, err := newEsHook(tt.args.cc)
			monkey.UnpatchAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("Fire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockNewClinetErr(options ...elastic.ClientOptionFunc) (*elastic.Client, error) {
	return nil, errors.New("test")
}

func mockNewClinetOk(options ...elastic.ClientOptionFunc) (*elastic.Client, error) {
	return nil, nil
}

func Test_newEsLog(t *testing.T) {
	type args struct {
		e *logrus.Entry
	}
	tests := []struct {
		name string
		args args
		want appLogDocModel
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEsLog(tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEsLog() = %v, want %v", got, tt.want)
			}
		})
	}
}
