package tokens

import (
	"time"

	"github.com/CXTACLYSM/hiring-api/internal/auth/user/domain/entities"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

type JwtTokenGenerator struct {
	secretKey []byte
}

func NewJwtTokenGenerator(secretKey string) *JwtTokenGenerator {
	return &JwtTokenGenerator{secretKey: []byte(secretKey)}
}

func (g *JwtTokenGenerator) Generate(user *entities.User) (string, error) {
	var userId string
	if user != nil {
		userId = user.Id
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "my-service",
		},
	})
	jwt.Parse()

	return token.SignedString(g.secretKey)
}
