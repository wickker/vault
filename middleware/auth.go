package middleware

import (
	"encoding/json"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"net/http"
	"vault/openapi"
)

func Auth() gin.HandlerFunc {
	return gin.WrapH(clerkhttp.WithHeaderAuthorization()(auth()))
}

func auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := clerk.SessionClaimsFromContext(ctx)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			writeError(w, "Unauthorized access")
			return
		}

		usr, err := user.Get(ctx, claims.Subject)
		if err != nil {
			log.Err(err).Msg("Unable to get user from Clerk.")
			w.WriteHeader(http.StatusUnauthorized)
			writeError(w, err.Error())
			return
		}
		if usr == nil {
			w.WriteHeader(http.StatusUnauthorized)
			writeError(w, "User does not exist")
			return
		}

		setContextValue(r, "userID", usr.ID)
		if len(usr.EmailAddresses) > 0 {
			setContextValue(r, "userEmail", usr.EmailAddresses[0].EmailAddress)
		}
	}
}

func writeError(w http.ResponseWriter, message string) {
	errMsg := openapi.Error{
		Message: message,
	}
	bytes, _ := json.Marshal(errMsg)
	_, _ = w.Write(bytes)
}

func setContextValue(r *http.Request, key, value string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}
