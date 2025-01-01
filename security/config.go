package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

// Add to both agent and manager
const (
	HMAC_KEY_SIZE = 32 // 256 bits
)

type SecurityConfig struct {
	SecretKey []byte
}

// Function to create HMAC
func CreateHMAC(message string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// Function to verify HMAC
func VerifyHMAC(message, messageHMAC string, key []byte) bool {
	expectedHMAC := CreateHMAC(message, key)
	return hmac.Equal([]byte(messageHMAC), []byte(expectedHMAC))
}

// Modified ping message structure
type SecurePing struct {
	IP   string
	Time string
	HMAC string
}

type AgentConfig struct {
	SecretKey string `json:"secret_key"`
}

func LoadConfig(filename string) (*AgentConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	var config AgentConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config: %v", err)
	}

	if len(config.SecretKey) != HMAC_KEY_SIZE*2 { // *2 because hex encoded
		return nil, fmt.Errorf("secret key must be %d bytes (hex encoded)", HMAC_KEY_SIZE)
	}

	return &config, nil
}
