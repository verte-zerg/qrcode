package qrcode

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGFMul(t *testing.T) {
	tests := []struct {
		a, b byte
		res  byte
	}{
		{3, 0, 0},
		{5, 8, 40},
		{50, 8, 141},
		{200, 255, 248},
		{100, 1, 100},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v * %v", test.a, test.b)
		t.Run(name, func(t *testing.T) {
			if res := GFMul(test.a, test.b); res != test.res {
				t.Errorf("Expected %v, got %v", test.res, res)
			}
		})
	}
}

func TestGFDiv(t *testing.T) {
	tests := []struct {
		a, b byte
		res  byte
	}{
		{0, 1, 0},
		{40, 8, 5},
		{141, 8, 50},
		{248, 255, 200},
		{100, 1, 100},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v / %v", test.a, test.b)
		t.Run(name, func(t *testing.T) {
			if res := GFDiv(test.a, test.b); res != test.res {
				t.Errorf("Expected %v, got %v", test.res, res)
			}
		})
	}
}

func TestPolynomialAdd(t *testing.T) {
	tests := []struct {
		a, b, res *Polynomial
	}{
		{&Polynomial{[]byte{1, 2, 3}}, &Polynomial{[]byte{1, 2, 3}}, &Polynomial{[]byte{0, 0, 0}}},
		{&Polynomial{[]byte{1, 2, 3}}, &Polynomial{[]byte{4, 1, 2, 3}}, &Polynomial{[]byte{4, 0, 0, 0}}},
		{&Polynomial{[]byte{1, 2, 3}}, &Polynomial{[]byte{4, 1, 2, 3, 4}}, &Polynomial{[]byte{4, 1, 3, 1, 7}}},
		{&Polynomial{[]byte{192, 150}}, &Polynomial{[]byte{54, 167, 42, 86, 95}}, &Polynomial{[]byte{54, 167, 42, 150, 201}}},
		{&Polynomial{[]byte{26}}, &Polynomial{[]byte{2}}, &Polynomial{[]byte{24}}},
		{&Polynomial{[]byte{188, 25, 255}}, &Polynomial{[]byte{187, 30, 220}}, &Polynomial{[]byte{7, 7, 35}}},
		{&Polynomial{[]byte{141, 246}}, &Polynomial{[]byte{71, 17, 111, 195}}, &Polynomial{[]byte{71, 17, 226, 53}}},
		{&Polynomial{[]byte{149, 191, 128}}, &Polynomial{[]byte{241, 192, 229, 179, 203}}, &Polynomial{[]byte{241, 192, 112, 12, 75}}},
		{&Polynomial{[]byte{19, 30, 39}}, &Polynomial{[]byte{73, 73, 106, 38, 87}}, &Polynomial{[]byte{73, 73, 121, 56, 112}}},
		{&Polynomial{[]byte{58}}, &Polynomial{[]byte{18, 92, 26, 183, 134}}, &Polynomial{[]byte{18, 92, 26, 183, 188}}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v + %v", test.a, test.b)
		t.Run(name, func(t *testing.T) {
			if res := test.a.Add(test.b); !bytes.Equal(res.Coefficients, test.res.Coefficients) {
				t.Errorf("Expected %v, got %v", test.res, res)
			}
		})
	}
}

func TestPolynomialMultiply(t *testing.T) {
	tests := []struct {
		a, b, res *Polynomial
	}{
		{&Polynomial{[]byte{1, 2, 3}}, &Polynomial{[]byte{1, 2, 3}}, &Polynomial{[]byte{1, 0, 4, 0, 5}}},
		{&Polynomial{[]byte{228, 230, 113, 50}}, &Polynomial{[]byte{0}}, &Polynomial{[]byte{0, 0, 0, 0}}},
		{&Polynomial{[]byte{185, 177}}, &Polynomial{[]byte{144}}, &Polynomial{[]byte{157, 105}}},
		{&Polynomial{[]byte{195, 254}}, &Polynomial{[]byte{145, 217, 174, 196}}, &Polynomial{[]byte{102, 253, 63, 97, 76}}},
		{&Polynomial{[]byte{59, 184}}, &Polynomial{[]byte{81, 88, 86, 108, 220, 144}}, &Polynomial{[]byte{202, 226, 172, 165, 216, 4, 13}}},
		{&Polynomial{[]byte{240, 22}}, &Polynomial{[]byte{84, 112, 12, 210}}, &Polynomial{[]byte{138, 202, 90, 201, 119}}},
		{&Polynomial{[]byte{90, 34, 26}}, &Polynomial{[]byte{47, 8, 66}}, &Polynomial{[]byte{254, 61, 75, 252, 250}}},
		{&Polynomial{[]byte{77, 31}}, &Polynomial{[]byte{227, 188, 192}}, &Polynomial{[]byte{97, 141, 118, 168}}},
		{&Polynomial{[]byte{32, 178, 92, 76}}, &Polynomial{[]byte{137, 176, 103, 232, 123, 49}}, &Polynomial{[]byte{240, 227, 220, 53, 1, 167, 212, 31, 141}}},
		{&Polynomial{[]byte{215, 81, 110, 45, 100}}, &Polynomial{[]byte{141, 59, 44, 248, 189}}, &Polynomial{[]byte{129, 201, 61, 60, 172, 76, 38, 145, 140}}},
		{&Polynomial{[]byte{160, 33}}, &Polynomial{[]byte{40, 164, 149}}, &Polynomial{[]byte{208, 156, 139, 194}}},
		{&Polynomial{[]byte{125, 44, 133, 160, 160, 134}}, &Polynomial{[]byte{180}}, &Polynomial{[]byte{29, 32, 82, 15, 15, 147}}},
		{&Polynomial{[]byte{153, 241, 255}}, &Polynomial{[]byte{195, 190, 0}}, &Polynomial{[]byte{48, 196, 89, 44, 0}}},
		{&Polynomial{[]byte{26, 205, 94, 155, 159, 201}}, &Polynomial{[]byte{216, 224, 38}}, &Polynomial{[]byte{34, 53, 181, 211, 242, 155, 39, 148}}},
		{&Polynomial{[]byte{67, 135, 38, 245}}, &Polynomial{[]byte{220}}, &Polynomial{[]byte{96, 28, 112, 99}}},
		{&Polynomial{[]byte{67, 247}}, &Polynomial{[]byte{88, 193, 147, 61}}, &Polynomial{[]byte{107, 127, 113, 101, 167}}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v * %v", test.a, test.b)
		t.Run(name, func(t *testing.T) {
			if res := test.a.Multiply(test.b); !bytes.Equal(res.Coefficients, test.res.Coefficients) {
				t.Errorf("Expected %v, got %v", test.res, res)
			}
		})
	}
}

func TestPolynomialDivide(t *testing.T) {
	tests := []struct {
		a, b, res *Polynomial
	}{
		{&Polynomial{[]byte{119, 226}}, &Polynomial{[]byte{145, 250, 254}}, &Polynomial{[]byte{119, 226}}},
		{&Polynomial{[]byte{216}}, &Polynomial{[]byte{253, 123, 170, 154, 44}}, &Polynomial{[]byte{216}}},
		{&Polynomial{[]byte{75, 115}}, &Polynomial{[]byte{5, 15, 19}}, &Polynomial{[]byte{75, 115}}},
		{&Polynomial{[]byte{64, 190, 60, 139}}, &Polynomial{[]byte{50, 140, 69}}, &Polynomial{[]byte{172, 132}}},
		{&Polynomial{[]byte{152, 10}}, &Polynomial{[]byte{229, 36, 246, 207}}, &Polynomial{[]byte{152, 10}}},
		{&Polynomial{[]byte{236, 109, 114, 176, 103}}, &Polynomial{[]byte{64, 121}}, &Polynomial{[]byte{173}}},
		{&Polynomial{[]byte{134, 18}}, &Polynomial{[]byte{136, 76, 48}}, &Polynomial{[]byte{134, 18}}},
		{&Polynomial{[]byte{95, 142, 23, 166, 84}}, &Polynomial{[]byte{113, 250}}, &Polynomial{[]byte{234}}},
		{&Polynomial{[]byte{107, 240, 13, 36}}, &Polynomial{[]byte{254}}, &Polynomial{[]byte{}}},
		{&Polynomial{[]byte{143}}, &Polynomial{[]byte{21, 92, 120, 251, 63, 33}}, &Polynomial{[]byte{143}}},
		{&Polynomial{[]byte{126, 192, 189, 52, 167, 108}}, &Polynomial{[]byte{175}}, &Polynomial{[]byte{}}},
		{&Polynomial{[]byte{23, 220, 153}}, &Polynomial{[]byte{221, 176, 96, 23}}, &Polynomial{[]byte{23, 220, 153}}},
		{&Polynomial{[]byte{60, 241, 52, 178, 12}}, &Polynomial{[]byte{127, 198, 84, 204}}, &Polynomial{[]byte{116, 242, 199}}},
		{&Polynomial{[]byte{153, 242, 127, 166, 48}}, &Polynomial{[]byte{92, 195, 15, 241, 42, 16}}, &Polynomial{[]byte{153, 242, 127, 166, 48}}},
		{&Polynomial{[]byte{15}}, &Polynomial{[]byte{108, 87, 19, 240, 205, 115}}, &Polynomial{[]byte{15}}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v / %v", test.a, test.b)
		t.Run(name, func(t *testing.T) {
			if res := test.a.Divide(test.b); !bytes.Equal(res.Coefficients, test.res.Coefficients) {
				t.Errorf("Expected %v, got %v", test.res, res)
			}
		})
	}
}

func TestPolynomialNormalize(t *testing.T) {
	tests := []struct {
		p, res *Polynomial
	}{
		{&Polynomial{[]byte{0, 0, 0, 0, 1}}, &Polynomial{[]byte{1}}},
		{&Polynomial{[]byte{0, 0, 0, 0, 0, 1}}, &Polynomial{[]byte{1}}},
		{&Polynomial{[]byte{0, 0, 0, 0, 0, 1, 0}}, &Polynomial{[]byte{1, 0}}},
		{&Polynomial{[]byte{0, 0, 1, 0, 0, 0, 0, 0}}, &Polynomial{[]byte{1, 0, 0, 0, 0, 0}}},
	}
	for _, test := range tests {
		if res := test.p.Normalize(); !bytes.Equal(res.Coefficients, test.res.Coefficients) {
			t.Errorf("Expected %v, got %v", test.res, res)
		}
	}

}

func TestConvertByteToPolynomial(t *testing.T) {
	tests := []struct {
		version byte
		res     *Polynomial
	}{
		{1, &Polynomial{[]byte{1}}},
		{2, &Polynomial{[]byte{1, 0}}},
		{3, &Polynomial{[]byte{1, 1}}},
		{4, &Polynomial{[]byte{1, 0, 0}}},
		{5, &Polynomial{[]byte{1, 0, 1}}},
		{6, &Polynomial{[]byte{1, 1, 0}}},
		{7, &Polynomial{[]byte{1, 1, 1}}},
		{15, &Polynomial{[]byte{1, 1, 1, 1}}},
	}

	for _, test := range tests {
		if res := ConvertByteToPolynomial(test.version).Normalize(); !bytes.Equal(res.Coefficients, test.res.Coefficients) {
			t.Errorf("Expected %v, got %v", test.res, res)
		}
	}
}

func TestPolynomialDivideToVersion(t *testing.T) {
	for version := 1; version <= 40; version++ {
		p := ConvertByteToPolynomial(byte(version))

		// try catch panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic: %v", r)
			}
		}()

		p.Divide(VersionPolynomial)
	}
}
