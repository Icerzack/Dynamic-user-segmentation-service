package app

import (
	"avito-backend-internship/internal/pkg/db"
	"avito-backend-internship/internal/pkg/model"
	"avito-backend-internship/internal/pkg/service"
	mock_service "avito-backend-internship/internal/pkg/service/mocks"
	"context"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestServer_mainHandler(t *testing.T) {
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name           string
		args           args
		expectedOutput string
		expectedCode   int
	}{
		{
			name: "Test mainPageHandler return 'it's good' and 200 status code",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			expectedOutput: "it's good",
			expectedCode:   http.StatusOK,
		},
		{
			name: "Test mainPageHandler return '405 Method not allowed' and 405 status code",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodPost, "/", nil),
			},
			expectedOutput: "405 Method not allowed",
			expectedCode:   http.StatusMethodNotAllowed,
		},
		{
			name: "Test mainPageHandler return '404 Not found' and 404 status code",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodGet, "/aaa", nil),
			},
			expectedOutput: "404 Not found",
			expectedCode:   http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database := db.NewDBStub()
			service := service.NewServiceStub()
			s := NewServer(
				context.Background(),
				&Config{Port: ""},
				database,
				service,
			)

			s.mainHandler(tt.args.writer, tt.args.request)
			res := tt.args.writer.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if string(data) != tt.expectedOutput {
				t.Errorf("Output -- WANT: %v GOT: %v", tt.expectedOutput, string(data))
			}
			if res.StatusCode != tt.expectedCode {
				t.Errorf("StatusCode -- WANT: %v GOT: %v", tt.expectedCode, strconv.Itoa(res.StatusCode))
			}
		})
	}
}

func TestServer_createSegmentHandler(t *testing.T) {
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	type fields struct {
		service *mock_service.MockService
	}
	tests := []struct {
		name           string
		prepare        func(f *fields)
		args           args
		expectedOutput string
		expectedCode   int
	}{
		{
			name: "Test createSegmentHandler perform insert in database and return 'Success' and 200 status code",
			prepare: func(f *fields) {
				f.service.EXPECT().InsertSegmentIntoDatabase(context.Background(), gomock.Any(), gomock.Any()).Return(nil)
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodPost, "/create-segment", strings.NewReader(`{"title": "abc", "description": "-"}`)),
			},
			expectedOutput: `{"status":"Success"}`,
			expectedCode:   http.StatusOK,
		},
		{
			name: "Test createSegmentHandler doesn't perform insert in database and return '405 Method not allowed' and 405 status code",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodGet, "/create-segment", strings.NewReader(`{"title": "abc", "description": "-"}`)),
			},
			expectedOutput: `405 Method not allowed`,
			expectedCode:   http.StatusMethodNotAllowed,
		},
		{
			name: "Test createSegmentHandler doesn't perform insert in database and return 'Invalid JSON' and 400 status code",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodPost, "/create-segment", strings.NewReader(`{"name": "abc", "description": "-"}`)),
			},
			expectedOutput: `{"status":"Invalid JSON"}`,
			expectedCode:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				service: mock_service.NewMockService(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			s := NewServer(
				context.Background(),
				&Config{Port: ""},
				db.NewDBStub(),
				f.service,
			)

			s.createSegmentHandler(tt.args.writer, tt.args.request)
			res := tt.args.writer.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if string(data) != tt.expectedOutput {
				t.Errorf("Output -- WANT: %v GOT: %v", tt.expectedOutput, string(data))
			}
			if res.StatusCode != tt.expectedCode {
				t.Errorf("StatusCode -- WANT: %v GOT: %v", tt.expectedCode, strconv.Itoa(res.StatusCode))
			}
		})
	}
}

func TestServer_deleteSegmentHandler(t *testing.T) {
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	type fields struct {
		service *mock_service.MockService
	}
	tests := []struct {
		name           string
		prepare        func(f *fields)
		args           args
		expectedOutput string
		expectedCode   int
	}{
		{
			name: "Test deleteSegmentHandler perform delete in database and return 'Success' and 200 status code",
			prepare: func(f *fields) {
				f.service.EXPECT().DeleteSegmentFromDatabase(context.Background(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodDelete, "/delete-segment", strings.NewReader(`{"title": "abc"}`)),
			},
			expectedOutput: `{"status":"Success"}`,
			expectedCode:   http.StatusOK,
		},
		{
			name: "Test deleteSegmentHandler perform delete in database and return 'Failed: no such element to delete' and 200 status code",
			prepare: func(f *fields) {
				f.service.EXPECT().DeleteSegmentFromDatabase(context.Background(), gomock.Any(), gomock.Any()).Return(false, nil)
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodDelete, "/delete-segment", strings.NewReader(`{"title": "abc"}`)),
			},
			expectedOutput: `{"status":"Failed: no such element to delete"}`,
			expectedCode:   http.StatusOK,
		},
		{
			name: "Test deleteSegmentHandler doesn't perform delete in database and return '405 Method not allowed' and 405 status code",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodGet, "/delete-segment", strings.NewReader(`{"title": "abc"}`)),
			},
			expectedOutput: `405 Method not allowed`,
			expectedCode:   http.StatusMethodNotAllowed,
		},
		{
			name: "Test deleteSegmentHandler doesn't perform delete in database and return 'Invalid JSON' and 400 status code",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodDelete, "/delete-segment", strings.NewReader(`{"name": "abc"}`)),
			},
			expectedOutput: `{"status":"Invalid JSON"}`,
			expectedCode:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				service: mock_service.NewMockService(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			s := NewServer(
				context.Background(),
				&Config{Port: ""},
				db.NewDBStub(),
				f.service,
			)

			s.deleteSegmentHandler(tt.args.writer, tt.args.request)
			res := tt.args.writer.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if string(data) != tt.expectedOutput {
				t.Errorf("Output -- WANT: %v GOT: %v", tt.expectedOutput, string(data))
			}
			if res.StatusCode != tt.expectedCode {
				t.Errorf("StatusCode -- WANT: %v GOT: %v", tt.expectedCode, strconv.Itoa(res.StatusCode))
			}
		})
	}
}

func TestServer_userInSegmentHandler(t *testing.T) {
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	type fields struct {
		service *mock_service.MockService
	}
	tests := []struct {
		name           string
		prepare        func(f *fields)
		args           args
		expectedOutput string
		expectedCode   int
	}{
		{
			name: "Test userInSegmentHandler perform modify in database and return 'Success' and 200 status code",
			prepare: func(f *fields) {
				f.service.EXPECT().ModifyUsersSegmentsInDatabase(context.Background(), gomock.Any(), gomock.Any()).Return([]string{}, []string{}, nil)
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodPut, "/user-in-segment", strings.NewReader(`{"user_id": 111, "seg_titles_add": ["a1"], "seg_titles_delete": ["a1"]}`)),
			},
			expectedOutput: `{"status":"Success"}`,
			expectedCode:   http.StatusOK,
		},
		{
			name: "Test userInSegmentHandler perform modify in database and return 'Success and seg_titles_not_exist_add:[a1]' and 200 status code",
			prepare: func(f *fields) {
				f.service.EXPECT().ModifyUsersSegmentsInDatabase(context.Background(), gomock.Any(), gomock.Any()).Return([]string{"a1"}, []string{}, nil)
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodPut, "/user-in-segment", strings.NewReader(`{"user_id": 111, "seg_titles_add": ["a1"], "seg_titles_delete": ["a1"]}`)),
			},
			expectedOutput: `{"status":"Success","seg_titles_not_exist_add":["a1"]}`,
			expectedCode:   http.StatusOK,
		},
		{
			name: "Test userInSegmentHandler perform modify in database and return 'Success and seg_titles_not_exist_delete:[a1]' and 200 status code",
			prepare: func(f *fields) {
				f.service.EXPECT().ModifyUsersSegmentsInDatabase(context.Background(), gomock.Any(), gomock.Any()).Return([]string{}, []string{"a1"}, nil)
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodPut, "/user-in-segment", strings.NewReader(`{"user_id": 111, "seg_titles_add": ["a1"], "seg_titles_delete": ["a1"]}`)),
			},
			expectedOutput: `{"status":"Success","seg_titles_not_exist_delete":["a1"]}`,
			expectedCode:   http.StatusOK,
		},
		{
			name: "Test userInSegmentHandler doesn't perform modify in database and return 'Invalid JSON' and 400 status code",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodPut, "/user-in-segment", strings.NewReader(`{"name": "abc"}`)),
			},
			expectedOutput: `{"status":"Invalid JSON"}`,
			expectedCode:   http.StatusBadRequest,
		},
		{
			name: "Test userInSegmentHandler doesn't perform modify in database and return '405 Method not allowed' and 405 status code",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodGet, "/user-in-segment", strings.NewReader(`{"user_id": "111"}`)),
			},
			expectedOutput: `405 Method not allowed`,
			expectedCode:   http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				service: mock_service.NewMockService(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			s := NewServer(
				context.Background(),
				&Config{Port: ""},
				db.NewDBStub(),
				f.service,
			)

			s.userInSegmentHandler(tt.args.writer, tt.args.request)
			res := tt.args.writer.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if string(data) != tt.expectedOutput {
				t.Errorf("Output -- WANT: %v GOT: %v", tt.expectedOutput, string(data))
			}
			if res.StatusCode != tt.expectedCode {
				t.Errorf("StatusCode -- WANT: %v GOT: %v", tt.expectedCode, strconv.Itoa(res.StatusCode))
			}
		})
	}
}

func TestServer_getUserSegmentsHandler(t *testing.T) {
	type args struct {
		writer  *httptest.ResponseRecorder
		request *http.Request
	}
	type fields struct {
		service *mock_service.MockService
	}
	tests := []struct {
		name           string
		prepare        func(f *fields)
		args           args
		expectedOutput string
		expectedCode   int
	}{
		{
			name: "Test getUserSegmentsHandler perform read from database and return 'Success and segments' and 200 status code",
			prepare: func(f *fields) {
				f.service.EXPECT().GetUserSegmentsFromDatabase(context.Background(), gomock.Any(), gomock.Any()).Return([]model.UserSegments{
					{
						Title:       "a1",
						Description: "-",
					},
				}, nil)
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodGet, "/get-user-segments", strings.NewReader(`{"user_id": 111}`)),
			},
			expectedOutput: `{"status":"Success","segments":[{"title":"a1","description":"-"}]}`,
			expectedCode:   http.StatusOK,
		},
		{
			name: "Test getUserSegmentsHandler perform read from database and return 'Failure: this user is not in any segment' and 200 status code",
			prepare: func(f *fields) {
				f.service.EXPECT().GetUserSegmentsFromDatabase(context.Background(), gomock.Any(), gomock.Any()).Return([]model.UserSegments{}, nil)
			},
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodGet, "/get-user-segments", strings.NewReader(`{"user_id": 111}`)),
			},
			expectedOutput: `{"status":"Failure: this user is not in any segment"}`,
			expectedCode:   http.StatusOK,
		},
		{
			name: "Test getUserSegmentsHandler doesn't perform read from database and return '405 Method not allowed' and 405 status code",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodPost, "/get-user-segments", strings.NewReader(`{"user_id": 111}`)),
			},
			expectedOutput: `405 Method not allowed`,
			expectedCode:   http.StatusMethodNotAllowed,
		},
		{
			name: "Test deleteSegmentHandler doesn't perform delete in database and return 'Invalid JSON' and 400 status code",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(http.MethodGet, "/get-user-segments", strings.NewReader(`{"user": 111`)),
			},
			expectedOutput: `{"status":"Invalid JSON"}`,
			expectedCode:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				service: mock_service.NewMockService(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			s := NewServer(
				context.Background(),
				&Config{Port: ""},
				db.NewDBStub(),
				f.service,
			)

			s.getUserSegmentsHandler(tt.args.writer, tt.args.request)
			res := tt.args.writer.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if string(data) != tt.expectedOutput {
				t.Errorf("Output -- WANT: %v GOT: %v", tt.expectedOutput, string(data))
			}
			if res.StatusCode != tt.expectedCode {
				t.Errorf("StatusCode -- WANT: %v GOT: %v", tt.expectedCode, strconv.Itoa(res.StatusCode))
			}
		})
	}
}
