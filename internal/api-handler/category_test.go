package apihandler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	"github.com/vrv501/simple-api/internal/db"
	dbErr "github.com/vrv501/simple-api/internal/db/errors"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

func TestAPIHandler_FindAnimalCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBClient := db.NewMockHandler(ctrl)

	type args struct {
		ctx     context.Context
		request genRouter.FindAnimalCategoryRequestObject
	}
	tests := []struct {
		name    string
		args    args
		prepare func()
		want    genRouter.FindAnimalCategoryResponseObject
		wantErr bool
	}{
		{
			name: "AnimalCategory not found",
			args: args{
				ctx: context.Background(),
				request: genRouter.FindAnimalCategoryRequestObject{
					Params: genRouter.FindAnimalCategoryParams{
						Name: "Dog",
					},
				},
			},
			prepare: func() {
				mockDBClient.EXPECT().FindAnimalCategory(gomock.Any(), "Dog").
					Return(nil, dbErr.ErrNotFound)
			},
			want: genRouter.FindAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "Animal category Dog not found",
				},
				StatusCode: http.StatusNotFound,
			},
			wantErr: false,
		},
		{
			name: "AnimalCategory internal error",
			args: args{
				ctx: context.Background(),
				request: genRouter.FindAnimalCategoryRequestObject{
					Params: genRouter.FindAnimalCategoryParams{
						Name: "Dog",
					},
				},
			},
			prepare: func() {
				mockDBClient.EXPECT().FindAnimalCategory(gomock.Any(), "Dog").
					Return(nil, errors.New(""))
			},
			want: genRouter.FindAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
			wantErr: false,
		},
		{
			name: "AnimalCategory success",
			args: args{
				ctx: context.Background(),
				request: genRouter.FindAnimalCategoryRequestObject{
					Params: genRouter.FindAnimalCategoryParams{
						Name: "Dog",
					},
				},
			},
			prepare: func() {
				mockDBClient.EXPECT().FindAnimalCategory(gomock.Any(), "Dog").
					Return(&genRouter.AnimalCategoryJSONResponse{
						Id:   "1",
						Name: "Dog",
					}, nil)
			},
			want: genRouter.FindAnimalCategory200JSONResponse{
				AnimalCategoryJSONResponse: genRouter.AnimalCategoryJSONResponse{
					Id:   "1",
					Name: "Dog",
				}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				dbClient: mockDBClient,
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			got, err := a.FindAnimalCategory(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIHandler.FindAnimalCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("APIHandler.FindAnimalCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIHandler_AddAnimalCategory(t *testing.T) {
	type fields struct {
		dbClient db.Handler
	}
	type args struct {
		ctx     context.Context
		request genRouter.AddAnimalCategoryRequestObject
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    genRouter.AddAnimalCategoryResponseObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				dbClient: tt.fields.dbClient,
			}
			got, err := a.AddAnimalCategory(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIHandler.AddAnimalCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("APIHandler.AddAnimalCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIHandler_DeleteAnimalCategory(t *testing.T) {
	type fields struct {
		dbClient db.Handler
	}
	type args struct {
		ctx     context.Context
		request genRouter.DeleteAnimalCategoryRequestObject
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    genRouter.DeleteAnimalCategoryResponseObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				dbClient: tt.fields.dbClient,
			}
			got, err := a.DeleteAnimalCategory(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIHandler.DeleteAnimalCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("APIHandler.DeleteAnimalCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIHandler_ReplaceAnimalCategory(t *testing.T) {
	type fields struct {
		dbClient db.Handler
	}
	type args struct {
		ctx     context.Context
		request genRouter.ReplaceAnimalCategoryRequestObject
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    genRouter.ReplaceAnimalCategoryResponseObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				dbClient: tt.fields.dbClient,
			}
			got, err := a.ReplaceAnimalCategory(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIHandler.ReplaceAnimalCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("APIHandler.ReplaceAnimalCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}
