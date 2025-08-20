package apihandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	"github.com/vrv501/simple-api/internal/db"
	dbErr "github.com/vrv501/simple-api/internal/db/errors"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

func TestAPIHandler_FindAnimalCategory(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBClient := db.NewMockHandler(ctrl)

	type args struct {
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
			got, err := a.FindAnimalCategory(context.Background(), tt.args.request)
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
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBClient := db.NewMockHandler(ctrl)

	type args struct {
		request genRouter.AddAnimalCategoryRequestObject
	}
	tests := []struct {
		name    string
		args    args
		prepare func()
		want    genRouter.AddAnimalCategoryResponseObject
		wantErr bool
	}{
		{
			name: "AddAnimalCategory conflict",
			args: args{
				request: genRouter.AddAnimalCategoryRequestObject{
					Body: &genRouter.AddAnimalCategoryJSONRequestBody{
						Name: "Dog",
					},
				},
			},
			want: genRouter.AddAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: fmt.Sprintf(errMsgAnimalCategoryExists, "Dog"),
				},
				StatusCode: http.StatusConflict,
			},
			prepare: func() {
				mockDBClient.EXPECT().AddAnimalCategory(gomock.Any(), "Dog").
					Return(nil, &dbErr.ConflictError{})
			},
			wantErr: false,
		},
		{
			name: "AddAnimalCategory internal error",
			args: args{
				request: genRouter.AddAnimalCategoryRequestObject{
					Body: &genRouter.AddAnimalCategoryJSONRequestBody{
						Name: "Dog",
					},
				},
			},
			want: genRouter.AddAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
			prepare: func() {
				mockDBClient.EXPECT().AddAnimalCategory(gomock.Any(), "Dog").
					Return(nil, errors.New(""))
			},
			wantErr: false,
		},
		{
			name: "AddAnimalCategory success",
			args: args{
				request: genRouter.AddAnimalCategoryRequestObject{
					Body: &genRouter.AddAnimalCategoryJSONRequestBody{
						Name: "Dog",
					},
				},
			},
			want: genRouter.AddAnimalCategory201JSONResponse{
				AnimalCategoryJSONResponse: genRouter.AnimalCategoryJSONResponse{
					Id:   "1",
					Name: "Dog",
				},
			},
			prepare: func() {
				mockDBClient.EXPECT().AddAnimalCategory(gomock.Any(), "Dog").
					Return(&genRouter.AnimalCategoryJSONResponse{
						Id:   "1",
						Name: "Dog",
					}, nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			a := &APIHandler{
				dbClient: mockDBClient,
			}
			got, err := a.AddAnimalCategory(context.Background(), tt.args.request)
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
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBClient := db.NewMockHandler(ctrl)

	type args struct {
		request genRouter.DeleteAnimalCategoryRequestObject
	}
	tests := []struct {
		name    string
		args    args
		prepare func()
		want    genRouter.DeleteAnimalCategoryResponseObject
		wantErr bool
	}{
		{
			name: "DeleteAnimalCategory invalid ID",
			args: args{
				request: genRouter.DeleteAnimalCategoryRequestObject{
					AnimalCategoryId: "invalid-id",
				},
			},
			want: genRouter.DeleteAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgInvalidAnimalCategoryID,
				},
				StatusCode: http.StatusBadRequest,
			},
			prepare: func() {
				mockDBClient.EXPECT().DeleteAnimalCategory(gomock.Any(), "invalid-id").
					Return(dbErr.ErrInvalidID)
			},
			wantErr: false,
		},
		{
			name: "DeleteAnimalCategory not found",
			args: args{
				request: genRouter.DeleteAnimalCategoryRequestObject{
					AnimalCategoryId: "1",
				},
			},
			want: genRouter.DeleteAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgAnimalCategoryNotFound + " 1",
				},
				StatusCode: http.StatusNotFound,
			},
			prepare: func() {
				mockDBClient.EXPECT().DeleteAnimalCategory(gomock.Any(), "1").
					Return(dbErr.ErrNotFound)
			},
			wantErr: false,
		},
		{
			name: "DeleteAnimalCategory foreign constraint",
			args: args{
				request: genRouter.DeleteAnimalCategoryRequestObject{
					AnimalCategoryId: "1",
				},
			},
			want: genRouter.DeleteAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "Pets found for animal category 1",
				},
				StatusCode: http.StatusUnprocessableEntity,
			},
			prepare: func() {
				mockDBClient.EXPECT().DeleteAnimalCategory(gomock.Any(), "1").
					Return(dbErr.ErrForeignKeyConstraint)
			},
			wantErr: false,
		},
		{
			name: "DeleteAnimalCategory internal error",
			args: args{
				request: genRouter.DeleteAnimalCategoryRequestObject{
					AnimalCategoryId: "1",
				},
			},
			want: genRouter.DeleteAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
			prepare: func() {
				mockDBClient.EXPECT().DeleteAnimalCategory(gomock.Any(), "1").
					Return(errors.New(""))
			},
			wantErr: false,
		},
		{
			name: "DeleteAnimalCategory success",
			args: args{
				request: genRouter.DeleteAnimalCategoryRequestObject{
					AnimalCategoryId: "1",
				},
			},
			want: genRouter.DeleteAnimalCategory204Response{},
			prepare: func() {
				mockDBClient.EXPECT().DeleteAnimalCategory(gomock.Any(), "1").
					Return(nil)
			},
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
			got, err := a.DeleteAnimalCategory(context.Background(), tt.args.request)
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
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBClient := db.NewMockHandler(ctrl)

	type args struct {
		request genRouter.ReplaceAnimalCategoryRequestObject
	}
	tests := []struct {
		name    string
		args    args
		prepare func()
		want    genRouter.ReplaceAnimalCategoryResponseObject
		wantErr bool
	}{
		{
			name: "ReplaceAnimalCategory invalid ID",
			args: args{
				request: genRouter.ReplaceAnimalCategoryRequestObject{
					AnimalCategoryId: "invalid-id",
					Body: &genRouter.ReplaceAnimalCategoryJSONRequestBody{
						Name: "Dog",
					},
				},
			},
			want: genRouter.ReplaceAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgInvalidAnimalCategoryID,
				},
				StatusCode: http.StatusBadRequest,
			},
			prepare: func() {
				mockDBClient.EXPECT().UpdateAnimalCategory(gomock.Any(), "invalid-id", "Dog").
					Return(nil, dbErr.ErrInvalidID)
			},
			wantErr: false,
		},
		{
			name: "ReplaceAnimalCategory not found",
			args: args{
				request: genRouter.ReplaceAnimalCategoryRequestObject{
					AnimalCategoryId: "1",
					Body: &genRouter.ReplaceAnimalCategoryJSONRequestBody{
						Name: "Dog",
					},
				},
			},
			want: genRouter.ReplaceAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgAnimalCategoryNotFound + " 1",
				},
				StatusCode: http.StatusNotFound,
			},
			prepare: func() {
				mockDBClient.EXPECT().UpdateAnimalCategory(gomock.Any(), "1", "Dog").
					Return(nil, dbErr.ErrNotFound)
			},
			wantErr: false,
		},
		{
			name: "ReplaceAnimalCategory error conflict",
			args: args{
				request: genRouter.ReplaceAnimalCategoryRequestObject{
					AnimalCategoryId: "1",
					Body: &genRouter.ReplaceAnimalCategoryJSONRequestBody{
						Name: "Dog",
					},
				},
			},
			want: genRouter.ReplaceAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: fmt.Sprintf(errMsgAnimalCategoryExists, "Dog"),
				},
				StatusCode: http.StatusUnprocessableEntity,
			},
			prepare: func() {
				mockDBClient.EXPECT().UpdateAnimalCategory(gomock.Any(), "1", "Dog").
					Return(nil, &dbErr.ConflictError{})
			},
			wantErr: false,
		},
		{
			name: "ReplaceAnimalCategory internal error",
			args: args{
				request: genRouter.ReplaceAnimalCategoryRequestObject{
					AnimalCategoryId: "1",
					Body: &genRouter.ReplaceAnimalCategoryJSONRequestBody{
						Name: "Dog",
					},
				},
			},
			want: genRouter.ReplaceAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
			prepare: func() {
				mockDBClient.EXPECT().UpdateAnimalCategory(gomock.Any(), "1", "Dog").
					Return(nil, errors.New(""))
			},
			wantErr: false,
		},
		{
			name: "ReplaceAnimalCategory success",
			args: args{
				request: genRouter.ReplaceAnimalCategoryRequestObject{
					AnimalCategoryId: "1",
					Body: &genRouter.ReplaceAnimalCategoryJSONRequestBody{
						Name: "Dog",
					},
				},
			},
			want: genRouter.ReplaceAnimalCategory200JSONResponse{
				AnimalCategoryJSONResponse: genRouter.AnimalCategoryJSONResponse{
					Id:   "1",
					Name: "Dog",
				},
			},
			prepare: func() {
				mockDBClient.EXPECT().UpdateAnimalCategory(gomock.Any(), "1", "Dog").
					Return(&genRouter.AnimalCategoryJSONResponse{
						Id:   "1",
						Name: "Dog",
					}, nil)
			},
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
			got, err := a.ReplaceAnimalCategory(context.Background(), tt.args.request)
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
