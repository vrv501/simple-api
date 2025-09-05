package apihandler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	dbErr "github.com/vrv501/simple-api/internal/db/errors"
	"github.com/vrv501/simple-api/internal/generated/mockdb"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

func TestAPIHandler_GetImageById(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBClient := mockdb.NewMockHandler(ctrl)

	tests := []struct {
		name    string
		request genRouter.GetImageByIdRequestObject
		prepare func()
		want    genRouter.GetImageByIdResponseObject
	}{
		{
			name: "image not found",
			prepare: func() {
				mockDBClient.EXPECT().GetPetImage(gomock.Any(), gomock.Any()).
					Return(nil, int64(0), dbErr.ErrNotFound)
			},
			want: genRouter.GetImageByIddefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "image not found",
				},
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "invalid imageid",
			prepare: func() {
				mockDBClient.EXPECT().GetPetImage(gomock.Any(), gomock.Any()).
					Return(nil, int64(0), dbErr.ErrInvalidValue)
			},
			want: genRouter.GetImageByIddefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "invalid image ID",
				},
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "internal error",
			prepare: func() {
				mockDBClient.EXPECT().GetPetImage(gomock.Any(), gomock.Any()).
					Return(nil, int64(0), errors.New(""))
			},
			want: genRouter.GetImageByIddefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "success",
			prepare: func() {
				mockDBClient.EXPECT().GetPetImage(gomock.Any(), gomock.Any()).
					Return(nil, int64(0), nil)
			},
			want: genRouter.GetImageById200ImagejpegResponse{
				Body:          nil,
				ContentLength: 0,
			},
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
			got, _ := a.GetImageById(context.Background(), tt.request)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetImageById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateImage(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Reader
		want    []byte
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := validateImage(tt.r)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("validateImage() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("validateImage() succeeded unexpectedly")
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("validateImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
