package domain

type OAuthUserInfo struct {
	ID            string
	Email         string
	FirstName     string
	LastName      string
	AvatarURL     string
	EmailVerified bool
}
