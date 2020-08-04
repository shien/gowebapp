package main

import (
	gomniauthtest "github.com/stretchr/gomniauth/test"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar

	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}

	url, err := authAvatar.GetAvatarURL(testChatUser)

	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、AuthAvatar.GetAvatarURL は" +
			"ErrNoAvatarURL を返すべき")
	}

	testUrl := "http://url-to-avatar/"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testUrl, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)

	if err != nil {
		t.Error("値が存在する場合、AuthAvatar.GetAvatarURLは" +
			"値を返すべきではありません")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURL は正しい値を返すべきです")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return error.")
	}
	if url != "https://www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvatar.GetAvatarURL が %s という誤った値を返しました", url)
	}

}

func TestFileSystemAvatar(t *testing.T) {

	filename := filepath.Join("avatars", "avc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar

	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatasrURLは'%s'という誤った値を返しました", url)
	}
}
