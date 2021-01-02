package hai

type haiAttribute string

var (
	// suhai
	suhai haiAttribute = "suhai"
	pinzu haiAttribute = "pinzu"
	souzu haiAttribute = "souzu"
	manzu haiAttribute = "manzu"
	one   haiAttribute = "1"
	two   haiAttribute = "2"
	three haiAttribute = "3"
	four  haiAttribute = "4"
	five  haiAttribute = "5"
	six   haiAttribute = "6"
	seven haiAttribute = "7"
	eight haiAttribute = "8"
	nine  haiAttribute = "9"

	// zihai
	zihai haiAttribute = "zihai"
	// kaze
	kaze haiAttribute = "kaze"
	ton  haiAttribute = "ton"
	nan  haiAttribute = "nan"
	sha  haiAttribute = "sha"
	pei  haiAttribute = "pei"
	// sangen
	sangen haiAttribute = "sangen"
	hatu   haiAttribute = "hatu"
	haku   haiAttribute = "haku"
	tyun   haiAttribute = "tyun"
)

type Hai struct {
	attributes []*haiAttribute
	name       string
}

var (
	Pinzu1 = Hai{
		attributes: []*haiAttribute{&suhai, &one, &pinzu},
		name:       "p1",
	}
	Pinzu2 = Hai{
		attributes: []*haiAttribute{&suhai, &two, &pinzu},
		name:       "p2",
	}
	Pinzu3 = Hai{
		attributes: []*haiAttribute{&suhai, &three, &pinzu},
		name:       "p3",
	}
	Pinzu4 = Hai{
		attributes: []*haiAttribute{&suhai, &four, &pinzu},
		name:       "p4",
	}
	Pinzu5 = Hai{
		attributes: []*haiAttribute{&suhai, &five, &pinzu},
		name:       "p5",
	}
	Pinzu6 = Hai{
		attributes: []*haiAttribute{&suhai, &six, &pinzu},
		name:       "p6",
	}
	Pinzu7 = Hai{
		attributes: []*haiAttribute{&suhai, &seven, &pinzu},
		name:       "p7",
	}
	Pinzu8 = Hai{
		attributes: []*haiAttribute{&suhai, &eight, &pinzu},
		name:       "p8",
	}
	Pinzu9 = Hai{
		attributes: []*haiAttribute{&suhai, &nine, &pinzu},
		name:       "p9",
	}
	Souzu1 = Hai{
		attributes: []*haiAttribute{&suhai, &one, &souzu},
		name:       "s1",
	}
	Souzu2 = Hai{
		attributes: []*haiAttribute{&suhai, &two, &souzu},
		name:       "s2",
	}
	Souzu3 = Hai{
		attributes: []*haiAttribute{&suhai, &three, &souzu},
		name:       "s3",
	}
	Souzu4 = Hai{
		attributes: []*haiAttribute{&suhai, &four, &souzu},
		name:       "s4",
	}
	Souzu5 = Hai{
		attributes: []*haiAttribute{&suhai, &five, &souzu},
		name:       "s5",
	}
	Souzu6 = Hai{
		attributes: []*haiAttribute{&suhai, &six, &souzu},
		name:       "s6",
	}
	Souzu7 = Hai{
		attributes: []*haiAttribute{&suhai, &seven, &souzu},
		name:       "s7",
	}
	Souzu8 = Hai{
		attributes: []*haiAttribute{&suhai, &eight, &souzu},
		name:       "s8",
	}
	Souzu9 = Hai{
		attributes: []*haiAttribute{&suhai, &nine, &souzu},
		name:       "s9",
	}
	Manzu1 = Hai{
		attributes: []*haiAttribute{&suhai, &one, &manzu},
		name:       "m1",
	}
	Manzu2 = Hai{
		attributes: []*haiAttribute{&suhai, &two, &manzu},
		name:       "m2",
	}
	Manzu3 = Hai{
		attributes: []*haiAttribute{&suhai, &three, &manzu},
		name:       "m3",
	}
	Manzu4 = Hai{
		attributes: []*haiAttribute{&suhai, &four, &manzu},
		name:       "m4",
	}
	Manzu5 = Hai{
		attributes: []*haiAttribute{&suhai, &five, &manzu},
		name:       "m5",
	}
	Manzu6 = Hai{
		attributes: []*haiAttribute{&suhai, &six, &manzu},
		name:       "m6",
	}
	Manzu7 = Hai{
		attributes: []*haiAttribute{&suhai, &seven, &manzu},
		name:       "m7",
	}
	Manzu8 = Hai{
		attributes: []*haiAttribute{&suhai, &eight, &manzu},
		name:       "m8",
	}
	Manzu9 = Hai{
		attributes: []*haiAttribute{&suhai, &nine, &manzu},
		name:       "m9",
	}

	Tyun = Hai{
		attributes: []*haiAttribute{&zihai, &sangen, &tyun},
		name:       "Tyun",
	}
	Hatu = Hai{
		attributes: []*haiAttribute{&zihai, &sangen, &hatu},
		name:       "hatu",
	}
	Haku = Hai{
		attributes: []*haiAttribute{&zihai, &sangen, &haku},
		name:       "haku",
	}
	Ton = Hai{
		attributes: []*haiAttribute{&zihai, &kaze, &ton},
		name:       "ton",
	}
	Nan = Hai{
		attributes: []*haiAttribute{&zihai, &kaze, &nan},
		name:       "nan",
	}
	Sha = Hai{
		attributes: []*haiAttribute{&zihai, &kaze, &sha},
		name:       "sha",
	}
	Pei = Hai{
		attributes: []*haiAttribute{&zihai, &kaze, &pei},
		name:       "pei",
	}
)
