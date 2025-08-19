package apihandler

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vrv501/simple-api/internal/db"
	"go.uber.org/mock/gomock"
)

func TestNewAPIHandler(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *APIHandler
	}{
		{
			name: "NewAPIHandler",
			args: args{ctx: context.Background()},
			want: &APIHandler{
				dbClient: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPIHandler(tt.args.ctx); !cmp.Equal(got, tt.want,
				cmpopts.IgnoreUnexported(APIHandler{})) {
				t.Errorf("NewAPIHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIHandler_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBClient := db.NewMockHandler(ctrl)
	type fields struct {
		dbClient db.Handler
	}
	tests := []struct {
		name    string
		fields  fields
		prepare func()
	}{
		{
			name: "Close APIHandler",
			fields: fields{
				dbClient: mockDBClient,
			},
			prepare: func() {
				mockDBClient.EXPECT().Close(gomock.Any()).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}

			a := &APIHandler{
				dbClient: tt.fields.dbClient,
			}
			a.Close()
		})
	}
}
