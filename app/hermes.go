package app

import "github.com/matcornic/hermes/v2"

var h = &hermes.Hermes{
	Product: hermes.Product{
		Name:      "Shopicano",
		Link:      "http://shopicano.com",
		Copyright: "Copyright © 2020 Shopicano. All rights reserved.",
	},
}

func Hermes() *hermes.Hermes {
	return h
}
