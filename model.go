/*
	model.go define the 'items' to store.
	All columns with getters and setters are defined here.

	ItemIn, represent rows from the Input data
	Item, the compact item stored in memmory
	ItemOut, defines how and which fields are exported out
	of the API. It is possible to ignore input columns

	Repeated values are stored in maps with int numbers
	as keys.  Optionally bitarrays are created for reapeated
	column values to do fast bit-wise filtering.

	A S2 geo index in created for lat, lon values.

	Unique values are stored as-is.

	The generated codes leaves room to create custom
	index functions yourself to create an API with an
	< 1 ms response time for your specific needs.

	This codebase solves: I need to have an API on this
	tabular dataset fast!
*/

package main

import (
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/Workiva/go-datastructures/bitarray"
)

type registerGroupByFunc map[string]func(*Item) string
type registerGettersMap map[string]func(*Item) string
type registerReduce map[string]func(Items) map[string]string

type registerBitArray map[string]func(s string) (bitarray.BitArray, error)
type fieldBitarrayMap map[uint32]bitarray.BitArray

type ItemIn struct {
	Pid                     string `json:"pid"`
	Vid                     string `json:"vid"`
	Numid                   string `json:"numid"`
	Straat                  string `json:"straat"`
	Postcode                string `json:"postcode"`
	Huisnummer              string `json:"huisnummer"`
	Huisletter              string `json:"huisletter"`
	Huisnummertoevoeging    string `json:"huisnummertoevoeging"`
	Oppervlakte             string `json:"oppervlakte"`
	Woningequivalent        string `json:"woningequivalent"`
	WoningType              string `json:"woning_type"`
	LabelscoreVoorlopig     string `json:"labelscore_voorlopig"`
	LabelscoreDefinitief    string `json:"labelscore_definitief"`
	Energieklasse           string `json:"energieklasse"`
	Gemeentecode            string `json:"gemeentecode"`
	Gemeentenaam            string `json:"gemeentenaam"`
	Buurtcode               string `json:"buurtcode"`
	Buurtnaam               string `json:"buurtnaam"`
	Wijkcode                string `json:"wijkcode"`
	Wijknaam                string `json:"wijknaam"`
	Provinciecode           string `json:"provinciecode"`
	Provincienaam           string `json:"provincienaam"`
	Point                   string `json:"point"`
	PandGasEanAansluitingen string `json:"pand_gas_ean_aansluitingen"`
	GroupId2020             string `json:"group_id_2020"`
	P6GasAansluitingen2020  string `json:"p6_gas_aansluitingen_2020"`
	P6Gasm32020             string `json:"p6_gasm3_2020"`
	P6Kwh2020               string `json:"p6_kwh_2020"`
	P6TotaalPandoppervlakM2 string `json:"p6_totaal_pandoppervlak_m2"`
	PandBouwjaar            string `json:"pand_bouwjaar"`
	PandGasAansluitingen    string `json:"pand_gas_aansluitingen"`
	Gebruiksdoelen          string `json:"gebruiksdoelen"`
}

type ItemOut struct {
	Pid                     string `json:"pid"`
	Vid                     string `json:"vid"`
	Numid                   string `json:"numid"`
	Straat                  string `json:"straat"`
	Postcode                string `json:"postcode"`
	Huisnummer              string `json:"huisnummer"`
	Huisletter              string `json:"huisletter"`
	Huisnummertoevoeging    string `json:"huisnummertoevoeging"`
	Adres                   string `json:"adres"` // should be removed soon
	Oppervlakte             string `json:"oppervlakte"`
	Woningequivalent        string `json:"woningequivalent"`
	WoningType              string `json:"woning_type"`
	LabelscoreVoorlopig     string `json:"labelscore_voorlopig"`
	LabelscoreDefinitief    string `json:"labelscore_definitief"`
	Energieklasse           string `json:"energieklasse"`
	Gemeentecode            string `json:"gemeentecode"`
	Gemeentenaam            string `json:"gemeentenaam"`
	Buurtcode               string `json:"buurtcode"`
	Buurtnaam               string `json:"buurtnaam"`
	Wijkcode                string `json:"wijkcode"`
	Wijknaam                string `json:"wijknaam"`
	Provinciecode           string `json:"provinciecode"`
	Provincienaam           string `json:"provincienaam"`
	Point                   string `json:"point"`
	PandGasEanAansluitingen string `json:"pand_gas_ean_aansluitingen"`
	GroupId2020             string `json:"group_id_2020"`
	P6GasAansluitingen2020  string `json:"p6_gas_aansluitingen_2020"`
	P6Gasm32020             string `json:"p6_gasm3_2020"`
	P6Kwh2020               string `json:"p6_kwh_2020"`
	P6TotaalPandoppervlakM2 string `json:"p6_totaal_pandoppervlak_m2"`
	PandBouwjaar            string `json:"pand_bouwjaar"`
	PandGasAansluitingen    string `json:"pand_gas_aansluitingen"`
	Gebruiksdoelen          string `json:"gebruiksdoelen"`
}

type Item struct {
	Label                   int // internal index in ITEMS
	Pid                     uint32
	Vid                     uint32
	Numid                   string
	Straat                  uint32
	Postcode                uint32
	Huisnummer              uint32
	Huisletter              uint32
	Huisnummertoevoeging    uint32
	Oppervlakte             uint32
	Woningequivalent        uint32
	WoningType              uint32
	LabelscoreVoorlopig     uint32
	LabelscoreDefinitief    uint32
	Energieklasse           uint32
	Gemeentecode            uint32
	Gemeentenaam            uint32
	Buurtcode               uint32
	Buurtnaam               uint32
	Wijkcode                uint32
	Wijknaam                uint32
	Provinciecode           uint32
	Provincienaam           uint32
	Point                   string
	PandGasEanAansluitingen uint32
	GroupId2020             string
	P6GasAansluitingen2020  uint32
	P6Gasm32020             uint32
	P6Kwh2020               uint32
	P6TotaalPandoppervlakM2 uint32
	PandBouwjaar            uint32
	PandGasAansluitingen    uint32
	Gebruiksdoelen          []uint32
}

func (i Item) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Serialize())
}

// Shrink create smaller Item using uint32
func (i ItemIn) Shrink(label int) Item {

	Pid.Store(i.Pid)
	Vid.Store(i.Vid)
	Straat.Store(i.Straat)
	Postcode.Store(i.Postcode)
	Huisnummer.Store(i.Huisnummer)
	Huisletter.Store(i.Huisletter)
	Huisnummertoevoeging.Store(i.Huisnummertoevoeging)
	Oppervlakte.Store(i.Oppervlakte)
	Woningequivalent.Store(i.Woningequivalent)
	WoningType.Store(i.WoningType)
	LabelscoreVoorlopig.Store(i.LabelscoreVoorlopig)
	LabelscoreDefinitief.Store(i.LabelscoreDefinitief)
	Energieklasse.Store(i.Energieklasse)
	Gemeentecode.Store(i.Gemeentecode)
	Gemeentenaam.Store(i.Gemeentenaam)
	Buurtcode.Store(i.Buurtcode)
	Buurtnaam.Store(i.Buurtnaam)
	Wijkcode.Store(i.Wijkcode)
	Wijknaam.Store(i.Wijknaam)
	Provinciecode.Store(i.Provinciecode)
	Provincienaam.Store(i.Provincienaam)
	PandGasEanAansluitingen.Store(i.PandGasEanAansluitingen)
	P6GasAansluitingen2020.Store(i.P6GasAansluitingen2020)
	P6Gasm32020.Store(i.P6Gasm32020)
	P6Kwh2020.Store(i.P6Kwh2020)
	P6TotaalPandoppervlakM2.Store(i.P6TotaalPandoppervlakM2)
	PandBouwjaar.Store(i.PandBouwjaar)
	PandGasAansluitingen.Store(i.PandGasAansluitingen)

	doelen := Gebruiksdoelen.StoreArray(i.Gebruiksdoelen)

	return Item{

		label,

		Pid.GetIndex(i.Pid),
		Vid.GetIndex(i.Vid),
		i.Numid,
		Straat.GetIndex(i.Straat),
		Postcode.GetIndex(i.Postcode),
		Huisnummer.GetIndex(i.Huisnummer),
		Huisletter.GetIndex(i.Huisletter),
		Huisnummertoevoeging.GetIndex(i.Huisnummertoevoeging),
		Oppervlakte.GetIndex(i.Oppervlakte),
		Woningequivalent.GetIndex(i.Woningequivalent),
		WoningType.GetIndex(i.WoningType),
		LabelscoreVoorlopig.GetIndex(i.LabelscoreVoorlopig),
		LabelscoreDefinitief.GetIndex(i.LabelscoreDefinitief),
		Energieklasse.GetIndex(i.Energieklasse),
		Gemeentecode.GetIndex(i.Gemeentecode),
		Gemeentenaam.GetIndex(i.Gemeentenaam),
		Buurtcode.GetIndex(i.Buurtcode),
		Buurtnaam.GetIndex(i.Buurtnaam),
		Wijkcode.GetIndex(i.Wijkcode),
		Wijknaam.GetIndex(i.Wijknaam),
		Provinciecode.GetIndex(i.Provinciecode),
		Provincienaam.GetIndex(i.Provincienaam),
		i.Point,
		PandGasEanAansluitingen.GetIndex(i.PandGasEanAansluitingen),
		i.GroupId2020,
		P6GasAansluitingen2020.GetIndex(i.P6GasAansluitingen2020),
		P6Gasm32020.GetIndex(i.P6Gasm32020),
		P6Kwh2020.GetIndex(i.P6Kwh2020),
		P6TotaalPandoppervlakM2.GetIndex(i.P6TotaalPandoppervlakM2),
		PandBouwjaar.GetIndex(i.PandBouwjaar),
		PandGasAansluitingen.GetIndex(i.PandGasAansluitingen),
		doelen,
	}
}

// Store selected columns in seperate map[columnvalue]bitarray
// for fast item selection
func (i *Item) StoreBitArrayColumns() {

	SetBitArray("pid", i.Pid, i.Label)
	SetBitArray("postcode", i.Postcode, i.Label)
	SetBitArray("woning_type", i.WoningType, i.Label)
	SetBitArray("labelscore_voorlopig", i.LabelscoreVoorlopig, i.Label)
	SetBitArray("labelscore_definitief", i.LabelscoreDefinitief, i.Label)
	SetBitArray("energieklasse", i.Energieklasse, i.Label)
	SetBitArray("gemeentecode", i.Gemeentecode, i.Label)
	SetBitArray("buurtcode", i.Buurtcode, i.Label)
	SetBitArray("wijkcode", i.Wijkcode, i.Label)
	SetBitArray("provinciecode", i.Provinciecode, i.Label)

}

func (i Item) Serialize() ItemOut {
	return ItemOut{

		Pid.GetValue(i.Pid),
		Vid.GetValue(i.Vid),
		i.Numid,
		Straat.GetValue(i.Straat),
		Postcode.GetValue(i.Postcode),
		Huisnummer.GetValue(i.Huisnummer),
		Huisletter.GetValue(i.Huisletter),
		Huisnummertoevoeging.GetValue(i.Huisnummertoevoeging),
		GetAdres(&i),
		Oppervlakte.GetValue(i.Oppervlakte),
		Woningequivalent.GetValue(i.Woningequivalent),
		WoningType.GetValue(i.WoningType),
		LabelscoreVoorlopig.GetValue(i.LabelscoreVoorlopig),
		LabelscoreDefinitief.GetValue(i.LabelscoreDefinitief),
		Energieklasse.GetValue(i.Energieklasse),
		Gemeentecode.GetValue(i.Gemeentecode),
		Gemeentenaam.GetValue(i.Gemeentenaam),
		Buurtcode.GetValue(i.Buurtcode),
		Buurtnaam.GetValue(i.Buurtnaam),
		Wijkcode.GetValue(i.Wijkcode),
		Wijknaam.GetValue(i.Wijknaam),
		Provinciecode.GetValue(i.Provinciecode),
		Provincienaam.GetValue(i.Provincienaam),
		i.Point,
		PandGasEanAansluitingen.GetValue(i.PandGasEanAansluitingen),
		i.GroupId2020,
		P6GasAansluitingen2020.GetValue(i.P6GasAansluitingen2020),
		P6Gasm32020.GetValue(i.P6Gasm32020),
		P6Kwh2020.GetValue(i.P6Kwh2020),
		P6TotaalPandoppervlakM2.GetValue(i.P6TotaalPandoppervlakM2),
		PandBouwjaar.GetValue(i.PandBouwjaar),
		PandGasAansluitingen.GetValue(i.PandGasAansluitingen),
		Gebruiksdoelen.GetArrayValue(i.Gebruiksdoelen),
	}
}

func (i ItemIn) Columns() []string {
	return []string{

		"pid",
		"vid",
		"numid",
		"straat",
		"postcode",
		"huisnummer",
		"huisletter",
		"huisnummertoevoeging",
		"oppervlakte",
		"woningequivalent",
		"woning_type",
		"labelscore_voorlopig",
		"labelscore_definitief",
		"energieklasse",
		"gemeentecode",
		"gemeentenaam",
		"buurtcode",
		"buurtnaam",
		"wijkcode",
		"wijknaam",
		"provinciecode",
		"provincienaam",
		"point",
		"pand_gas_ean_aansluitingen",
		"group_id_2020",
		"p6_gas_aansluitingen_2020",
		"p6_gasm3_2020",
		"p6_kwh_2020",
		"p6_totaal_pandoppervlak_m2",
		"pand_bouwjaar",
		"pand_gas_aansluitingen",
		"gebruiksdoelen",
	}
}

func (i ItemOut) Columns() []string {
	return []string{

		"pid",
		"vid",
		"numid",
		"straat",
		"postcode",
		"huisnummer",
		"huisletter",
		"huisnummertoevoeging",
		"oppervlakte",
		"woningequivalent",
		"woning_type",
		"labelscore_voorlopig",
		"labelscore_definitief",
		"energieklasse",
		"gemeentecode",
		"gemeentenaam",
		"buurtcode",
		"buurtnaam",
		"wijkcode",
		"wijknaam",
		"provinciecode",
		"provincienaam",
		"point",
		"pand_gas_ean_aansluitingen",
		"group_id_2020",
		"p6_gas_aansluitingen_2020",
		"p6_gasm3_2020",
		"p6_kwh_2020",
		"p6_totaal_pandoppervlak_m2",
		"pand_bouwjaar",
		"pand_gas_aansluitingen",
		"gebruiksdoelen",
	}
}

func (i Item) Row() []string {

	return []string{

		Pid.GetValue(i.Pid),
		Vid.GetValue(i.Vid),
		i.Numid,
		Straat.GetValue(i.Straat),
		Postcode.GetValue(i.Postcode),
		Huisnummer.GetValue(i.Huisnummer),
		Huisletter.GetValue(i.Huisletter),
		Huisnummertoevoeging.GetValue(i.Huisnummertoevoeging),
		Oppervlakte.GetValue(i.Oppervlakte),
		Woningequivalent.GetValue(i.Woningequivalent),
		WoningType.GetValue(i.WoningType),
		LabelscoreVoorlopig.GetValue(i.LabelscoreVoorlopig),
		LabelscoreDefinitief.GetValue(i.LabelscoreDefinitief),
		Energieklasse.GetValue(i.Energieklasse),
		Gemeentecode.GetValue(i.Gemeentecode),
		Gemeentenaam.GetValue(i.Gemeentenaam),
		Buurtcode.GetValue(i.Buurtcode),
		Buurtnaam.GetValue(i.Buurtnaam),
		Wijkcode.GetValue(i.Wijkcode),
		Wijknaam.GetValue(i.Wijknaam),
		Provinciecode.GetValue(i.Provinciecode),
		Provincienaam.GetValue(i.Provincienaam),
		i.Point,
		PandGasEanAansluitingen.GetValue(i.PandGasEanAansluitingen),
		i.GroupId2020,
		P6GasAansluitingen2020.GetValue(i.P6GasAansluitingen2020),
		P6Gasm32020.GetValue(i.P6Gasm32020),
		P6Kwh2020.GetValue(i.P6Kwh2020),
		P6TotaalPandoppervlakM2.GetValue(i.P6TotaalPandoppervlakM2),
		PandBouwjaar.GetValue(i.PandBouwjaar),
		PandGasAansluitingen.GetValue(i.PandGasAansluitingen),
		Gebruiksdoelen.GetArrayValue(i.Gebruiksdoelen),
	}
}

func (i Item) GetIndex() string {
	return GettersPid(&i)
}

func (i Item) GetGeometry() string {
	return GettersPoint(&i)
}

// contain filter Pid
func FilterPidContains(i *Item, s string) bool {
	return strings.Contains(Pid.GetValue(i.Pid), s)
}

// startswith filter Pid
func FilterPidStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Pid.GetValue(i.Pid), s)
}

// match filters Pid
func FilterPidMatch(i *Item, s string) bool {
	return Pid.GetValue(i.Pid) == s
}

// getter Pid
func GettersPid(i *Item) string {
	return Pid.GetValue(i.Pid)
}

// contain filter Vid
func FilterVidContains(i *Item, s string) bool {
	return strings.Contains(Vid.GetValue(i.Vid), s)
}

// startswith filter Vid
func FilterVidStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Vid.GetValue(i.Vid), s)
}

// match filters Vid
func FilterVidMatch(i *Item, s string) bool {
	return Vid.GetValue(i.Vid) == s
}

// getter Vid
func GettersVid(i *Item) string {
	return Vid.GetValue(i.Vid)
}

// contain filter Numid
func FilterNumidContains(i *Item, s string) bool {
	return strings.Contains(i.Numid, s)
}

// startswith filter Numid
func FilterNumidStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Numid, s)
}

// match filters Numid
func FilterNumidMatch(i *Item, s string) bool {
	return i.Numid == s
}

// getter Numid
func GettersNumid(i *Item) string {
	return i.Numid
}

// contain filter Straat
func FilterStraatContains(i *Item, s string) bool {
	return strings.Contains(Straat.GetValue(i.Straat), s)
}

// startswith filter Straat
func FilterStraatStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Straat.GetValue(i.Straat), s)
}

// match filters Straat
func FilterStraatMatch(i *Item, s string) bool {
	return Straat.GetValue(i.Straat) == s
}

// getter Straat
func GettersStraat(i *Item) string {
	return Straat.GetValue(i.Straat)
}

// contain filter Postcode
func FilterPostcodeContains(i *Item, s string) bool {
	return strings.Contains(Postcode.GetValue(i.Postcode), s)
}

// startswith filter Postcode
func FilterPostcodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Postcode.GetValue(i.Postcode), s)
}

// match filters Postcode
func FilterPostcodeMatch(i *Item, s string) bool {
	return Postcode.GetValue(i.Postcode) == s
}

// getter Postcode
func GettersPostcode(i *Item) string {
	return Postcode.GetValue(i.Postcode)
}

// contain filter Huisnummer
func FilterHuisnummerContains(i *Item, s string) bool {
	return strings.Contains(Huisnummer.GetValue(i.Huisnummer), s)
}

// startswith filter Huisnummer
func FilterHuisnummerStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Huisnummer.GetValue(i.Huisnummer), s)
}

// match filters Huisnummer
func FilterHuisnummerMatch(i *Item, s string) bool {
	return Huisnummer.GetValue(i.Huisnummer) == s
}

// getter Huisnummer
func GettersHuisnummer(i *Item) string {
	return Huisnummer.GetValue(i.Huisnummer)
}

// contain filter Huisletter
func FilterHuisletterContains(i *Item, s string) bool {
	return strings.Contains(Huisletter.GetValue(i.Huisletter), s)
}

// startswith filter Huisletter
func FilterHuisletterStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Huisletter.GetValue(i.Huisletter), s)
}

// match filters Huisletter
func FilterHuisletterMatch(i *Item, s string) bool {
	return Huisletter.GetValue(i.Huisletter) == s
}

// getter Huisletter
func GettersHuisletter(i *Item) string {
	return Huisletter.GetValue(i.Huisletter)
}

// contain filter Huisnummertoevoeging
func FilterHuisnummertoevoegingContains(i *Item, s string) bool {
	return strings.Contains(Huisnummertoevoeging.GetValue(i.Huisnummertoevoeging), s)
}

// startswith filter Huisnummertoevoeging
func FilterHuisnummertoevoegingStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Huisnummertoevoeging.GetValue(i.Huisnummertoevoeging), s)
}

// match filters Huisnummertoevoeging
func FilterHuisnummertoevoegingMatch(i *Item, s string) bool {
	return Huisnummertoevoeging.GetValue(i.Huisnummertoevoeging) == s
}

// getter Huisnummertoevoeging
func GettersHuisnummertoevoeging(i *Item) string {
	return Huisnummertoevoeging.GetValue(i.Huisnummertoevoeging)
}

// contain filter Oppervlakte
func FilterOppervlakteContains(i *Item, s string) bool {
	return strings.Contains(Oppervlakte.GetValue(i.Oppervlakte), s)
}

// startswith filter Oppervlakte
func FilterOppervlakteStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Oppervlakte.GetValue(i.Oppervlakte), s)
}

// match filters Oppervlakte
func FilterOppervlakteMatch(i *Item, s string) bool {
	return Oppervlakte.GetValue(i.Oppervlakte) == s
}

// getter Oppervlakte
func GettersOppervlakte(i *Item) string {
	return Oppervlakte.GetValue(i.Oppervlakte)
}

// contain filter Woningequivalent
func FilterWoningequivalentContains(i *Item, s string) bool {
	return strings.Contains(Woningequivalent.GetValue(i.Woningequivalent), s)
}

// startswith filter Woningequivalent
func FilterWoningequivalentStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Woningequivalent.GetValue(i.Woningequivalent), s)
}

// match filters Woningequivalent
func FilterWoningequivalentMatch(i *Item, s string) bool {
	return Woningequivalent.GetValue(i.Woningequivalent) == s
}

// getter Woningequivalent
func GettersWoningequivalent(i *Item) string {
	return Woningequivalent.GetValue(i.Woningequivalent)
}

/*
// contain filter Adres
func FilterAdresContains(i *Item, s string) bool {
	return strings.Contains(i.Adres, s)
}

// startswith filter Adres
func FilterAdresStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Adres, s)
}

// match filters Adres
func FilterAdresMatch(i *Item, s string) bool {
	return i.Adres == s
}

// getter Adres
func GettersAdres(i *Item) string {
	return i.Adres
}
*/

// contain filter WoningType
func FilterWoningTypeContains(i *Item, s string) bool {
	return strings.Contains(WoningType.GetValue(i.WoningType), s)
}

// startswith filter WoningType
func FilterWoningTypeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(WoningType.GetValue(i.WoningType), s)
}

// match filters WoningType
func FilterWoningTypeMatch(i *Item, s string) bool {
	return WoningType.GetValue(i.WoningType) == s
}

// getter WoningType
func GettersWoningType(i *Item) string {
	return WoningType.GetValue(i.WoningType)
}

// contain filter LabelscoreVoorlopig
func FilterLabelscoreVoorlopigContains(i *Item, s string) bool {
	return strings.Contains(LabelscoreVoorlopig.GetValue(i.LabelscoreVoorlopig), s)
}

// startswith filter LabelscoreVoorlopig
func FilterLabelscoreVoorlopigStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(LabelscoreVoorlopig.GetValue(i.LabelscoreVoorlopig), s)
}

// match filters LabelscoreVoorlopig
func FilterLabelscoreVoorlopigMatch(i *Item, s string) bool {
	return LabelscoreVoorlopig.GetValue(i.LabelscoreVoorlopig) == s
}

// getter LabelscoreVoorlopig
func GettersLabelscoreVoorlopig(i *Item) string {
	return LabelscoreVoorlopig.GetValue(i.LabelscoreVoorlopig)
}

// contain filter LabelscoreDefinitief
func FilterLabelscoreDefinitiefContains(i *Item, s string) bool {
	return strings.Contains(LabelscoreDefinitief.GetValue(i.LabelscoreDefinitief), s)
}

// startswith filter LabelscoreDefinitief
func FilterLabelscoreDefinitiefStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(LabelscoreDefinitief.GetValue(i.LabelscoreDefinitief), s)
}

// match filters LabelscoreDefinitief
func FilterLabelscoreDefinitiefMatch(i *Item, s string) bool {
	return LabelscoreDefinitief.GetValue(i.LabelscoreDefinitief) == s
}

// getter LabelscoreDefinitief
func GettersLabelscoreDefinitief(i *Item) string {
	return LabelscoreDefinitief.GetValue(i.LabelscoreDefinitief)
}

// contain filter Energieklasse
func FilterEnergieklasseContains(i *Item, s string) bool {
	return strings.Contains(Energieklasse.GetValue(i.Energieklasse), s)
}

// startswith filter Energieklasse
func FilterEnergieklasseStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Energieklasse.GetValue(i.Energieklasse), s)
}

// match filters Energieklasse
func FilterEnergieklasseMatch(i *Item, s string) bool {
	return Energieklasse.GetValue(i.Energieklasse) == s
}

// getter Energieklasse
func GettersEnergieklasse(i *Item) string {
	return Energieklasse.GetValue(i.Energieklasse)
}

// contain filter Gemeentecode
func FilterGemeentecodeContains(i *Item, s string) bool {
	return strings.Contains(Gemeentecode.GetValue(i.Gemeentecode), s)
}

// startswith filter Gemeentecode
func FilterGemeentecodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Gemeentecode.GetValue(i.Gemeentecode), s)
}

// match filters Gemeentecode
func FilterGemeentecodeMatch(i *Item, s string) bool {
	return Gemeentecode.GetValue(i.Gemeentecode) == s
}

// getter Gemeentecode
func GettersGemeentecode(i *Item) string {
	return Gemeentecode.GetValue(i.Gemeentecode)
}

// contain filter Gemeentenaam
func FilterGemeentenaamContains(i *Item, s string) bool {
	return strings.Contains(Gemeentenaam.GetValue(i.Gemeentenaam), s)
}

// startswith filter Gemeentenaam
func FilterGemeentenaamStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Gemeentenaam.GetValue(i.Gemeentenaam), s)
}

// match filters Gemeentenaam
func FilterGemeentenaamMatch(i *Item, s string) bool {
	return Gemeentenaam.GetValue(i.Gemeentenaam) == s
}

// getter Gemeentenaam
func GettersGemeentenaam(i *Item) string {
	return Gemeentenaam.GetValue(i.Gemeentenaam)
}

// contain filter Buurtcode
func FilterBuurtcodeContains(i *Item, s string) bool {
	return strings.Contains(Buurtcode.GetValue(i.Buurtcode), s)
}

// startswith filter Buurtcode
func FilterBuurtcodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Buurtcode.GetValue(i.Buurtcode), s)
}

// match filters Buurtcode
func FilterBuurtcodeMatch(i *Item, s string) bool {
	return Buurtcode.GetValue(i.Buurtcode) == s
}

// getter Buurtcode
func GettersBuurtcode(i *Item) string {
	return Buurtcode.GetValue(i.Buurtcode)
}

// contain filter Buurtnaam
func FilterBuurtnaamContains(i *Item, s string) bool {
	return strings.Contains(Buurtnaam.GetValue(i.Buurtnaam), s)
}

// startswith filter Buurtnaam
func FilterBuurtnaamStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Buurtnaam.GetValue(i.Buurtnaam), s)
}

// match filters Buurtnaam
func FilterBuurtnaamMatch(i *Item, s string) bool {
	return Buurtnaam.GetValue(i.Buurtnaam) == s
}

// getter Buurtnaam
func GettersBuurtnaam(i *Item) string {
	return Buurtnaam.GetValue(i.Buurtnaam)
}

// contain filter Wijkcode
func FilterWijkcodeContains(i *Item, s string) bool {
	return strings.Contains(Wijkcode.GetValue(i.Wijkcode), s)
}

// startswith filter Wijkcode
func FilterWijkcodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Wijkcode.GetValue(i.Wijkcode), s)
}

// match filters Wijkcode
func FilterWijkcodeMatch(i *Item, s string) bool {
	return Wijkcode.GetValue(i.Wijkcode) == s
}

// getter Wijkcode
func GettersWijkcode(i *Item) string {
	return Wijkcode.GetValue(i.Wijkcode)
}

// contain filter Wijknaam
func FilterWijknaamContains(i *Item, s string) bool {
	return strings.Contains(Wijknaam.GetValue(i.Wijknaam), s)
}

// startswith filter Wijknaam
func FilterWijknaamStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Wijknaam.GetValue(i.Wijknaam), s)
}

// match filters Wijknaam
func FilterWijknaamMatch(i *Item, s string) bool {
	return Wijknaam.GetValue(i.Wijknaam) == s
}

// getter Wijknaam
func GettersWijknaam(i *Item) string {
	return Wijknaam.GetValue(i.Wijknaam)
}

// contain filter Provinciecode
func FilterProvinciecodeContains(i *Item, s string) bool {
	return strings.Contains(Provinciecode.GetValue(i.Provinciecode), s)
}

// startswith filter Provinciecode
func FilterProvinciecodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Provinciecode.GetValue(i.Provinciecode), s)
}

// match filters Provinciecode
func FilterProvinciecodeMatch(i *Item, s string) bool {
	return Provinciecode.GetValue(i.Provinciecode) == s
}

// getter Provinciecode
func GettersProvinciecode(i *Item) string {
	return Provinciecode.GetValue(i.Provinciecode)
}

// contain filter Provincienaam
func FilterProvincienaamContains(i *Item, s string) bool {
	return strings.Contains(Provincienaam.GetValue(i.Provincienaam), s)
}

// startswith filter Provincienaam
func FilterProvincienaamStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Provincienaam.GetValue(i.Provincienaam), s)
}

// match filters Provincienaam
func FilterProvincienaamMatch(i *Item, s string) bool {
	return Provincienaam.GetValue(i.Provincienaam) == s
}

// getter Provincienaam
func GettersProvincienaam(i *Item) string {
	return Provincienaam.GetValue(i.Provincienaam)
}

// contain filter Point
func FilterPointContains(i *Item, s string) bool {
	return strings.Contains(i.Point, s)
}

// startswith filter Point
func FilterPointStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Point, s)
}

// match filters Point
func FilterPointMatch(i *Item, s string) bool {
	return i.Point == s
}

// getter Point
func GettersPoint(i *Item) string {
	return i.Point
}

// contain filter PandGasEanAansluitingen
func FilterPandGasEanAansluitingenContains(i *Item, s string) bool {
	return strings.Contains(PandGasEanAansluitingen.GetValue(i.PandGasEanAansluitingen), s)
}

// startswith filter PandGasEanAansluitingen
func FilterPandGasEanAansluitingenStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(PandGasEanAansluitingen.GetValue(i.PandGasEanAansluitingen), s)
}

// match filters PandGasEanAansluitingen
func FilterPandGasEanAansluitingenMatch(i *Item, s string) bool {
	return PandGasEanAansluitingen.GetValue(i.PandGasEanAansluitingen) == s
}

// getter PandGasEanAansluitingen
func GettersPandGasEanAansluitingen(i *Item) string {
	return PandGasEanAansluitingen.GetValue(i.PandGasEanAansluitingen)
}

// contain filter GroupId2020
func FilterGroupId2020Contains(i *Item, s string) bool {
	return strings.Contains(i.GroupId2020, s)
}

// startswith filter GroupId2020
func FilterGroupId2020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.GroupId2020, s)
}

// match filters GroupId2020
func FilterGroupId2020Match(i *Item, s string) bool {
	return i.GroupId2020 == s
}

// getter GroupId2020
func GettersGroupId2020(i *Item) string {
	return i.GroupId2020
}

// contain filter P6GasAansluitingen2020
func FilterP6GasAansluitingen2020Contains(i *Item, s string) bool {
	return strings.Contains(P6GasAansluitingen2020.GetValue(i.P6GasAansluitingen2020), s)
}

// startswith filter P6GasAansluitingen2020
func FilterP6GasAansluitingen2020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(P6GasAansluitingen2020.GetValue(i.P6GasAansluitingen2020), s)
}

// match filters P6GasAansluitingen2020
func FilterP6GasAansluitingen2020Match(i *Item, s string) bool {
	return P6GasAansluitingen2020.GetValue(i.P6GasAansluitingen2020) == s
}

// getter P6GasAansluitingen2020
func GettersP6GasAansluitingen2020(i *Item) string {
	return P6GasAansluitingen2020.GetValue(i.P6GasAansluitingen2020)
}

// contain filter P6Gasm32020
func FilterP6Gasm32020Contains(i *Item, s string) bool {
	return strings.Contains(P6Gasm32020.GetValue(i.P6Gasm32020), s)
}

// startswith filter P6Gasm32020
func FilterP6Gasm32020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(P6Gasm32020.GetValue(i.P6Gasm32020), s)
}

// match filters P6Gasm32020
func FilterP6Gasm32020Match(i *Item, s string) bool {
	return P6Gasm32020.GetValue(i.P6Gasm32020) == s
}

// getter P6Gasm32020
func GettersP6Gasm32020(i *Item) string {
	return P6Gasm32020.GetValue(i.P6Gasm32020)
}

// contain filter P6Kwh2020
func FilterP6Kwh2020Contains(i *Item, s string) bool {
	return strings.Contains(P6Kwh2020.GetValue(i.P6Kwh2020), s)
}

// startswith filter P6Kwh2020
func FilterP6Kwh2020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(P6Kwh2020.GetValue(i.P6Kwh2020), s)
}

// match filters P6Kwh2020
func FilterP6Kwh2020Match(i *Item, s string) bool {
	return P6Kwh2020.GetValue(i.P6Kwh2020) == s
}

// getter P6Kwh2020
func GettersP6Kwh2020(i *Item) string {
	return P6Kwh2020.GetValue(i.P6Kwh2020)
}

// contain filter P6TotaalPandoppervlakM2
func FilterP6TotaalPandoppervlakM2Contains(i *Item, s string) bool {
	return strings.Contains(P6TotaalPandoppervlakM2.GetValue(i.P6TotaalPandoppervlakM2), s)
}

// startswith filter P6TotaalPandoppervlakM2
func FilterP6TotaalPandoppervlakM2StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(P6TotaalPandoppervlakM2.GetValue(i.P6TotaalPandoppervlakM2), s)
}

// match filters P6TotaalPandoppervlakM2
func FilterP6TotaalPandoppervlakM2Match(i *Item, s string) bool {
	return P6TotaalPandoppervlakM2.GetValue(i.P6TotaalPandoppervlakM2) == s
}

// getter P6TotaalPandoppervlakM2
func GettersP6TotaalPandoppervlakM2(i *Item) string {
	return P6TotaalPandoppervlakM2.GetValue(i.P6TotaalPandoppervlakM2)
}

// contain filter PandBouwjaar
func FilterPandBouwjaarContains(i *Item, s string) bool {
	return strings.Contains(PandBouwjaar.GetValue(i.PandBouwjaar), s)
}

// startswith filter PandBouwjaar
func FilterPandBouwjaarStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(PandBouwjaar.GetValue(i.PandBouwjaar), s)
}

// match filters PandBouwjaar
func FilterPandBouwjaarMatch(i *Item, s string) bool {
	return PandBouwjaar.GetValue(i.PandBouwjaar) == s
}

// getter PandBouwjaar
func GettersPandBouwjaar(i *Item) string {
	return PandBouwjaar.GetValue(i.PandBouwjaar)
}

// contain filter PandGasAansluitingen
func FilterPandGasAansluitingenContains(i *Item, s string) bool {
	return strings.Contains(PandGasAansluitingen.GetValue(i.PandGasAansluitingen), s)
}

// startswith filter PandGasAansluitingen
func FilterPandGasAansluitingenStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(PandGasAansluitingen.GetValue(i.PandGasAansluitingen), s)
}

// match filters PandGasAansluitingen
func FilterPandGasAansluitingenMatch(i *Item, s string) bool {
	return PandGasAansluitingen.GetValue(i.PandGasAansluitingen) == s
}

// getter PandGasAansluitingen
func GettersPandGasAansluitingen(i *Item) string {
	return PandGasAansluitingen.GetValue(i.PandGasAansluitingen)
}

// contain filter Gebruiksdoelen
func FilterGebruiksdoelenContains(i *Item, s string) bool {
	for _, v := range i.Gebruiksdoelen {
		vs := Gebruiksdoelen.GetValue(v)
		if strings.Contains(vs, s) {
			return true
		}
	}
	return false
}

// startswith filter Gebruiksdoelen
func FilterGebruiksdoelenStartsWith(i *Item, s string) bool {
	for _, v := range i.Gebruiksdoelen {
		vs := Gebruiksdoelen.GetValue(v)
		if strings.HasPrefix(vs, s) {
			return true
		}
	}
	return false

}

// match filters Gebruiksdoelen
func FilterGebruiksdoelenMatch(i *Item, s string) bool {
	for _, v := range i.Gebruiksdoelen {
		vs := Gebruiksdoelen.GetValue(v)
		if vs == s {
			return true
		}
	}
	return false
}

// getter Gebruiksdoelen
func GettersGebruiksdoelen(i *Item) string {
	return Gebruiksdoelen.GetArrayValue(i.Gebruiksdoelen)
}

/*
// contain filters
func FilterEkeyContains(i *Item, s string) bool {
	return strings.Contains(i.Ekey, s)
}


// startswith filters
func FilterEkeyStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Ekey, s)
}


// match filters
func FilterEkeyMatch(i *Item, s string) bool {
	return i.Ekey == s
}

// getters
func GettersEkey(i *Item) string {
	return i.Ekey
}
*/

// reduce functions
func reduceCount(items Items) map[string]string {
	result := make(map[string]string)
	result["count"] = strconv.Itoa(len(items))
	return result
}

type GroupedOperations struct {
	Funcs     registerFuncType
	GroupBy   registerGroupByFunc
	Getters   registerGettersMap
	Reduce    registerReduce
	BitArrays registerBitArray
}

var Operations GroupedOperations

var RegisterFuncMap registerFuncType
var RegisterGroupBy registerGroupByFunc
var RegisterGetters registerGettersMap
var RegisterReduce registerReduce
var RegisterBitArray registerBitArray

// ValidateRegsiters validate exposed columns do match filter names
func validateRegisters() error {
	var i = ItemOut{}
	var filters = []string{"match", "contains", "startswith"}
	for _, c := range i.Columns() {
		for _, f := range filters {
			if _, ok := RegisterFuncMap[f+"-"+c]; !ok {
				return errors.New(c + " is missing in RegisterMap")
			}
		}
	}
	return nil
}

func init() {

	RegisterFuncMap = make(registerFuncType)
	RegisterGroupBy = make(registerGroupByFunc)
	RegisterGetters = make(registerGettersMap)
	RegisterReduce = make(registerReduce)

	// register search filter.
	//RegisterFuncMap["search"] = 'EDITYOURSELF'
	// example RegisterFuncMap["search"] = FilterEkeyStartsWith

	//RegisterFuncMap["value"] = 'EDITYOURSELF'
	RegisterGetters["value"] = GettersGemeentecode

	// register filters

	//register filters for Pid
	RegisterFuncMap["match-pid"] = FilterPidMatch
	RegisterFuncMap["contains-pid"] = FilterPidContains
	RegisterFuncMap["startswith-pid"] = FilterPidStartsWith
	RegisterGetters["pid"] = GettersPid
	RegisterGroupBy["pid"] = GettersPid

	//register filters for Vid
	RegisterFuncMap["match-vid"] = FilterVidMatch
	RegisterFuncMap["contains-vid"] = FilterVidContains
	RegisterFuncMap["startswith-vid"] = FilterVidStartsWith
	RegisterGetters["vid"] = GettersVid
	RegisterGroupBy["vid"] = GettersVid

	//register filters for Numid
	RegisterFuncMap["match-numid"] = FilterNumidMatch
	RegisterFuncMap["contains-numid"] = FilterNumidContains
	RegisterFuncMap["startswith-numid"] = FilterNumidStartsWith
	RegisterGetters["numid"] = GettersNumid
	RegisterGroupBy["numid"] = GettersNumid

	//register filters for Straat
	RegisterFuncMap["match-straat"] = FilterStraatMatch
	RegisterFuncMap["contains-straat"] = FilterStraatContains
	RegisterFuncMap["startswith-straat"] = FilterStraatStartsWith
	RegisterGetters["straat"] = GettersStraat
	RegisterGroupBy["straat"] = GettersStraat

	//register filters for Postcode
	RegisterFuncMap["match-postcode"] = FilterPostcodeMatch
	RegisterFuncMap["contains-postcode"] = FilterPostcodeContains
	RegisterFuncMap["startswith-postcode"] = FilterPostcodeStartsWith
	RegisterGetters["postcode"] = GettersPostcode
	RegisterGroupBy["postcode"] = GettersPostcode

	//register filters for Huisnummer
	RegisterFuncMap["match-huisnummer"] = FilterHuisnummerMatch
	RegisterFuncMap["contains-huisnummer"] = FilterHuisnummerContains
	RegisterFuncMap["startswith-huisnummer"] = FilterHuisnummerStartsWith
	RegisterGetters["huisnummer"] = GettersHuisnummer
	RegisterGroupBy["huisnummer"] = GettersHuisnummer

	//register filters for Huisletter
	RegisterFuncMap["match-huisletter"] = FilterHuisletterMatch
	RegisterFuncMap["contains-huisletter"] = FilterHuisletterContains
	RegisterFuncMap["startswith-huisletter"] = FilterHuisletterStartsWith
	RegisterGetters["huisletter"] = GettersHuisletter
	RegisterGroupBy["huisletter"] = GettersHuisletter

	//register filters for Huisnummertoevoeging
	RegisterFuncMap["match-huisnummertoevoeging"] = FilterHuisnummertoevoegingMatch
	RegisterFuncMap["contains-huisnummertoevoeging"] = FilterHuisnummertoevoegingContains
	RegisterFuncMap["startswith-huisnummertoevoeging"] = FilterHuisnummertoevoegingStartsWith
	RegisterGetters["huisnummertoevoeging"] = GettersHuisnummertoevoeging
	RegisterGroupBy["huisnummertoevoeging"] = GettersHuisnummertoevoeging

	//register filters for Oppervlakte
	RegisterFuncMap["match-oppervlakte"] = FilterOppervlakteMatch
	RegisterFuncMap["contains-oppervlakte"] = FilterOppervlakteContains
	RegisterFuncMap["startswith-oppervlakte"] = FilterOppervlakteStartsWith
	RegisterGetters["oppervlakte"] = GettersOppervlakte
	RegisterGroupBy["oppervlakte"] = GettersOppervlakte

	//register filters for Woningequivalent
	RegisterFuncMap["match-woningequivalent"] = FilterWoningequivalentMatch
	RegisterFuncMap["contains-woningequivalent"] = FilterWoningequivalentContains
	RegisterFuncMap["startswith-woningequivalent"] = FilterWoningequivalentStartsWith
	RegisterGetters["woningequivalent"] = GettersWoningequivalent
	RegisterGroupBy["woningequivalent"] = GettersWoningequivalent
	//register filters for WoningType
	RegisterFuncMap["match-woning_type"] = FilterWoningTypeMatch
	RegisterFuncMap["contains-woning_type"] = FilterWoningTypeContains
	RegisterFuncMap["startswith-woning_type"] = FilterWoningTypeStartsWith
	RegisterGetters["woning_type"] = GettersWoningType
	RegisterGroupBy["woning_type"] = GettersWoningType

	//register filters for LabelscoreVoorlopig
	RegisterFuncMap["match-labelscore_voorlopig"] = FilterLabelscoreVoorlopigMatch
	RegisterFuncMap["contains-labelscore_voorlopig"] = FilterLabelscoreVoorlopigContains
	RegisterFuncMap["startswith-labelscore_voorlopig"] = FilterLabelscoreVoorlopigStartsWith
	RegisterGetters["labelscore_voorlopig"] = GettersLabelscoreVoorlopig
	RegisterGroupBy["labelscore_voorlopig"] = GettersLabelscoreVoorlopig

	//register filters for LabelscoreDefinitief
	RegisterFuncMap["match-labelscore_definitief"] = FilterLabelscoreDefinitiefMatch
	RegisterFuncMap["contains-labelscore_definitief"] = FilterLabelscoreDefinitiefContains
	RegisterFuncMap["startswith-labelscore_definitief"] = FilterLabelscoreDefinitiefStartsWith
	RegisterGetters["labelscore_definitief"] = GettersLabelscoreDefinitief
	RegisterGroupBy["labelscore_definitief"] = GettersLabelscoreDefinitief

	//register filters for Energieklasse
	RegisterFuncMap["match-energieklasse"] = FilterEnergieklasseMatch
	RegisterFuncMap["contains-energieklasse"] = FilterEnergieklasseContains
	RegisterFuncMap["startswith-energieklasse"] = FilterEnergieklasseStartsWith
	RegisterGetters["energieklasse"] = GettersEnergieklasse
	RegisterGroupBy["energieklasse"] = GettersEnergieklasse

	//register filters for Gemeentecode
	RegisterFuncMap["match-gemeentecode"] = FilterGemeentecodeMatch
	RegisterFuncMap["contains-gemeentecode"] = FilterGemeentecodeContains
	RegisterFuncMap["startswith-gemeentecode"] = FilterGemeentecodeStartsWith
	RegisterGetters["gemeentecode"] = GettersGemeentecode
	RegisterGroupBy["gemeentecode"] = GettersGemeentecode

	//register filters for Gemeentenaam
	RegisterFuncMap["match-gemeentenaam"] = FilterGemeentenaamMatch
	RegisterFuncMap["contains-gemeentenaam"] = FilterGemeentenaamContains
	RegisterFuncMap["startswith-gemeentenaam"] = FilterGemeentenaamStartsWith
	RegisterGetters["gemeentenaam"] = GettersGemeentenaam
	RegisterGroupBy["gemeentenaam"] = GettersGemeentenaam

	//register filters for Buurtcode
	RegisterFuncMap["match-buurtcode"] = FilterBuurtcodeMatch
	RegisterFuncMap["contains-buurtcode"] = FilterBuurtcodeContains
	RegisterFuncMap["startswith-buurtcode"] = FilterBuurtcodeStartsWith
	RegisterGetters["buurtcode"] = GettersBuurtcode
	RegisterGroupBy["buurtcode"] = GettersBuurtcode

	//register filters for Buurtnaam
	RegisterFuncMap["match-buurtnaam"] = FilterBuurtnaamMatch
	RegisterFuncMap["contains-buurtnaam"] = FilterBuurtnaamContains
	RegisterFuncMap["startswith-buurtnaam"] = FilterBuurtnaamStartsWith
	RegisterGetters["buurtnaam"] = GettersBuurtnaam
	RegisterGroupBy["buurtnaam"] = GettersBuurtnaam

	//register filters for Wijkcode
	RegisterFuncMap["match-wijkcode"] = FilterWijkcodeMatch
	RegisterFuncMap["contains-wijkcode"] = FilterWijkcodeContains
	RegisterFuncMap["startswith-wijkcode"] = FilterWijkcodeStartsWith
	RegisterGetters["wijkcode"] = GettersWijkcode
	RegisterGroupBy["wijkcode"] = GettersWijkcode

	//register filters for Wijknaam
	RegisterFuncMap["match-wijknaam"] = FilterWijknaamMatch
	RegisterFuncMap["contains-wijknaam"] = FilterWijknaamContains
	RegisterFuncMap["startswith-wijknaam"] = FilterWijknaamStartsWith
	RegisterGetters["wijknaam"] = GettersWijknaam
	RegisterGroupBy["wijknaam"] = GettersWijknaam

	//register filters for Provinciecode
	RegisterFuncMap["match-provinciecode"] = FilterProvinciecodeMatch
	RegisterFuncMap["contains-provinciecode"] = FilterProvinciecodeContains
	RegisterFuncMap["startswith-provinciecode"] = FilterProvinciecodeStartsWith
	RegisterGetters["provinciecode"] = GettersProvinciecode
	RegisterGroupBy["provinciecode"] = GettersProvinciecode

	//register filters for Provincienaam
	RegisterFuncMap["match-provincienaam"] = FilterProvincienaamMatch
	RegisterFuncMap["contains-provincienaam"] = FilterProvincienaamContains
	RegisterFuncMap["startswith-provincienaam"] = FilterProvincienaamStartsWith
	RegisterGetters["provincienaam"] = GettersProvincienaam
	RegisterGroupBy["provincienaam"] = GettersProvincienaam

	//register filters for Point
	RegisterFuncMap["match-point"] = FilterPointMatch
	RegisterFuncMap["contains-point"] = FilterPointContains
	RegisterFuncMap["startswith-point"] = FilterPointStartsWith
	RegisterGetters["point"] = GettersPoint
	RegisterGroupBy["point"] = GettersPoint

	//register filters for PandGasEanAansluitingen
	RegisterFuncMap["match-pand_gas_ean_aansluitingen"] = FilterPandGasEanAansluitingenMatch
	RegisterFuncMap["contains-pand_gas_ean_aansluitingen"] = FilterPandGasEanAansluitingenContains
	RegisterFuncMap["startswith-pand_gas_ean_aansluitingen"] = FilterPandGasEanAansluitingenStartsWith
	RegisterGetters["pand_gas_ean_aansluitingen"] = GettersPandGasEanAansluitingen
	RegisterGroupBy["pand_gas_ean_aansluitingen"] = GettersPandGasEanAansluitingen

	//register filters for GroupId2020
	RegisterFuncMap["match-group_id_2020"] = FilterGroupId2020Match
	RegisterFuncMap["contains-group_id_2020"] = FilterGroupId2020Contains
	RegisterFuncMap["startswith-group_id_2020"] = FilterGroupId2020StartsWith
	RegisterGetters["group_id_2020"] = GettersGroupId2020
	RegisterGroupBy["group_id_2020"] = GettersGroupId2020

	//register filters for P6GasAansluitingen2020
	RegisterFuncMap["match-p6_gas_aansluitingen_2020"] = FilterP6GasAansluitingen2020Match
	RegisterFuncMap["contains-p6_gas_aansluitingen_2020"] = FilterP6GasAansluitingen2020Contains
	RegisterFuncMap["startswith-p6_gas_aansluitingen_2020"] = FilterP6GasAansluitingen2020StartsWith
	RegisterGetters["p6_gas_aansluitingen_2020"] = GettersP6GasAansluitingen2020
	RegisterGroupBy["p6_gas_aansluitingen_2020"] = GettersP6GasAansluitingen2020

	//register filters for P6Gasm32020
	RegisterFuncMap["match-p6_gasm3_2020"] = FilterP6Gasm32020Match
	RegisterFuncMap["contains-p6_gasm3_2020"] = FilterP6Gasm32020Contains
	RegisterFuncMap["startswith-p6_gasm3_2020"] = FilterP6Gasm32020StartsWith
	RegisterGetters["p6_gasm3_2020"] = GettersP6Gasm32020
	RegisterGroupBy["p6_gasm3_2020"] = GettersP6Gasm32020

	//register filters for P6Kwh2020
	RegisterFuncMap["match-p6_kwh_2020"] = FilterP6Kwh2020Match
	RegisterFuncMap["contains-p6_kwh_2020"] = FilterP6Kwh2020Contains
	RegisterFuncMap["startswith-p6_kwh_2020"] = FilterP6Kwh2020StartsWith
	RegisterGetters["p6_kwh_2020"] = GettersP6Kwh2020
	RegisterGroupBy["p6_kwh_2020"] = GettersP6Kwh2020

	//register filters for P6TotaalPandoppervlakM2
	RegisterFuncMap["match-p6_totaal_pandoppervlak_m2"] = FilterP6TotaalPandoppervlakM2Match
	RegisterFuncMap["contains-p6_totaal_pandoppervlak_m2"] = FilterP6TotaalPandoppervlakM2Contains
	RegisterFuncMap["startswith-p6_totaal_pandoppervlak_m2"] = FilterP6TotaalPandoppervlakM2StartsWith
	RegisterGetters["p6_totaal_pandoppervlak_m2"] = GettersP6TotaalPandoppervlakM2
	RegisterGroupBy["p6_totaal_pandoppervlak_m2"] = GettersP6TotaalPandoppervlakM2

	//register filters for PandBouwjaar
	RegisterFuncMap["match-pand_bouwjaar"] = FilterPandBouwjaarMatch
	RegisterFuncMap["contains-pand_bouwjaar"] = FilterPandBouwjaarContains
	RegisterFuncMap["startswith-pand_bouwjaar"] = FilterPandBouwjaarStartsWith
	RegisterGetters["pand_bouwjaar"] = GettersPandBouwjaar
	RegisterGroupBy["pand_bouwjaar"] = GettersPandBouwjaar

	//register filters for PandGasAansluitingen
	RegisterFuncMap["match-pand_gas_aansluitingen"] = FilterPandGasAansluitingenMatch
	RegisterFuncMap["contains-pand_gas_aansluitingen"] = FilterPandGasAansluitingenContains
	RegisterFuncMap["startswith-pand_gas_aansluitingen"] = FilterPandGasAansluitingenStartsWith
	RegisterGetters["pand_gas_aansluitingen"] = GettersPandGasAansluitingen
	RegisterGroupBy["pand_gas_aansluitingen"] = GettersPandGasAansluitingen

	//register filters for Gebruiksdoelen
	RegisterFuncMap["match-gebruiksdoelen"] = FilterGebruiksdoelenMatch
	RegisterFuncMap["contains-gebruiksdoelen"] = FilterGebruiksdoelenContains
	RegisterFuncMap["startswith-gebruiksdoelen"] = FilterGebruiksdoelenStartsWith
	RegisterGetters["gebruiksdoelen"] = GettersGebruiksdoelen
	RegisterGroupBy["gebruiksdoelen"] = GettersGebruiksdoelen

	RegisterGroupBy["postcodehuisnummer"] = GettersToevoegingen

	validateRegisters()

	/*
		RegisterFuncMap["match-ekey"] = FilterEkeyMatch
		RegisterFuncMap["contains-ekey"] = FilterEkeyContains
		// register startswith filters
		RegisterFuncMap["startswith-ekey"] = FilterEkeyStartsWith
		// register getters
		RegisterGetters["ekey"] = GettersEkey
		// register groupby
		RegisterGroupBy["ekey"] = GettersEkey

	*/

	// register reduce functions
	RegisterReduce["count"] = reduceCount
	RegisterReduce["woningequivalent"] = reduceWEQ
}

type sortLookup map[string]func(int, int) bool

func createSort(items Items) sortLookup {

	sortFuncs := sortLookup{

		"pid":  func(i, j int) bool { return Pid.GetValue(items[i].Pid) < Pid.GetValue(items[j].Pid) },
		"-pid": func(i, j int) bool { return Pid.GetValue(items[i].Pid) > Pid.GetValue(items[j].Pid) },

		"vid":  func(i, j int) bool { return Vid.GetValue(items[i].Vid) < Vid.GetValue(items[j].Vid) },
		"-vid": func(i, j int) bool { return Vid.GetValue(items[i].Vid) > Vid.GetValue(items[j].Vid) },

		"numid":  func(i, j int) bool { return items[i].Numid < items[j].Numid },
		"-numid": func(i, j int) bool { return items[i].Numid > items[j].Numid },

		"straat":  func(i, j int) bool { return Straat.GetValue(items[i].Straat) < Straat.GetValue(items[j].Straat) },
		"-straat": func(i, j int) bool { return Straat.GetValue(items[i].Straat) > Straat.GetValue(items[j].Straat) },

		"postcode": func(i, j int) bool {
			return Postcode.GetValue(items[i].Postcode) < Postcode.GetValue(items[j].Postcode)
		},
		"-postcode": func(i, j int) bool {
			return Postcode.GetValue(items[i].Postcode) > Postcode.GetValue(items[j].Postcode)
		},

		"huisnummer": func(i, j int) bool {
			return Huisnummer.GetValue(items[i].Huisnummer) < Huisnummer.GetValue(items[j].Huisnummer)
		},
		"-huisnummer": func(i, j int) bool {
			return Huisnummer.GetValue(items[i].Huisnummer) > Huisnummer.GetValue(items[j].Huisnummer)
		},

		"huisletter": func(i, j int) bool {
			return Huisletter.GetValue(items[i].Huisletter) < Huisletter.GetValue(items[j].Huisletter)
		},
		"-huisletter": func(i, j int) bool {
			return Huisletter.GetValue(items[i].Huisletter) > Huisletter.GetValue(items[j].Huisletter)
		},

		"huisnummertoevoeging": func(i, j int) bool {
			return Huisnummertoevoeging.GetValue(items[i].Huisnummertoevoeging) < Huisnummertoevoeging.GetValue(items[j].Huisnummertoevoeging)
		},
		"-huisnummertoevoeging": func(i, j int) bool {
			return Huisnummertoevoeging.GetValue(items[i].Huisnummertoevoeging) > Huisnummertoevoeging.GetValue(items[j].Huisnummertoevoeging)
		},

		"oppervlakte": func(i, j int) bool {
			return Oppervlakte.GetValue(items[i].Oppervlakte) < Oppervlakte.GetValue(items[j].Oppervlakte)
		},
		"-oppervlakte": func(i, j int) bool {
			return Oppervlakte.GetValue(items[i].Oppervlakte) > Oppervlakte.GetValue(items[j].Oppervlakte)
		},

		"woningequivalent": func(i, j int) bool {
			return Woningequivalent.GetValue(items[i].Woningequivalent) < Woningequivalent.GetValue(items[j].Woningequivalent)
		},
		"-woningequivalent": func(i, j int) bool {
			return Woningequivalent.GetValue(items[i].Woningequivalent) > Woningequivalent.GetValue(items[j].Woningequivalent)
		},

		"woning_type": func(i, j int) bool {
			return WoningType.GetValue(items[i].WoningType) < WoningType.GetValue(items[j].WoningType)
		},
		"-woning_type": func(i, j int) bool {
			return WoningType.GetValue(items[i].WoningType) > WoningType.GetValue(items[j].WoningType)
		},

		"labelscore_voorlopig": func(i, j int) bool {
			return LabelscoreVoorlopig.GetValue(items[i].LabelscoreVoorlopig) < LabelscoreVoorlopig.GetValue(items[j].LabelscoreVoorlopig)
		},
		"-labelscore_voorlopig": func(i, j int) bool {
			return LabelscoreVoorlopig.GetValue(items[i].LabelscoreVoorlopig) > LabelscoreVoorlopig.GetValue(items[j].LabelscoreVoorlopig)
		},

		"labelscore_definitief": func(i, j int) bool {
			return LabelscoreDefinitief.GetValue(items[i].LabelscoreDefinitief) < LabelscoreDefinitief.GetValue(items[j].LabelscoreDefinitief)
		},
		"-labelscore_definitief": func(i, j int) bool {
			return LabelscoreDefinitief.GetValue(items[i].LabelscoreDefinitief) > LabelscoreDefinitief.GetValue(items[j].LabelscoreDefinitief)
		},

		"energieklasse": func(i, j int) bool {
			return Energieklasse.GetValue(items[i].Energieklasse) < Energieklasse.GetValue(items[j].Energieklasse)
		},
		"-energieklasse": func(i, j int) bool {
			return Energieklasse.GetValue(items[i].Energieklasse) > Energieklasse.GetValue(items[j].Energieklasse)
		},

		"gemeentecode": func(i, j int) bool {
			return Gemeentecode.GetValue(items[i].Gemeentecode) < Gemeentecode.GetValue(items[j].Gemeentecode)
		},
		"-gemeentecode": func(i, j int) bool {
			return Gemeentecode.GetValue(items[i].Gemeentecode) > Gemeentecode.GetValue(items[j].Gemeentecode)
		},

		"gemeentenaam": func(i, j int) bool {
			return Gemeentenaam.GetValue(items[i].Gemeentenaam) < Gemeentenaam.GetValue(items[j].Gemeentenaam)
		},
		"-gemeentenaam": func(i, j int) bool {
			return Gemeentenaam.GetValue(items[i].Gemeentenaam) > Gemeentenaam.GetValue(items[j].Gemeentenaam)
		},

		"buurtcode": func(i, j int) bool {
			return Buurtcode.GetValue(items[i].Buurtcode) < Buurtcode.GetValue(items[j].Buurtcode)
		},
		"-buurtcode": func(i, j int) bool {
			return Buurtcode.GetValue(items[i].Buurtcode) > Buurtcode.GetValue(items[j].Buurtcode)
		},

		"buurtnaam": func(i, j int) bool {
			return Buurtnaam.GetValue(items[i].Buurtnaam) < Buurtnaam.GetValue(items[j].Buurtnaam)
		},
		"-buurtnaam": func(i, j int) bool {
			return Buurtnaam.GetValue(items[i].Buurtnaam) > Buurtnaam.GetValue(items[j].Buurtnaam)
		},

		"wijkcode": func(i, j int) bool {
			return Wijkcode.GetValue(items[i].Wijkcode) < Wijkcode.GetValue(items[j].Wijkcode)
		},
		"-wijkcode": func(i, j int) bool {
			return Wijkcode.GetValue(items[i].Wijkcode) > Wijkcode.GetValue(items[j].Wijkcode)
		},

		"wijknaam": func(i, j int) bool {
			return Wijknaam.GetValue(items[i].Wijknaam) < Wijknaam.GetValue(items[j].Wijknaam)
		},
		"-wijknaam": func(i, j int) bool {
			return Wijknaam.GetValue(items[i].Wijknaam) > Wijknaam.GetValue(items[j].Wijknaam)
		},

		"provinciecode": func(i, j int) bool {
			return Provinciecode.GetValue(items[i].Provinciecode) < Provinciecode.GetValue(items[j].Provinciecode)
		},
		"-provinciecode": func(i, j int) bool {
			return Provinciecode.GetValue(items[i].Provinciecode) > Provinciecode.GetValue(items[j].Provinciecode)
		},

		"provincienaam": func(i, j int) bool {
			return Provincienaam.GetValue(items[i].Provincienaam) < Provincienaam.GetValue(items[j].Provincienaam)
		},
		"-provincienaam": func(i, j int) bool {
			return Provincienaam.GetValue(items[i].Provincienaam) > Provincienaam.GetValue(items[j].Provincienaam)
		},

		"point":  func(i, j int) bool { return items[i].Point < items[j].Point },
		"-point": func(i, j int) bool { return items[i].Point > items[j].Point },

		"pand_gas_ean_aansluitingen": func(i, j int) bool {
			return PandGasEanAansluitingen.GetValue(items[i].PandGasEanAansluitingen) < PandGasEanAansluitingen.GetValue(items[j].PandGasEanAansluitingen)
		},
		"-pand_gas_ean_aansluitingen": func(i, j int) bool {
			return PandGasEanAansluitingen.GetValue(items[i].PandGasEanAansluitingen) > PandGasEanAansluitingen.GetValue(items[j].PandGasEanAansluitingen)
		},

		"group_id_2020":  func(i, j int) bool { return items[i].GroupId2020 < items[j].GroupId2020 },
		"-group_id_2020": func(i, j int) bool { return items[i].GroupId2020 > items[j].GroupId2020 },

		"p6_gas_aansluitingen_2020": func(i, j int) bool {
			return P6GasAansluitingen2020.GetValue(items[i].P6GasAansluitingen2020) < P6GasAansluitingen2020.GetValue(items[j].P6GasAansluitingen2020)
		},
		"-p6_gas_aansluitingen_2020": func(i, j int) bool {
			return P6GasAansluitingen2020.GetValue(items[i].P6GasAansluitingen2020) > P6GasAansluitingen2020.GetValue(items[j].P6GasAansluitingen2020)
		},

		"p6_gasm3_2020": func(i, j int) bool {
			return P6Gasm32020.GetValue(items[i].P6Gasm32020) < P6Gasm32020.GetValue(items[j].P6Gasm32020)
		},
		"-p6_gasm3_2020": func(i, j int) bool {
			return P6Gasm32020.GetValue(items[i].P6Gasm32020) > P6Gasm32020.GetValue(items[j].P6Gasm32020)
		},

		"p6_kwh_2020": func(i, j int) bool {
			return P6Kwh2020.GetValue(items[i].P6Kwh2020) < P6Kwh2020.GetValue(items[j].P6Kwh2020)
		},
		"-p6_kwh_2020": func(i, j int) bool {
			return P6Kwh2020.GetValue(items[i].P6Kwh2020) > P6Kwh2020.GetValue(items[j].P6Kwh2020)
		},

		"p6_totaal_pandoppervlak_m2": func(i, j int) bool {
			return P6TotaalPandoppervlakM2.GetValue(items[i].P6TotaalPandoppervlakM2) < P6TotaalPandoppervlakM2.GetValue(items[j].P6TotaalPandoppervlakM2)
		},
		"-p6_totaal_pandoppervlak_m2": func(i, j int) bool {
			return P6TotaalPandoppervlakM2.GetValue(items[i].P6TotaalPandoppervlakM2) > P6TotaalPandoppervlakM2.GetValue(items[j].P6TotaalPandoppervlakM2)
		},

		"pand_bouwjaar": func(i, j int) bool {
			return PandBouwjaar.GetValue(items[i].PandBouwjaar) < PandBouwjaar.GetValue(items[j].PandBouwjaar)
		},
		"-pand_bouwjaar": func(i, j int) bool {
			return PandBouwjaar.GetValue(items[i].PandBouwjaar) > PandBouwjaar.GetValue(items[j].PandBouwjaar)
		},

		"pand_gas_aansluitingen": func(i, j int) bool {
			return PandGasAansluitingen.GetValue(items[i].PandGasAansluitingen) < PandGasAansluitingen.GetValue(items[j].PandGasAansluitingen)
		},
		"-pand_gas_aansluitingen": func(i, j int) bool {
			return PandGasAansluitingen.GetValue(items[i].PandGasAansluitingen) > PandGasAansluitingen.GetValue(items[j].PandGasAansluitingen)
		},

		"gebruiksdoelen": func(i, j int) bool {
			return Gebruiksdoelen.GetArrayValue(items[i].Gebruiksdoelen) < Gebruiksdoelen.GetArrayValue(items[j].Gebruiksdoelen)
		},
		"-gebruiksdoelen": func(i, j int) bool {
			return Gebruiksdoelen.GetArrayValue(items[i].Gebruiksdoelen) > Gebruiksdoelen.GetArrayValue(items[j].Gebruiksdoelen)
		},
	}
	return sortFuncs
}

func sortBy(items Items, sortingL []string) (Items, []string) {
	sortFuncs := createSort(items)

	for _, sortFuncName := range sortingL {
		sortFunc, ok := sortFuncs[sortFuncName]
		if ok {
			sort.Slice(items, sortFunc)
		}
	}

	// TODO must be nicer way
	keys := []string{}
	for key := range sortFuncs {
		keys = append(keys, key)
	}

	return items, keys
}
