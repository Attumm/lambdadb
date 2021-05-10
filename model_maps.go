/*
	Transforming ItemsIn -> Items -> ItemsOut
	Where Items has column values ar integers to save memmory
	maps are needed to restore integers back to the actual values.
	those are generated and stored here.
*/

package main

type ModelMaps struct {
	Pid                     MappedColumn
	Vid                     MappedColumn
	Straat                  MappedColumn
	Postcode                MappedColumn
	Huisnummer              MappedColumn
	Huisletter              MappedColumn
	Huisnummertoevoeging    MappedColumn
	Oppervlakte             MappedColumn
	Woningequivalent        MappedColumn
	WoningType              MappedColumn
	LabelscoreVoorlopig     MappedColumn
	LabelscoreDefinitief    MappedColumn
	Energieklasse           MappedColumn
	Gemeentecode            MappedColumn
	Gemeentenaam            MappedColumn
	Buurtcode               MappedColumn
	Buurtnaam               MappedColumn
	Wijkcode                MappedColumn
	Wijknaam                MappedColumn
	Provinciecode           MappedColumn
	Provincienaam           MappedColumn
	PandGasEanAansluitingen MappedColumn
	P6GasAansluitingen2020  MappedColumn
	P6Gasm32020             MappedColumn
	P6Kwh2020               MappedColumn
	P6TotaalPandoppervlakM2 MappedColumn
	PandBouwjaar            MappedColumn
	PandGasAansluitingen    MappedColumn
	Gebruiksdoelen          MappedColumn
}

var BitArrays map[string]fieldBitarrayMap

var Pid MappedColumn
var Vid MappedColumn
var Straat MappedColumn
var Postcode MappedColumn
var Huisnummer MappedColumn
var Huisletter MappedColumn
var Huisnummertoevoeging MappedColumn
var Oppervlakte MappedColumn
var Woningequivalent MappedColumn
var WoningType MappedColumn
var LabelscoreVoorlopig MappedColumn
var LabelscoreDefinitief MappedColumn
var Energieklasse MappedColumn
var Gemeentecode MappedColumn
var Gemeentenaam MappedColumn
var Buurtcode MappedColumn
var Buurtnaam MappedColumn
var Wijkcode MappedColumn
var Wijknaam MappedColumn
var Provinciecode MappedColumn
var Provincienaam MappedColumn
var PandGasEanAansluitingen MappedColumn
var P6GasAansluitingen2020 MappedColumn
var P6Gasm32020 MappedColumn
var P6Kwh2020 MappedColumn
var P6TotaalPandoppervlakM2 MappedColumn
var PandBouwjaar MappedColumn
var PandGasAansluitingen MappedColumn
var Gebruiksdoelen MappedColumn

func clearBitArrays() {
	BitArrays = make(map[string]fieldBitarrayMap)
}

func init() {
	clearBitArrays()
	setUpRepeatedColumns()
}

func setUpRepeatedColumns() {
	Pid = NewReapeatedColumn("pid")
	Vid = NewReapeatedColumn("vid")
	Straat = NewReapeatedColumn("straat")
	Postcode = NewReapeatedColumn("postcode")
	Huisnummer = NewReapeatedColumn("huisnummer")
	Huisletter = NewReapeatedColumn("huisletter")
	Huisnummertoevoeging = NewReapeatedColumn("huisnummertoevoeging")
	Oppervlakte = NewReapeatedColumn("oppervlakte")
	Woningequivalent = NewReapeatedColumn("woningequivalent")
	WoningType = NewReapeatedColumn("woning_type")
	LabelscoreVoorlopig = NewReapeatedColumn("labelscore_voorlopig")
	LabelscoreDefinitief = NewReapeatedColumn("labelscore_definitief")
	Energieklasse = NewReapeatedColumn("energieklasse")
	Gemeentecode = NewReapeatedColumn("gemeentecode")
	Gemeentenaam = NewReapeatedColumn("gemeentenaam")
	Buurtcode = NewReapeatedColumn("buurtcode")
	Buurtnaam = NewReapeatedColumn("buurtnaam")
	Wijkcode = NewReapeatedColumn("wijkcode")
	Wijknaam = NewReapeatedColumn("wijknaam")
	Provinciecode = NewReapeatedColumn("provinciecode")
	Provincienaam = NewReapeatedColumn("provincienaam")
	PandGasEanAansluitingen = NewReapeatedColumn("pand_gas_ean_aansluitingen")
	P6GasAansluitingen2020 = NewReapeatedColumn("p6_gas_aansluitingen_2020")
	P6Gasm32020 = NewReapeatedColumn("p6_gasm3_2020")
	P6Kwh2020 = NewReapeatedColumn("p6_kwh_2020")
	P6TotaalPandoppervlakM2 = NewReapeatedColumn("p6_totaal_pandoppervlak_m2")
	PandBouwjaar = NewReapeatedColumn("pand_bouwjaar")
	PandGasAansluitingen = NewReapeatedColumn("pand_gas_aansluitingen")
	Gebruiksdoelen = NewReapeatedColumn("gebruiksdoelen")

}

func CreateMapstore() ModelMaps {
	return ModelMaps{
		Pid,
		Vid,
		Straat,
		Postcode,
		Huisnummer,
		Huisletter,
		Huisnummertoevoeging,
		Oppervlakte,
		Woningequivalent,
		WoningType,
		LabelscoreVoorlopig,
		LabelscoreDefinitief,
		Energieklasse,
		Gemeentecode,
		Gemeentenaam,
		Buurtcode,
		Buurtnaam,
		Wijkcode,
		Wijknaam,
		Provinciecode,
		Provincienaam,
		PandGasEanAansluitingen,
		P6GasAansluitingen2020,
		P6Gasm32020,
		P6Kwh2020,
		P6TotaalPandoppervlakM2,
		PandBouwjaar,
		PandGasAansluitingen,
		Gebruiksdoelen,
	}
}

func LoadMapstore(m ModelMaps) {

	Pid = m.Pid
	Vid = m.Vid
	Straat = m.Straat
	Postcode = m.Postcode
	Huisnummer = m.Huisnummer
	Huisletter = m.Huisletter
	Huisnummertoevoeging = m.Huisnummertoevoeging
	Oppervlakte = m.Oppervlakte
	Woningequivalent = m.Woningequivalent
	WoningType = m.WoningType
	LabelscoreVoorlopig = m.LabelscoreVoorlopig
	LabelscoreDefinitief = m.LabelscoreDefinitief
	Energieklasse = m.Energieklasse
	Gemeentecode = m.Gemeentecode
	Gemeentenaam = m.Gemeentenaam
	Buurtcode = m.Buurtcode
	Buurtnaam = m.Buurtnaam
	Wijkcode = m.Wijkcode
	Wijknaam = m.Wijknaam
	Provinciecode = m.Provinciecode
	Provincienaam = m.Provincienaam
	PandGasEanAansluitingen = m.PandGasEanAansluitingen
	P6GasAansluitingen2020 = m.P6GasAansluitingen2020
	P6Gasm32020 = m.P6Gasm32020
	P6Kwh2020 = m.P6Kwh2020
	P6TotaalPandoppervlakM2 = m.P6TotaalPandoppervlakM2
	PandBouwjaar = m.PandBouwjaar
	PandGasAansluitingen = m.PandGasAansluitingen
	Gebruiksdoelen = m.Gebruiksdoelen

	// register the columns
	RegisteredColumns[Pid.Name] = Pid
	RegisteredColumns[Vid.Name] = Vid
	RegisteredColumns[Straat.Name] = Straat
	RegisteredColumns[Postcode.Name] = Postcode
	RegisteredColumns[Huisnummer.Name] = Huisnummer
	RegisteredColumns[Huisletter.Name] = Huisletter
	RegisteredColumns[Huisnummertoevoeging.Name] = Huisnummertoevoeging
	RegisteredColumns[Oppervlakte.Name] = Oppervlakte
	RegisteredColumns[Woningequivalent.Name] = Woningequivalent
	RegisteredColumns[WoningType.Name] = WoningType
	RegisteredColumns[LabelscoreVoorlopig.Name] = LabelscoreVoorlopig
	RegisteredColumns[LabelscoreDefinitief.Name] = LabelscoreDefinitief
	RegisteredColumns[Energieklasse.Name] = Energieklasse
	RegisteredColumns[Gemeentecode.Name] = Gemeentecode
	RegisteredColumns[Gemeentenaam.Name] = Gemeentenaam
	RegisteredColumns[Buurtcode.Name] = Buurtcode
	RegisteredColumns[Buurtnaam.Name] = Buurtnaam
	RegisteredColumns[Wijkcode.Name] = Wijkcode
	RegisteredColumns[Wijknaam.Name] = Wijknaam
	RegisteredColumns[Provinciecode.Name] = Provinciecode
	RegisteredColumns[Provincienaam.Name] = Provincienaam
	RegisteredColumns[PandGasEanAansluitingen.Name] = PandGasEanAansluitingen
	RegisteredColumns[P6GasAansluitingen2020.Name] = P6GasAansluitingen2020
	RegisteredColumns[P6Gasm32020.Name] = P6Gasm32020
	RegisteredColumns[P6Kwh2020.Name] = P6Kwh2020
	RegisteredColumns[P6TotaalPandoppervlakM2.Name] = P6TotaalPandoppervlakM2
	RegisteredColumns[PandBouwjaar.Name] = PandBouwjaar
	RegisteredColumns[PandGasAansluitingen.Name] = PandGasAansluitingen
	RegisteredColumns[Gebruiksdoelen.Name] = Gebruiksdoelen
}
