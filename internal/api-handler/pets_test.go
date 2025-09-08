package apihandler

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	"github.com/vrv501/simple-api/internal/constants"
	contextKeys "github.com/vrv501/simple-api/internal/context-keys"
	dbErr "github.com/vrv501/simple-api/internal/db/errors"
	"github.com/vrv501/simple-api/internal/generated/mockdb"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

// createTestJPEG creates a simple JPEG image with the specified dimensions for testing
func createTestJPEG(t *testing.T, width, height int) []byte {
	t.Helper()

	// Create a simple image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with a simple pattern
	for y := range height {
		for x := range width {
			// Create a simple checkerboard pattern
			if (x/32+y/32)%2 == 0 {
				img.Set(x, y, color.RGBA{255, 0, 0, 255}) // Red
			} else {
				img.Set(x, y, color.RGBA{0, 255, 0, 255}) // Green
			}
		}
	}

	// Encode to JPEG
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	if err != nil {
		t.Fatalf("Failed to create test JPEG: %v", err)
	}

	return buf.Bytes()
}

// errorReader is a mock reader that always returns an error
type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

// createMalformedMultipartReader creates a multipart reader that will fail during NextPart
func createMalformedMultipartReader(t *testing.T) *multipart.Reader {
	t.Helper()
	// Create a reader that will cause a bufio error during NextPart
	// Use data that's too short and will cause unexpected EOF during header parsing
	malformedData := "--boundary\r\nContent-"
	return multipart.NewReader(strings.NewReader(malformedData), "boundary")
}

// createMultipartReader creates a multipart reader for testing
func createMultipartReader(t *testing.T, stringFields map[string]string, fileFields map[string][]byte) *multipart.Reader {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	defer writer.Close()

	// Add string fields
	for fieldName, value := range stringFields {
		err := writer.WriteField(fieldName, value)
		if err != nil {
			t.Fatalf("Failed to write field %s: %v", fieldName, err)
		}
	}

	// Add file fields
	for fieldName, data := range fileFields {
		part, err := writer.CreateFormFile(fieldName, "test.jpg")
		if err != nil {
			t.Fatalf("Failed to create form file %s: %v", fieldName, err)
		}
		_, err = part.Write(data)
		if err != nil {
			t.Fatalf("Failed to write file data for %s: %v", fieldName, err)
		}
	}

	return multipart.NewReader(&body, writer.Boundary())
}

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
	t.Parallel()
	// Create a valid small JPEG image (256x256)
	validSmallJPEG := createTestJPEG(t, 256, 256)

	// Create an oversized JPEG image (exceeding MaxImgSize)
	oversizedJPEG := make([]byte, constants.MaxImgSize+10)
	copy(oversizedJPEG, validSmallJPEG)
	for i := len(validSmallJPEG); i < len(oversizedJPEG); i++ {
		oversizedJPEG[i] = 0xFF // Fill with dummy data
	}

	// Create a JPEG with valid header but corrupted image data
	// This will pass DecodeConfig but fail on actual Decode
	corruptedJPEG := make([]byte, len(validSmallJPEG))
	copy(corruptedJPEG, validSmallJPEG)

	// Find the Start of Scan (SOS) marker (0xFF 0xDA) and corrupt the data after it
	// The SOS marker indicates the start of the actual image data
	for i := 0; i < len(corruptedJPEG)-1; i++ {
		if corruptedJPEG[i] == 0xFF && corruptedJPEG[i+1] == 0xDA {
			// Found SOS marker, corrupt the image data that follows
			// Skip the SOS marker and its length field, then corrupt the actual data
			start := i + 12 // Skip SOS marker + typical header
			if start < len(corruptedJPEG) {
				// Corrupt middle section of image data while
				// preserving EOI marker at end
				end := len(corruptedJPEG) - 10
				if end > start {
					for j := start; j < end; j++ {
						corruptedJPEG[j] = 0x00 // Corrupt scan data
					}
					break
				}
			}
		}
	}

	tests := []struct {
		name    string
		r       io.Reader
		want    []byte
		wantErr bool
	}{
		{
			name:    "reader error",
			r:       &errorReader{err: errors.New("read error")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "oversized image",
			r:       bytes.NewReader(oversizedJPEG),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid JPEG data",
			r:       bytes.NewReader([]byte("not a jpeg image")),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "image height too large",
			r:       bytes.NewReader(createTestJPEG(t, 257, 1081)),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "corrupted JPEG image data",
			r:       bytes.NewReader(corruptedJPEG),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "valid small JPEG",
			r:       bytes.NewReader(validSmallJPEG),
			want:    validSmallJPEG,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := validateImage(tt.r)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("validateImage() failed: %v", gotErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("validateImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIHandler_AddPet(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBClient := mockdb.NewMockHandler(ctrl)
	ctxU := contextKeys.ContextWithUserID(context.Background(), "1")
	validRequestBody := createMultipartReader(t,
		map[string]string{
			"pet": `{"name":"test"}`,
		},
		map[string][]byte{
			"photos": createTestJPEG(t, 256, 256),
		})

	tests := []struct {
		name    string
		ctx     context.Context
		request genRouter.AddPetRequestObject
		prepare func()
		want    genRouter.AddPetResponseObject
	}{
		{
			name: "userID not in context",
			ctx:  context.Background(),
			want: genRouter.AddPetdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "multipart parsing error",
			ctx:  ctxU,
			request: genRouter.AddPetRequestObject{
				Body: createMalformedMultipartReader(t),
			},
			want: genRouter.AddPetdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgIncorrectReqEncoding,
				},
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "invalid JSON in pet field",
			ctx:  ctxU,
			request: genRouter.AddPetRequestObject{
				Body: createMultipartReader(t, map[string]string{
					"pet": "invalid json",
				}, nil),
			},
			want: genRouter.AddPetdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgIncorrectReqEncoding,
				},
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "invalid image in photos field",
			ctx:  ctxU,
			request: genRouter.AddPetRequestObject{
				Body: createMultipartReader(t, nil, map[string][]byte{
					"photos": []byte("not an image"),
				}),
			},
			want: genRouter.AddPetdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "jpeg image is corrupted",
				},
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "unknown multipart field",
			ctx:  ctxU,
			request: genRouter.AddPetRequestObject{
				Body: createMultipartReader(t, map[string]string{
					"unknown": "field",
				}, nil),
			},
			want: genRouter.AddPetdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "unknown multipart field unknown",
				},
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "internal DB error",
			ctx:  ctxU,
			request: genRouter.AddPetRequestObject{
				Body: validRequestBody,
			},
			prepare: func() {
				mockDBClient.EXPECT().AddPet(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			want: genRouter.AddPetdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "err invalid value",
			ctx:  ctxU,
			request: genRouter.AddPetRequestObject{
				Body: validRequestBody,
			},
			prepare: func() {
				mockDBClient.EXPECT().AddPet(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&dbErr.HintError{Err: dbErr.ErrInvalidValue})
			},
			want: genRouter.AddPetdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "invalid value for ",
				},
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "err not found",
			ctx:  ctxU,
			request: genRouter.AddPetRequestObject{
				Body: validRequestBody,
			},
			prepare: func() {
				mockDBClient.EXPECT().AddPet(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&dbErr.HintError{Err: dbErr.ErrNotFound})
			},
			want: genRouter.AddPetdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: " not found",
				},
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "err conflict",
			ctx:  ctxU,
			request: genRouter.AddPetRequestObject{
				Body: validRequestBody,
			},
			prepare: func() {
				mockDBClient.EXPECT().AddPet(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dbErr.ErrConflict)
			},
			want: genRouter.AddPetdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "Similar pet already exists",
				},
				StatusCode: http.StatusConflict,
			},
		},
		{
			name: "success",
			ctx:  ctxU,
			request: genRouter.AddPetRequestObject{
				Body: validRequestBody,
			},
			prepare: func() {
				mockDBClient.EXPECT().AddPet(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			want: genRouter.AddPet202Response{},
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
			got, _ := a.AddPet(tt.ctx, tt.request)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("AddPet() = %v, want %v", got, tt.want)
			}
		})
	}
}
