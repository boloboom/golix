package utilsx

import (
	"encoding/json"
	"testing"
	"time"
)

// 待转换的struct
type User struct {
	Id        uint64
	Name      string
	CreatedAt time.Time
}

// resource资源转换需要定义struct并实现transform方法
type UserResource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (*UserResource) transform(model interface{}) interface{} {
	var user = model.(*User)
	idTransformer := NewIdTransformer()
	return &UserResource{
		Id:   idTransformer.Encode(user.Id),
		Name: user.Name,
	}
}

func TestMakeResource(t *testing.T) {
	result := NewResourceTransformer(new(UserResource)).Make(&User{
		Id:        1,
		Name:      "test_user",
		CreatedAt: time.Now(),
	})
	jsonResult, err := json.Marshal(result)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(jsonResult))
}

func TestCollectionResource(t *testing.T) {
	var users = []*User{
		{
			Id:        1,
			Name:      "test_user1",
			CreatedAt: time.Now(),
		},
		{
			Id:        2,
			Name:      "test_user2",
			CreatedAt: time.Now(),
		},
	}
	result := NewResourceTransformer(new(UserResource)).Collection(users)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(jsonResult))
}
