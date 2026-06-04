package github

import (
	"crypto/hmac"
	"encoding/hex"
	"hash"
)

// validates the GitHub webhook signature.
func VerifyHMAC(body []byte, secret, providedHex string, hashFunc func() hash.Hash) bool {
	provided, err := hex.DecodeString(providedHex)
	if err != nil {
		return false
	}
	mac := hmac.New(hashFunc, []byte(secret))
	mac.Write(body)
	expected := mac.Sum(nil)
	return hmac.Equal(expected, provided)
}

// pulls the clone_url from the GitHub webhook payload.
func ExtractRepoURL(payload map[string]interface{}) string {
	if repo, ok := payload["repository"]; ok {
		if repoMap, ok := repo.(map[string]interface{}); ok {
			if url, ok := repoMap["clone_url"]; ok {
				if s, ok := url.(string); ok {
					return s
				}
			}
		}
	}
	return ""
}

// pulls the html_url of the pull request from the payload.
func ExtractPRURL(payload map[string]interface{}) string {
	if pr, ok := payload["pull_request"]; ok {
		if prMap, ok := pr.(map[string]interface{}); ok {
			if url, ok := prMap["html_url"]; ok {
				if s, ok := url.(string); ok {
					return s
				}
			}
		}
	}
	return ""
}
