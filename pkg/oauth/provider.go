package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	apiv1 "github.com/welaw/welaw/api/v1"
	"github.com/welaw/welaw/pkg/errs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/amazon"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
)

const (
	// google
	ProviderGoogle    = "google"
	GoogleURL         = "https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s"
	ProviderAmazon    = "amazon"
	AmazonURL         = "https://api.amazon.com/user/profile?access_token=%s"
	ProviderMicrosoft = "microsoft"
	MicrosoftURL      = "https://apis.live.net/v5.0/me?access_token=%s"
)

type Provider interface {
	GetLoginURL(state string) (url string)
	LoginUser(state, code string) (*apiv1.User, error)
}

type provider struct {
	ident          string
	accessTokenURL string
	oauthConfig    *oauth2.Config
}

func NewProviderAmazon(
	redirectURL string,
	clientId string,
	clientSecret string,
) Provider {
	return &provider{
		ident:          ProviderAmazon,
		accessTokenURL: AmazonURL,
		oauthConfig: &oauth2.Config{
			RedirectURL:  redirectURL,
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Scopes: []string{
				"profile",
			},
			Endpoint: amazon.Endpoint,
		},
	}
}

func NewProviderGoogle(
	redirectURL string,
	clientId string,
	clientSecret string,
) Provider {
	return &provider{
		ident:          ProviderGoogle,
		accessTokenURL: GoogleURL,
		oauthConfig: &oauth2.Config{
			RedirectURL:  redirectURL,
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func NewProviderMicrosoft(
	redirectURL string,
	clientId string,
	clientSecret string,
) Provider {
	return &provider{
		ident:          ProviderMicrosoft,
		accessTokenURL: MicrosoftURL,
		oauthConfig: &oauth2.Config{
			RedirectURL:  redirectURL,
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Scopes:       []string{"wl.emails"},
			Endpoint:     microsoft.LiveConnectEndpoint,
		},
	}
}

func (p *provider) GetLoginURL(state string) string {
	return p.oauthConfig.AuthCodeURL(state)
}

func (p *provider) LoginUser(state, code string) (*apiv1.User, error) {
	token, err := p.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	response, err := http.Get(fmt.Sprintf(p.accessTokenURL, token.AccessToken))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	switch p.ident {
	case ProviderGoogle:
		var r GoogleResponse
		err = json.Unmarshal(contents, &r)
		return &apiv1.User{
			ProviderId: r.ID,
			Name:       r.Name,
			Email:      r.Email,
			PictureUrl: r.Picture,
		}, nil
	case ProviderAmazon:
		var r AmazonResponse
		err = json.Unmarshal(contents, &r)
		return &apiv1.User{
			ProviderId: r.UserID,
			Name:       r.Name,
			Email:      r.Email,
		}, nil
	case ProviderMicrosoft:
		var r MicrosoftResponse
		err = json.Unmarshal(contents, &r)
		return &apiv1.User{
			ProviderId: r.ID,
			Name:       r.Name,
			Email:      r.Emails["account"],
		}, nil

	default:
		return nil, errs.ErrBadRequest
	}
}

type GoogleResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type AmazonResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type MicrosoftResponse struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	UserID string            `json:"user_id"`
	Emails map[string]string `json:"emails"`
}
