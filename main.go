package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"google.golang.org/api/idtoken"

	"github.com/dotenv-org/godotenvvault"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/oneto6/vcp-be/meet"
	"resty.dev/v3"
)

type OAuth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	IDToken     string `json:"id_token"`
}

const ADDR = "127.0.0.1:4373"
const domainBaseURL = `https://swim-brochure-accessed-aluminium.trycloudflare.com`

const pathAuthGoogleCallback = `/auth/google/callback`

type Config struct {
}

func block() (err error) {

	err = godotenvvault.Load()
	if err != nil {
		return
	}
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		err = fmt.Errorf("clientID or clientSecret is empty")
		return
	}
	fmt.Println(clientID, clientSecret)

	engin := gin.New()
	engin.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))
	engin.GET("/", func(c *gin.Context) {
		c.String(200, `Hello World`)
	})

	engin.GET(pathAuthGoogleCallback, func(c *gin.Context) {
		code := c.Query(`code`)
		if code == `` {
			gin.DefaultWriter.Write([]byte(`code isn't present;`))
			c.Status(http.StatusExpectationFailed)
			return
		}
		uri := `https://oauth2.googleapis.com/token`

		var formData = map[string]string{
			`client_id`:     clientID,
			`client_secret`: clientSecret,
			`grant_type`:    `authorization_code`,
			`code`:          code,
			`redirect_uri`:  domainBaseURL + pathAuthGoogleCallback,
		}
		res, err := resty.New().R().SetFormData(formData).Post(uri)
		if err != nil {
			c.Status(http.StatusExpectationFailed)
			return
		}
		gin.DefaultWriter.Write(res.Bytes())
		tokenResponse := OAuth2TokenResponse{}
		err = json.Unmarshal(res.Bytes(), &tokenResponse)
		if err != nil {
			c.Status(http.StatusExpectationFailed)
			return
		}
		payload, err := idtoken.Validate(c, tokenResponse.IDToken, clientID)
		if err != nil {
			c.Status(http.StatusExpectationFailed)
			return
		}
		fmt.Fprintln(gin.DefaultWriter, payload)
	})

	engin.GET("/auth/google", func(c *gin.Context) {
		uri := url.URL{
			Scheme: `https`,
			Host:   `accounts.google.com`,
			Path:   `o/oauth2/v2/auth`,
		}
		query := url.Values{}

		// const scope =
		//     "openid https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile";
		//     '?client_id=$clientId'
		//     '&redirect_uri=$redirectUri'
		//     '&response_type=token'
		//     '&scope=$scope';

		redirectURI := domainBaseURL + pathAuthGoogleCallback
		query.Set(`client_id`, clientID)
		query.Set(`redirect_uri`, redirectURI)
		query.Set(`response_type`, `code`)
		query.Set(`scope`, `openid https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile`)
		uri.RawQuery = query.Encode()
		c.Redirect(http.StatusTemporaryRedirect, uri.String())
	})
	engin.GET("/meet", func(c *gin.Context) {
		meet := meet.GetMeet()
		c.JSON(200, meet)
	})
	engin.Run(ADDR)
	return
}

func main() {
	err := block()
	if err != nil {
		panic(err)
	}
}
