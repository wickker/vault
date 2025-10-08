package middleware

import (
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"vault/openapi"
)

func Auth(frontendOrigins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.Request.URL.Path, "protected") {
			return
		}

		token := strings.TrimSpace(c.GetHeader("Authorization"))
		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := jwt.Verify(c, &jwt.VerifyParams{
			Token: token,
			AuthorizedPartyHandler: func(azpClaim string) bool {
				return strings.Contains(strings.ToLower(frontendOrigins),
					strings.ToLower(azpClaim))
			},
		})
		if err != nil {
			log.Err(err).Msg("Unable to verify Clerk JWT.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, openapi.Error{
				Message: err.Error(),
			})
			return
		}

		u, err := user.Get(c, claims.Subject)
		if err != nil {
			log.Err(err).Msg("Unable to get user from Clerk.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, openapi.Error{
				Message: err.Error(),
			})
			return
		}
		if u == nil {
			log.Error().Msg("Clerk user is nil.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, openapi.Error{
				Message: "Clerk user is nil",
			})
			return
		}

		c.Set(ContextKeys.User, u)
		c.Next()
	}
}
