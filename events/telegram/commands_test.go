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

func TestCutTextToData(t *testing.T) {
	str := "change_wishes 5"
	command, id := cutTextToData(str)
	assert(t, command, "change_wishes")
	assert(t, id, 5)

}

func TestCheckIfStartHasID(t *testing.T) {
	req1 := "/start 1234567"
	ok, resp1 := checkIfStartHasID(req1)
	assert(t, ok, true)
	assert(t, resp1, "1234567")

	req2 := "/help"
	ok, resp2 := checkIfStartHasID(req2)
	assert(t, ok, false)
	assert(t, resp2, "")

	req3 := "Привіт, мене звати Вася"
	ok, resp3 := checkIfStartHasID(req3)
	assert(t, ok, false)
	assert(t, resp3, "")
}
