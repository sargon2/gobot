package gobot

type MessageSource struct {
	ChannelID string
	Username  string
	response  string // for testing
}

// For testing
func (m *MessageSource) GetResponse() string {
	return m.response
}

func (m *MessageSource) SetResponse(t string) {
	m.response = t
}
