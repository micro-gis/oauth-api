package access_token

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccessTokenConstants (t *testing.T) {
	assert.EqualValues(t, 24, expirationTime, "Expiration time should be 24 hours")
}

func TestGetNewAccessToken(t *testing.T) {
	at:= GetNewAccessToken()
	assert.False(t, at.IsExpired(), "Brand new access token should not be expired")
	assert.EqualValues(t, "",at.AccessToken, "Brand new access token should not be expired")
	assert.EqualValues(t, 0, at.UserId, "new access token should not have an associated user id")
}

func TestAccessToken_IsExpired(t *testing.T) {
	at:= AccessToken{}
	assert.True(t, at.IsExpired(), "empty access token should be expired by default")
	at.Expires = time.Now().UTC().Add(3*time.Hour).Unix()
	assert.False(t, at.IsExpired(), "AT expiring in 3 hours from now should NOT be expired")
}
