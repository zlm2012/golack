package rtmapi

// Channel represents Slack channel.
type Channel struct {
	Name string
}

// UnmarshalText parses a given slack chanel information to Channel
// This method is mainly used by encode/json.
func (channel *Channel) UnmarshalText(b []byte) error {
	str := string(b)

	channel.Name = str

	return nil
}

// MarshalText returns the stringified value of Channel
func (channel *Channel) MarshalText() ([]byte, error) {
	return []byte(channel.Name), nil
}
