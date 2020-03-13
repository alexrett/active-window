package activeWindow

type ActiveWindow struct{}

func (a *ActiveWindow) GetActiveWindowTitle() (string, string) {
	return a.getActiveWindowTitle()
}
