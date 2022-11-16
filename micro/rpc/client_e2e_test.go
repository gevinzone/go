package rpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient(":8081")
	require.NoError(t, err)
	us := &UserService{}
	err = InitClientProxy(us, c)
	require.NoError(t, err)

	//resp, err := us.GetById(context.Background(), &GetByIdReq{Id: 100})
	resp, err := us.GetById(context.Background(), &GetByIdReq{
		Id: 100,
	})
	require.NoError(t, err)
	t.Log(resp)
}
