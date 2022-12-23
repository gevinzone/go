package example1

import (
	person_mock "github.com/gevinzone/go/testing/gomock/example1/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_GetUserInfo(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockMale := person_mock.NewMockPerson(ctl)
	expectName := "gevin"
	mockMale.EXPECT().GetName(gomock.Any()).Return(expectName, nil)
	user := NewUser(mockMale)
	name, err := user.GetName(1)
	assert.NoError(t, err)
	assert.Equal(t, expectName, name)
}
