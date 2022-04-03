package gobot

type HistoryGrabber struct {
}

func NewHistoryGrabber() *HistoryGrabber {
	grabber := &HistoryGrabber{}
	grabber.grabHistory()
	return grabber
}

func (*HistoryGrabber) grabHistory() {
	print("Hello, world!\n")
}
