package middleware

/*
import (
	"com.setlog/internal/configuration"
	"gitlab.setlog.lan/osca-dc/golang/auth.git"
	"log/slog"
	"net/http"
)

func AuthMiddleware(config *configuration.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authServer := auth.AuthServer{
				Host:      config.OAuthIssuer,
				Realm:     "DC",
				ClientIds: config.OAuthClientIds,
			}
			req := auth.GetLoginDataFromHTTPRequest(r)

			_, err := authServer.ValidateRequest(req)
			if err != nil {
				slog.Warn("Unauthorized", err)
				http.Error(w, "Unauthorized", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
*/
