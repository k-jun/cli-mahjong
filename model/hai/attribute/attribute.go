package attribute

type Attribute string

var (
	// suhai
	Suhai Attribute = "suhai"
	Pinzu Attribute = "pinzu"
	Souzu Attribute = "souzu"
	Manzu Attribute = "manzu"
	One   Attribute = "1"
	Two   Attribute = "2"
	Three Attribute = "3"
	Four  Attribute = "4"
	Five  Attribute = "5"
	Six   Attribute = "6"
	Seven Attribute = "7"
	Eight Attribute = "8"
	Nine  Attribute = "9"

	// zihai
	Jihai Attribute = "zihai"
	// kaze
	Kaze Attribute = "kaze"
	Ton  Attribute = "ton"
	Nan  Attribute = "nan"
	Sha  Attribute = "sha"
	Pei  Attribute = "pei"
	// sangen
	Sangen Attribute = "sangen"
	Hatsu  Attribute = "hatsu"
	Haku   Attribute = "haku"
	Chun   Attribute = "tyun"
)

var (
	Numbers = []*Attribute{&One, &Two, &Three, &Four, &Five, &Six, &Seven, &Eight, &Nine}
	Suits   = []*Attribute{&Pinzu, &Souzu, &Manzu}
)
