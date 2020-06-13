package eventsapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

// RequestValidator receives requests from Slack and check its validity.
type RequestValidator interface {
	Validate(*SlackRequest) bool
}

// SignatureValidator verifies Slack request.
// See https://api.slack.com/authentication/verifying-requests-from-slack
type SignatureValidator struct {
	Secret string
}

var _ RequestValidator = (*SignatureValidator)(nil)

// Validate checks the validity of given request.
// This returns true when the incoming request is valid.
func (rv *SignatureValidator) Validate(request *SlackRequest) bool {
	hash := hmac.New(sha256.New, []byte(rv.Secret))
	fmt.Fprintf(hash, "v0:%d:%s", request.TimeStamp.Unix(), request.Payload)
	expected := fmt.Sprintf("v0=%x", hash.Sum(nil))
	return request.Signature == expected
}
