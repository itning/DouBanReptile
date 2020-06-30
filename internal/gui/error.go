package gui

type ErrorHandler struct {
}

func (deh ErrorHandler) Handler(e error) {
	if e != nil {
		Print(e.Error())
	}
}
