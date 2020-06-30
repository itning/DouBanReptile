package error2

import "github.com/itning/DouBanReptile/internal/log"

type Error interface {
	Handler(e error)
}

var handlerImpl Error = DefaultErrorHandler{}

type DefaultErrorHandler struct {
}

func (deh DefaultErrorHandler) Handler(e error) {
	if e != nil {
		log.GetImpl().Printf("Have Error %s", e.Error())
		panic(e)
	}
}

func SetImpl(e Error) {
	handlerImpl = e
}

func GetImpl() Error {
	return handlerImpl
}
