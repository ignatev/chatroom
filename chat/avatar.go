package main

import "errors"

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get na avatar URL")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}
var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}