package main

import (
	"fmt"
	"strconv"
	"strings"
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

func GettersToevoegingen(i *Item) string {
	return Postcode.GetValue(i.Postcode) + " " + Huisnummer.GetValue(i.Huisnummer)
}

// getter Gebruiksdoelen
func GroupByGettersGebruiksdoelen(item *Item, grouping ItemsGroupedBy) {

	for i := range item.Gebruiksdoelen {
		groupkey := Gebruiksdoelen.GetValue(item.Gebruiksdoelen[i])
		grouping[groupkey] = append(grouping[groupkey], item)
	}
}

func GetAdres(i *Item) string {
	adres := fmt.Sprintf("%s %s %s %s %s %s",
		Straat.GetValue(i.Straat),
		Huisnummer.GetValue(i.Huisnummer),
		Huisletter.GetValue(i.Huisletter),
		Huisnummertoevoeging.GetValue(i.Huisnummertoevoeging),
		Postcode.GetValue(i.Postcode),
		Gemeentenaam.GetValue(i.Gemeentenaam))

	adres = strings.ReplaceAll(adres, "  ", " ")
	return adres
}
