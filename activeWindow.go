package activeWindow

type ActiveWindow struct {}

func (a *ActiveWindow) GetActiveWindowTitle() string {
	return a.getActiveWindowTitle()
}