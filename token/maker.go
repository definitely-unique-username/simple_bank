package token

import "time"

type Maker interface {
	CreateToken(userId int64, duartion time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
