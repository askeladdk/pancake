package mathx

import "testing"

func Test_Ident3Inv(t *testing.T) {
	m := Ident3()
	i := m.Inv()
	if i != m {
		t.Fatal("m!=i")
	}
}

func Test_Aff3Inv(t *testing.T) {
	m := Aff3{1, 2, 3, 4, 5, 6}
	i := m.Inv().Inv()
	if i != m {
		t.Fatal("m!=i")
	}
}
