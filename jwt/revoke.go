package jwt

import (
	"time"

	"github.com/zicare/rgm/config"
)

// Keeps a registry of users with a JWT revoke alert.
var revokedJWTMap map[string]map[string]time.Time

// Initializes and cleans up the revokedJWTMap registry.
// revokeJWTMap negates the use of fresh JWTs issued before access revoke.
// revokedJWTMap keys represent the users and the values a revoke alert timestamp.
// Any JWT issued before said timestamp will be reported as revoked by IsRevoked.
// Entries' lifetime is equal to the JWT lifetime,
// this garantees that all JWT issued before the revoke alert
// will be reported as revoked. Obsolete entries are deleted
// automatically by the go-routine.
// It is the responsability of the client app to add the entries.
// JWT lifetime is set in the configuration files.
func Init() {

	revokedJWTMap = map[string]map[string]time.Time{}

	go func() {
		mcl := time.Duration(60) * time.Second
		jwtDuration, _ := time.ParseDuration(config.Config().GetString("jwt_duration"))
		for {
			for k1, v1 := range revokedJWTMap {
				for k2, v2 := range v1 {
					if v2.Before(time.Now().Add(-1 * jwtDuration)) {
						delete(revokedJWTMap[k1], k2)
					}
				}
			}
			time.Sleep(mcl)
		}
	}()
}

//RevokeJWT exported
func RevokeJWT(t string, uid string) {

	if _, ok := revokedJWTMap[t]; !ok {
		revokedJWTMap[t] = map[string]time.Time{}
	}

	revokedJWTMap[t][uid] = time.Now()
}

// RevokedJWTReset exported
func RevokedJWTReset() {

	revokedJWTMap = map[string]map[string]time.Time{}
}

// IsRevoked exported
func IsRevoked(payload Payload) bool {

	ts, revoked := revokedJWTMap[payload.Type][payload.UID]

	if revoked {
		return ts.After(payload.Iat)
	}
	return false
}
