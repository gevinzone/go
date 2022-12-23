package person

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_GetUserInfo(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockMale := NewMockPerson(ctl)
	var id int64 = 1
	mockMale.EXPECT().Get(id).Return(mockMale, nil)
	user := NewUser(mockMale)
	_, err := user.GetPersonInfo(id)
	assert.NoError(t, err)
}
