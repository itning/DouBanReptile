package ini

import (
	"fmt"
	"github.com/itning/DouBanReptile/internal/preference"
	"testing"
)

func TestRead(t *testing.T) {
	var config = Config{}
	pre := config.Read()
	fmt.Println(*pre)
}

func TestWrite1(t *testing.T) {
	var config = Config{}
	var pre = preference.Preference{
		GroupEntityURL:             "a",
		MaxPrice:                   0,
		IncludeNoContentPriceCheck: true,
		ExcludeKeyArray:            []string{"a", "cc"},
		IncludeKeyArray:            []string{},
		MaxPage:                    1,
		SavePreference:             true,
	}
	config.Write(&pre)
}
