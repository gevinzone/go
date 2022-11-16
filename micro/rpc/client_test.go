package rpc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitClientProxy(t *testing.T) {
	testCases := []struct {
		name        string
		service     *UserService
		p           *mockProxy
		wantReq     *Request
		wantInitErr error
		wantErr     error
		wantResp    *GetByIdResp
	}{
		{
			name:    "user service",
			service: &UserService{},
			p: &mockProxy{
				result: []byte(`{"name":"Tom"}`),
			},
			wantReq: &Request{
				ServiceName: "user-service",
				MethodName:  "GetById",
				Data:        []byte(`{"Id":0}`),
			},
			wantResp: &GetByIdResp{
				Name: "Tom",
			},
		},
		{
			name: "proxy return error",
			p: &mockProxy{
				err: errors.New("mock error"),
			},
			service: &UserService{},
			wantErr: errors.New("mock error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := InitClientProxy(tc.service, tc.p)
			assert.Equal(t, tc.wantInitErr, err)
			if err != nil {
				return
			}
			resp, err := tc.service.GetById(context.Background(), &GetByIdReq{})
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantReq.Data, tc.p.req.Data)
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}

type mockProxy struct {
	req    *Request
	err    error
	result []byte
}

func (m *mockProxy) Invoke(_ context.Context, req *Request) (*Response, error) {
	m.req = req
	return &Response{Data: m.result}, m.err
}

type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

func (u *UserService) Name() string {
	return "user-service"
}

type GetByIdReq struct {
	Id int64
}

type GetByIdResp struct {
	Name string `json:"name"`
}
