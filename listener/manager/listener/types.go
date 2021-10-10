package listener

type ListenerManager interface {
	DecryptAndEmit(splitArray []string)
}
