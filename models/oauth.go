package models

import (
	"database/sql"
	"time"

	"github.com/tientruongcao51/oauth2-sever/util"
	"github.com/tientruongcao51/oauth2-sever/uuid"
)

// OauthClient ...
type OauthClient struct {
	MyGormModel
	Name        string //unique;not null
	Key         string //unique;not null
	Secret      string //not null
	Mail        string //not null
	RedirectURI string
}

// OauthClient ...
type MailToken struct {
	MyGormModel
	Mail  string //not null
	Token string
}

// TableName specifies table name
func (c *OauthClient) TableName() string {
	return "oauth_clients"
}

// OauthScope ...
type OauthScope struct {
	MyGormModel
	Scope       string //unique;not null
	Description sql.NullString
	IsDefault   bool //default:false
}

// TableName specifies table name
func (s *OauthScope) TableName() string {
	return "oauth_scopes"
}

// OauthRole is a one of roles user can have (currently superuser or user)
type OauthRole struct {
	TimestampModel
	ID   string //primary_key
	Name string // unique;not null
}

// TableName specifies table name
func (r *OauthRole) TableName() string {
	return "oauth_roles"
}

// OauthUser ...
type OauthUser struct {
	MyGormModel
	RoleID   sql.NullString //index;not null
	Role     *OauthRole
	Username string //unique;not null
	Mail     string //unique;not null
	Password sql.NullString
}

// TableName specifies table name
func (u *OauthUser) TableName() string {
	return "oauth_users"
}

// OauthRefreshToken ...
type OauthRefreshToken struct {
	MyGormModel
	BsKey     string
	ClientID  sql.NullString //not null
	UserID    sql.NullString
	Client    *OauthClient
	User      *OauthUser
	Token     string //unique;not null
	ExpiresAt time.Time
	Scope     string //not null
}

// TableName specifies table name
func (rt *OauthRefreshToken) TableName() string {
	return "oauth_refresh_tokens"
}

// OauthAccessToken ...
type OauthAccessToken struct {
	MyGormModel
	BsKey     string
	ClientID  sql.NullString //index;not null
	UserID    sql.NullString
	Client    *OauthClient
	User      *OauthUser
	Token     string //unique;not null
	ExpiresAt time.Time
	Scope     string //not null
}

// TableName specifies table name
func (at *OauthAccessToken) TableName() string {
	return "oauth_access_tokens"
}

// OauthAuthorizationCode ...
type OauthAuthorizationCode struct {
	MyGormModel
	BsKey       string
	ClientID    sql.NullString //not null
	UserID      sql.NullString //not null
	Client      *OauthClient
	User        *OauthUser
	Code        string //unique;not null
	RedirectURI sql.NullString
	ExpiresAt   time.Time
	Scope       string //not null
}

// TableName specifies table name
func (ac *OauthAuthorizationCode) TableName() string {
	return "oauth_authorization_codes"
}

// NewOauthRefreshToken creates new OauthRefreshToken instance
func NewOauthRefreshToken(client *OauthClient, user *OauthUser, expiresIn int, scope string) *OauthRefreshToken {
	token := uuid.New()
	bsKey := GetItemKeyRefreshToken(client.ID, "")
	refreshToken := &OauthRefreshToken{
		MyGormModel: MyGormModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		Client:    client,
		ClientID:  util.StringOrNull(string(client.ID)),
		Token:     token,
		ExpiresAt: time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		Scope:     scope,
	}
	if user != nil {
		refreshToken.UserID = util.StringOrNull(string(user.ID))
		refreshToken.User = user
		bsKey = GetItemKeyRefreshToken(client.ID, user.ID)
	}
	refreshToken.BsKey = bsKey
	return refreshToken
}

func GetItemKeyRefreshToken(clientId string, userId string) string {
	return clientId + "_" + userId
}

// NewOauthAccessToken creates new OauthAccessToken instance
func NewOauthAccessToken(client *OauthClient, user *OauthUser, expiresIn int, scope string) *OauthAccessToken {
	token := uuid.New()
	bsKey := GetItemKeyAccessToken(client.ID, "")
	accessToken := &OauthAccessToken{
		MyGormModel: MyGormModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		Client:    client,
		ClientID:  util.StringOrNull(string(client.ID)),
		Token:     token,
		ExpiresAt: time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		Scope:     scope,
	}
	if user != nil {
		accessToken.UserID = util.StringOrNull(string(user.ID))
		accessToken.User = user
		bsKey = GetItemKeyAccessToken(client.ID, user.ID)
	}
	accessToken.BsKey = bsKey

	return accessToken
}

func GetItemKeyAccessToken(clientId string, userId string) string {
	return clientId + "_" + userId
}

// NewOauthAuthorizationCode creates new OauthAuthorizationCode instance
func NewOauthAuthorizationCode(client *OauthClient, user *OauthUser, expiresIn int, redirectURI, scope string) *OauthAuthorizationCode {
	code := uuid.New()
	bsKey := GetItemKeyAuthorizationToken(client.ID, user.ID)
	return &OauthAuthorizationCode{
		MyGormModel: MyGormModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		BsKey:       bsKey,
		ClientID:    util.StringOrNull(string(client.ID)),
		UserID:      util.StringOrNull(string(user.ID)),
		User:        user,
		Client:      client,
		Code:        code,
		ExpiresAt:   time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		RedirectURI: util.StringOrNull(redirectURI),
		Scope:       scope,
	}
}

func GetItemKeyAuthorizationToken(clientId string, userId string) string {
	return clientId + "_" + userId
}

/*// OauthAuthorizationCodePreload sets up Gorm preloads for an auth code object
func OauthAuthorizationCodePreload(db *gorm.DB) *gorm.DB {
	return OauthAuthorizationCodePreloadWithPrefix(db, "")
}

// OauthAuthorizationCodePreloadWithPrefix sets up Gorm preloads for an auth code object,
// and prefixes with prefix for nested objects
func OauthAuthorizationCodePreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.
		Preload(prefix + "Client").Preload(prefix + "User")
}

// OauthAccessTokenPreload sets up Gorm preloads for an access token object
func OauthAccessTokenPreload(db *gorm.DB) *gorm.DB {
	return OauthAccessTokenPreloadWithPrefix(db, "")
}

// OauthAccessTokenPreloadWithPrefix sets up Gorm preloads for an access token object,
// and prefixes with prefix for nested objects
func OauthAccessTokenPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.
		Preload(prefix + "Client").Preload(prefix + "User")
}

// OauthRefreshTokenPreload sets up Gorm preloads for a refresh token object
func OauthRefreshTokenPreload(db *gorm.DB) *gorm.DB {
	return OauthRefreshTokenPreloadWithPrefix(db, "")
}

// OauthRefreshTokenPreloadWithPrefix sets up Gorm preloads for a refresh token object,
// and prefixes with prefix for nested objects
func OauthRefreshTokenPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.
		Preload(prefix + "Client").Preload(prefix + "User")
}
*/
