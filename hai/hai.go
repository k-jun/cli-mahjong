package hai

import (
	"mahjong/attribute"
)

type Hai struct {
	attributes []*attribute.HaiAttribute
	name       string
}

func (h *Hai) HasAttribute(attr *attribute.HaiAttribute) bool {
	for _, a := range h.attributes {
		if a == attr {
			return true
		}
	}

	return false
}

var (
	Pinzu1 = Hai{
		attributes: []*attribute.HaiAttribute{
			&attribute.Suhai,
			&attribute.One,
			&attribute.Pinzu,
		},
		name: "p1",
	}
	Pinzu2 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Two, &attribute.Pinzu},
		name:       "p2",
	}
	Pinzu3 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Three, &attribute.Pinzu},
		name:       "p3",
	}
	Pinzu4 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Four, &attribute.Pinzu},
		name:       "p4",
	}
	Pinzu5 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Five, &attribute.Pinzu},
		name:       "p5",
	}
	Pinzu6 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Six, &attribute.Pinzu},
		name:       "p6",
	}
	Pinzu7 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Seven, &attribute.Pinzu},
		name:       "p7",
	}
	Pinzu8 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Eight, &attribute.Pinzu},
		name:       "p8",
	}
	Pinzu9 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Nine, &attribute.Pinzu},
		name:       "p9",
	}
	Souzu1 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.One, &attribute.Souzu},
		name:       "s1",
	}
	Souzu2 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Two, &attribute.Souzu},
		name:       "s2",
	}
	Souzu3 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Three, &attribute.Souzu},
		name:       "s3",
	}
	Souzu4 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Four, &attribute.Souzu},
		name:       "s4",
	}
	Souzu5 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Five, &attribute.Souzu},
		name:       "s5",
	}
	Souzu6 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Six, &attribute.Souzu},
		name:       "s6",
	}
	Souzu7 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Seven, &attribute.Souzu},
		name:       "s7",
	}
	Souzu8 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Eight, &attribute.Souzu},
		name:       "s8",
	}
	Souzu9 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Nine, &attribute.Souzu},
		name:       "s9",
	}
	Manzu1 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.One, &attribute.Manzu},
		name:       "m1",
	}
	Manzu2 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Two, &attribute.Manzu},
		name:       "m2",
	}
	Manzu3 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Three, &attribute.Manzu},
		name:       "m3",
	}
	Manzu4 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Four, &attribute.Manzu},
		name:       "m4",
	}
	Manzu5 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Five, &attribute.Manzu},
		name:       "m5",
	}
	Manzu6 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Six, &attribute.Manzu},
		name:       "m6",
	}
	Manzu7 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Seven, &attribute.Manzu},
		name:       "m7",
	}
	Manzu8 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Eight, &attribute.Manzu},
		name:       "m8",
	}
	Manzu9 = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Suhai, &attribute.Nine, &attribute.Manzu},
		name:       "m9",
	}

	Tyun = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Zihai, &attribute.Sangen, &attribute.Tyun},
		name:       "Tyun",
	}
	Hatu = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Zihai, &attribute.Sangen, &attribute.Hatu},
		name:       "hatu",
	}
	Haku = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Zihai, &attribute.Sangen, &attribute.Haku},
		name:       "haku",
	}
	Ton = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Zihai, &attribute.Kaze, &attribute.Ton},
		name:       "ton",
	}
	Nan = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Zihai, &attribute.Kaze, &attribute.Nan},
		name:       "nan",
	}
	Sha = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Zihai, &attribute.Kaze, &attribute.Sha},
		name:       "sha",
	}
	Pei = Hai{
		attributes: []*attribute.HaiAttribute{&attribute.Zihai, &attribute.Kaze, &attribute.Pei},
		name:       "pei",
	}
)
