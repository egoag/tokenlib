package tokenlib

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/hkdf"
)

func toInt64(v interface{}) (int64, error) {
	switch i := v.(type) {
	case int:
		return int64(i), nil
	case int64:
		return int64(i), nil
	case float32:
		return int64(i), nil
	case float64:
		return int64(i), nil
	default:
		return 0, errors.New("not Number")
	}
}

func signSecret(secret string, len int) []byte {
	hash := sha256.New
	info := []byte("services.mozilla.com/tokenlib/v1/signing")
	hkdf := hkdf.New(hash, []byte(secret), nil, info)
	result := make([]byte, len)
	hkdf.Read(result)
	return result
}

func generateSig(secret, payload []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	sig := mac.Sum(nil)
	return sig
}

// MakeToken generate token with given data v, secret and timeout
func MakeToken(v map[string]interface{}, secret string, timeout int64) (string, error) {
	if _, ok := v["salt"].(string); !ok {
		salt := make([]byte, 3)
		rand.Read(salt)
		dstSalt := make([]byte, hex.EncodedLen(3))
		hex.Encode(dstSalt, salt)
		v["salt"] = dstSalt
	}

	if _, ok := v["expires"].(int); !ok {
		expires := int64(time.Now().Unix()) + timeout
		v["expires"] = expires
	}

	payload, err := json.Marshal(v)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	sig := generateSig(signSecret(secret, sha256.Size), payload)
	raw := append(payload, sig[:]...)
	token := base64.URLEncoding.EncodeToString(raw)

	return token, nil
}

// ParseToken parse string token with secret,
// return error when secret dismatch or expires
func ParseToken(token string, secret string) (map[string]interface{}, error) {
	decodedToken, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	sig := decodedToken[len(decodedToken)-sha256.Size:]
	payload := decodedToken[:len(decodedToken)-sha256.Size]
	realSig := generateSig(signSecret(secret, sha256.Size), payload)

	if !bytes.Equal(sig, realSig) {
		return nil, errors.New("invalid signiture")
	}

	var v interface{}
	err = json.Unmarshal(payload, &v)
	if err != nil {
		return nil, err
	}

	data, ok := v.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid json")
	}

	expires, err := toInt64(data["expires"])
	if err != nil {
		return nil, errors.New("invalid json")
	}
	if expires < time.Now().Unix() {
		return nil, errors.New("token has expired")
	}

	delete(data, "salt")
	delete(data, "expires")

	return data, nil
}
