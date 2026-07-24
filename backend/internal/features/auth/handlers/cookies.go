package handlers



import (

	"net/http"



	"github.com/PococodoOrg/PocoClinic/internal/pkg/config"

	"github.com/gin-gonic/gin"

)



const (

	authCookiePath = "/api/v1/auth"

	apiCookiePath  = "/api/v1"

)



func setRefreshTokenCookie(c *gin.Context, cfg config.AuthConfig, token string) {

	if token == "" {

		return

	}



	maxAge := int(cfg.RefreshTokenTTL.Seconds())

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(

		cfg.RefreshTokenCookieName,

		token,

		maxAge,

		authCookiePath,

		"",

		cfg.CookieSecure,

		true,

	)

}



func setAccessTokenCookie(c *gin.Context, cfg config.AuthConfig, token string) {

	if token == "" {

		return

	}



	maxAge := int(cfg.AccessTokenTTL.Seconds())

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(

		cfg.AccessTokenCookieName,

		token,

		maxAge,

		apiCookiePath,

		"",

		cfg.CookieSecure,

		true,

	)

}



func clearRefreshTokenCookie(c *gin.Context, cfg config.AuthConfig) {

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(

		cfg.RefreshTokenCookieName,

		"",

		-1,

		authCookiePath,

		"",

		cfg.CookieSecure,

		true,

	)

}



func clearAccessTokenCookie(c *gin.Context, cfg config.AuthConfig) {

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(

		cfg.AccessTokenCookieName,

		"",

		-1,

		apiCookiePath,

		"",

		cfg.CookieSecure,

		true,

	)

}



func readRefreshTokenCookie(c *gin.Context, cfg config.AuthConfig) string {

	token, err := c.Cookie(cfg.RefreshTokenCookieName)

	if err != nil {

		return ""

	}

	return token

}



func readAccessTokenCookie(c *gin.Context, cfg config.AuthConfig) string {

	token, err := c.Cookie(cfg.AccessTokenCookieName)

	if err != nil {

		return ""

	}

	return token

}



func setAuthCookies(c *gin.Context, cfg config.AuthConfig, accessToken, refreshToken string) {

	setAccessTokenCookie(c, cfg, accessToken)

	setRefreshTokenCookie(c, cfg, refreshToken)

}



func clearAuthCookies(c *gin.Context, cfg config.AuthConfig) {

	clearAccessTokenCookie(c, cfg)

	clearRefreshTokenCookie(c, cfg)

}


