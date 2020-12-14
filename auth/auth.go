package auth

const (
	SUB = "1"
	PUB = "2"
)

func CheckTopicAuth(action, clientID, username, topic string) bool {
	return true
}

func CheckConnectAuth(clientID, username, password string) bool {
	return true
}
