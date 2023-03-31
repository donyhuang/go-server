package reflectx

import (
	"errors"
	"reflect"
	"testing"
)

type a struct {
	Name     string
	Done     bool
	Num      int
	NumChan  chan int
	DoneFunc func() interface{}
}

func TestResetStructPointer(t *testing.T) {
	aa := &a{
		Name:    "donyhuang",
		Done:    true,
		Num:     123,
		NumChan: make(chan int, 1),
		DoneFunc: func() interface{} {
			return "test"
		},
	}
	t.Logf("%+v", aa)
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				v: aa,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ResetStructPointer(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("ResetStructPointer() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("%v %+v", reflect.DeepEqual(tt.args.v, &a{}), *(tt.args.v.(*a)))
		})
	}
}

func TestCopyStruct(t *testing.T) {
	type TestType struct {
		Name string
		Age  int
	}
	type s1 struct {
		Name string
		Age  int
	}
	type args struct {
		l interface{}
		r interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "ok",
			args: args{
				l: &TestType{
					Name: "donyhuang",
					Age:  18,
				},
				r: TestType{
					Name: "pzfyp",
					Age:  0,
				},
			},
		},
		{
			name: "not pointer",
			args: args{
				l: TestType{
					Name: "donyhuang",
					Age:  18,
				},
				r: TestType{
					Name: "pzfyp",
					Age:  0,
				},
			},
			wantErr: ErrorNeedPointerType,
		},
		{
			name: "not same type",
			args: args{
				l: &struct {
					Name string
					Age  int
				}{
					Name: "donyhuang",
					Age:  18,
				},
				r: TestType{
					Name: "pzfyp",
					Age:  0,
				},
			},
			wantErr: ErrorTwoTypeNotEqual,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetStructNotEmpty(tt.args.l, tt.args.r)
			if err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("%v", err)
				}
			} else {
				t.Logf("l %+v", *(tt.args.l.(*TestType)))
			}
		})
	}
}
