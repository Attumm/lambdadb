package main

import (
	"strconv"
)

type registerCustomGroupByFunc map[string]func(*Item, ItemsGroupedBy)

var RegisterGroupByCustom registerCustomGroupByFunc

func init() {

	RegisterGroupByCustom = make(registerCustomGroupByFunc)
	RegisterGroupByCustom["gebruiksdoelen-mixed"] = GroupByGettersGebruiksdoelen

}

func reduceWEQ(items Items) map[string]string {
	result := make(map[string]string)
	weq := 0
	for i := range items {
		_weq, err := strconv.ParseInt(Woningequivalent.GetValue(items[i].Woningequivalent), 10, 64)
		if err != nil {
			panic(err)
		}
		weq += int(_weq)
	}
	result["woningenquivalent"] = strconv.Itoa(weq)
	return result
}
