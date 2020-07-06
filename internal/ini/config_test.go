package ini

import (
	"fmt"
	"github.com/itning/DouBanReptile/internal/gui"
	"testing"
)

func TestRead(t *testing.T) {
	var config = Config{}
	preference := config.Read()
	fmt.Println(*preference)
}

func TestWrite1(t *testing.T) {
	var config = Config{}
	var preference = gui.Preference{
		GroupEntityURL:             "a",
		MaxPrice:                   0,
		IncludeNoContentPriceCheck: true,
		ExcludeKeyArray:            []string{"a", "cc"},
		IncludeKeyArray:            []string{},
		MaxPage:                    1,
		SavePreference:             true,
	}
	config.Write(&preference)
}
