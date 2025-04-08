package main

import (
	"testing"
)

type normalizaCasoTeste struct {
	input string
	want  string
}

func TestNormaliza(t *testing.T) {
	casosTeste := []normalizaCasoTeste{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"123-456-7890", "1234567890"},
		{"1234567892", "1234567892"},
		{"(123)456-7892", "1234567892"},
	}

	for _, ct := range casosTeste {
		t.Run(ct.input, func(t *testing.T) {
			res_correto := normalize(ct.input)
			if res_correto != ct.want {
				t.Errorf("Recebido %s; Queríamos: %s", res_correto, ct.want)

			}

		})
	}
}
