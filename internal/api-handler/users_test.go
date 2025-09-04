package apihandler

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	contextKeys "github.com/vrv501/simple-api/internal/context-keys"
	dbErr "github.com/vrv501/simple-api/internal/db/errors"
	"github.com/vrv501/simple-api/internal/generated/mockdb"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

func Test_hashPassword(t *testing.T) {
	t.Parallel()
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test hashing a password",
			args: args{
				password: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("hashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err = comparePasswords(got, tt.args.password); err != nil {
				t.Errorf("comparePasswords() = %v", err)
			}
		})
	}
}

func TestAPIHandler_CreateUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBClient := mockdb.NewMockHandler(ctrl)
	validReqBody := genRouter.CreateUserRequestObject{
		Body: &genRouter.CreateUserJSONRequestBody{
			Password: "test",
		},
	}

	type args struct {
		request genRouter.CreateUserRequestObject
	}
	tests := []struct {
		name     string
		args     args
		prepFunc func()
		want     genRouter.CreateUserResponseObject
	}{
		{
			name: "password too long",
			args: args{
				request: genRouter.CreateUserRequestObject{
					Body: &genRouter.CreateUserJSONRequestBody{
						Password: string(make([]byte, 100)),
					},
				},
			},
			want: genRouter.CreateUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "conflict error",
			args: args{
				request: validReqBody,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().AddUser(gomock.Any(),
					gomock.Any()).Return(nil, &dbErr.HintError{Err: dbErr.ErrConflict})
			},
			want: genRouter.CreateUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: " " + errMsgAlreadyInUse,
				},
				StatusCode: http.StatusConflict,
			},
		},
		{
			name: "internal error",
			args: args{
				request: validReqBody,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().AddUser(gomock.Any(),
					gomock.Any()).Return(nil, errors.New(""))
			},
			want: genRouter.CreateUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "success",
			args: args{
				request: validReqBody,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().AddUser(gomock.Any(),
					gomock.Any()).Return(&genRouter.UserJSONResponse{}, nil)
			},
			want: genRouter.CreateUser201JSONResponse{
				UserJSONResponse: genRouter.UserJSONResponse{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				dbClient: mockDBClient,
			}
			if tt.prepFunc != nil {
				tt.prepFunc()
			}
			got, _ := a.CreateUser(context.Background(), tt.args.request)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("APIHandler.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIHandler_DeleteUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBClient := mockdb.NewMockHandler(ctrl)
	ctxU := contextKeys.ContextWithUserID(context.Background(), "1")

	type args struct {
		ctx context.Context
		in1 genRouter.DeleteUserRequestObject
	}
	tests := []struct {
		name     string
		args     args
		prepFunc func()
		want     genRouter.DeleteUserResponseObject
	}{
		{
			name: "userID not in context",
			args: args{
				ctx: context.Background(),
			},
			want: genRouter.DeleteUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "invalid userID",
			args: args{
				ctx: ctxU,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().DeleteUser(gomock.Any(),
					gomock.Any()).Return(dbErr.ErrInvalidValue)
			},
			want: genRouter.DeleteUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgInvalidUserID,
				},
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "userID not found",
			args: args{
				ctx: ctxU,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().DeleteUser(gomock.Any(),
					gomock.Any()).Return(dbErr.ErrNotFound)
			},
			want: genRouter.DeleteUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgUserNotFound,
				},
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "fKey error",
			args: args{
				ctx: ctxU,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().DeleteUser(gomock.Any(),
					gomock.Any()).Return(&dbErr.HintError{Err: dbErr.ErrForeignKeyViolation})
			},
			want: genRouter.DeleteUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "User cannot be deleted as there are pending ",
				},
				StatusCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "internal error",
			args: args{
				ctx: ctxU,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().DeleteUser(gomock.Any(),
					gomock.Any()).Return(errors.New(""))
			},
			want: genRouter.DeleteUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "success",
			args: args{
				ctx: ctxU,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().DeleteUser(gomock.Any(),
					gomock.Any()).Return(nil)
			},
			want: genRouter.DeleteUser204Response{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				dbClient: mockDBClient,
			}
			if tt.prepFunc != nil {
				tt.prepFunc()
			}
			got, _ := a.DeleteUser(tt.args.ctx, tt.args.in1)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("APIHandler.DeleteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIHandler_GetUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBClient := mockdb.NewMockHandler(ctrl)
	ctxU := contextKeys.ContextWithUserID(context.Background(), "1")

	type args struct {
		ctx context.Context
		in1 genRouter.GetUserRequestObject
	}
	tests := []struct {
		name     string
		args     args
		prepFunc func()
		want     genRouter.GetUserResponseObject
	}{
		{
			name: "userid not in context",
			args: args{
				ctx: context.Background(),
			},
			want: genRouter.GetUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "invalid userid",
			args: args{
				ctx: ctxU,
			},
			want: genRouter.GetUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgInvalidUserID,
				},
				StatusCode: http.StatusBadRequest,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().GetUser(gomock.Any(),
					gomock.Any()).Return(nil, dbErr.ErrInvalidValue)
			},
		},
		{
			name: "userid not found",
			args: args{
				ctx: ctxU,
			},
			want: genRouter.GetUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgUserNotFound,
				},
				StatusCode: http.StatusNotFound,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().GetUser(gomock.Any(),
					gomock.Any()).Return(nil, dbErr.ErrNotFound)
			},
		},
		{
			name: "internal error",
			args: args{
				ctx: ctxU,
			},
			want: genRouter.GetUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().GetUser(gomock.Any(),
					gomock.Any()).Return(nil, errors.New(""))
			},
		},
		{
			name: "success",
			args: args{
				ctx: ctxU,
			},
			want: genRouter.GetUser200JSONResponse{
				UserJSONResponse: genRouter.UserJSONResponse{},
			},
			prepFunc: func() {
				mockDBClient.EXPECT().GetUser(gomock.Any(),
					gomock.Any()).
					Return(&genRouter.UserJSONResponse{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				a := &APIHandler{
					dbClient: mockDBClient,
				}
				if tt.prepFunc != nil {
					tt.prepFunc()
				}
				got, _ := a.GetUser(tt.args.ctx, tt.args.in1)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("APIHandler.GetUser() = %v, want %v", got, tt.want)
				}
			})
	}
}

func TestAPIHandler_PatchUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBClient := mockdb.NewMockHandler(ctrl)
	pswdTest := "test"
	ctxU := contextKeys.ContextWithUserID(context.Background(), "1")
	validReq := genRouter.PatchUserRequestObject{
		Body: &genRouter.PatchUserApplicationMergePatchPlusJSONRequestBody{
			Password: &pswdTest,
		},
	}
	tooLongPswd := string(make([]byte, 100))

	type args struct {
		ctx     context.Context
		request genRouter.PatchUserRequestObject
	}
	tests := []struct {
		name     string
		args     args
		prepFunc func()
		want     genRouter.PatchUserResponseObject
	}{
		{
			name: "userID not in context",
			args: args{
				ctx: context.Background(),
			},
			want: genRouter.PatchUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "nil request body",
			args: args{
				ctx: ctxU,
				request: genRouter.PatchUserRequestObject{
					Body: &genRouter.PatchUserApplicationMergePatchPlusJSONRequestBody{},
				},
			},
			want: genRouter.PatchUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "Nothing to Update",
				},
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "pswd too long",
			args: args{
				ctx: ctxU,
				request: genRouter.PatchUserRequestObject{
					Body: &genRouter.PatchUserApplicationMergePatchPlusJSONRequestBody{
						Password: &tooLongPswd,
					},
				},
			},
			want: genRouter.PatchUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "invalidID",
			args: args{
				ctx:     ctxU,
				request: validReq,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().PatchUser(gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil, dbErr.ErrInvalidValue)
			},
			want: genRouter.PatchUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgInvalidUserID,
				},
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "user ID not found",
			args: args{
				ctx:     ctxU,
				request: validReq,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().PatchUser(gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil, dbErr.ErrNotFound)
			},
			want: genRouter.PatchUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgUserNotFound,
				},
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "conflict error",
			args: args{
				ctx:     ctxU,
				request: validReq,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().PatchUser(gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil, dbErr.ErrConflict)
			},
			want: genRouter.PatchUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "phone_number " + errMsgAlreadyInUse,
				},
				StatusCode: http.StatusConflict,
			},
		},
		{
			name: "internal error",
			args: args{
				ctx:     ctxU,
				request: validReq,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().PatchUser(gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil, errors.New(""))
			},
			want: genRouter.PatchUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "success",
			args: args{
				ctx:     ctxU,
				request: validReq,
			},
			prepFunc: func() {
				mockDBClient.EXPECT().PatchUser(gomock.Any(),
					gomock.Any(), gomock.Any()).Return(&genRouter.UserJSONResponse{}, nil)
			},
			want: genRouter.PatchUser200JSONResponse{
				UserJSONResponse: genRouter.UserJSONResponse{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIHandler{
				dbClient: mockDBClient,
			}
			if tt.prepFunc != nil {
				tt.prepFunc()
			}
			got, _ := a.PatchUser(tt.args.ctx, tt.args.request)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("APIHandler.PatchUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
