package jwt

import (
	"time"

	"github.com/zicare/rgm/config"
)

// Keeps a registry of users with a JWT revoke alert.
var revokedJWTMap map[string]map[int64]int64

// Initializes and cleans up the revokedJWTMap registry.
// revokedJWTMap keys represent the users and the values
// a revoke alert timestamp. Any JWT issued
// before said timestamp will be reported as revoked
// by IsRevoked.
// Entries' lifetime is equal to the JWT lifetime,
// this garantees that all JWT issued before the revoke alert
// will be reported as revoked. Obsolete entries are deleted
// automatically by the go-routine.
// It is the responsability of the client app to add
// the entries.
// JWT lifetime is set in the configuration files.
func Init() {

	revokedJWTMap = map[string]map[int64]int64{}

	go func() {
		mcl := time.Duration(60) * time.Second
		jwtDuration, _ := time.ParseDuration(config.Config().GetString("jwt_duration"))
		for {
			for k1, v1 := range revokedJWTMap {
				for k2, v2 := range v1 {
					if v2 < time.Now().Add(-1*jwtDuration).Unix() {
						delete(revokedJWTMap[k1], k2)
					}
				}
			}
			time.Sleep(mcl)
		}
	}()
}

//RevokeJWT exported
func RevokeJWT(src string, id int64) {

	if _, ok := revokedJWTMap[src]; !ok {
		revokedJWTMap[src] = map[int64]int64{}
	}

	revokedJWTMap[src][id] = time.Now().Unix()
}

// RevokedJWTReset exported
func RevokedJWTReset() {

	revokedJWTMap = map[string]map[int64]int64{}
}

// IsRevoked exported
func IsRevoked(payload Payload) bool {

	ts, revoked := revokedJWTMap[payload.Src][payload.Id]

	if revoked {
		return payload.Iat < ts
	}
	return false
}
