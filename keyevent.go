package minerva

// KeyEvent is generic struct to normalize all events happening in remote config
// (could be redis, etcd and so on)
type KeyEvent struct {
	AffectedKey string
	Type        string
	Value       string
}
