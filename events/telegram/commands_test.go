package telegram

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

func TestExtractIDFromStringSettings(t *testing.T) {
	str := "Налаштування гри: Ім‘я (ID:12345)"
	res := ExtractIDFromStringSettings(str)
	fmt.Println(res)
	assert(t, res, "12345")

	if len(str) > 17 {
		asRunes := []rune(str)
		reqStr := string(asRunes[:17])
		fmt.Printf("%+v\n", reqStr)
		assert(t, string(reqStr), "Налаштування гри:")
		if reqStr == "Налаштування гри:" {
			fmt.Println("Yes")
		}
	}
}
