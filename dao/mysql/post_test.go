package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
)

func init() {
	dbCfg := settings.MysqlConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "root",
		DbName:       "bluebell",
		Port:         3306,
		MaxOpenConns: 20,
		MaxIdleConns: 10,
	}
	if err := Init(&dbCfg); err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:          123,
		AuthorID:    123,
		CommunityID: 1,
		Title:       "test",
		Content:     "just a test",
	}
	if err := CreatePost(&post); err != nil {
		t.Fatalf("CreatePost insert record into mysql failed, err:%v\n", err)
	}
	t.Logf("CreatePost insert record into mysql success")
}
