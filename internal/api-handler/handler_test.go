package apihandler

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/mock/gomock"

	"github.com/vrv501/simple-api/internal/db"
	"github.com/vrv501/simple-api/internal/generated/mockdb"
)

func TestNewAPIHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want *APIHandler
	}{
		{
			name: "NewAPIHandler",
			want: &APIHandler{
				dbClient: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPIHandler(context.Background()); !cmp.Equal(got, tt.want,
				cmpopts.IgnoreUnexported(APIHandler{})) {
				t.Errorf("NewAPIHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIHandler_Close(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBClient := mockdb.NewMockHandler(ctrl)
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
