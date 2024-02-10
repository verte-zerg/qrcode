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
			if res := gfMul(test.a, test.b); res != test.res {
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
			if res := gfDiv(test.a, test.b); res != test.res {
				t.Errorf("Expected %v, got %v", test.res, res)
			}
		})
	}
}

func TestPolynomialAdd(t *testing.T) {
	tests := []struct {
		a, b, res *polynomial
	}{
		{&polynomial{[]byte{1, 2, 3}}, &polynomial{[]byte{1, 2, 3}}, &polynomial{[]byte{0, 0, 0}}},
		{&polynomial{[]byte{1, 2, 3}}, &polynomial{[]byte{4, 1, 2, 3}}, &polynomial{[]byte{4, 0, 0, 0}}},
		{&polynomial{[]byte{1, 2, 3}}, &polynomial{[]byte{4, 1, 2, 3, 4}}, &polynomial{[]byte{4, 1, 3, 1, 7}}},
		{&polynomial{[]byte{192, 150}}, &polynomial{[]byte{54, 167, 42, 86, 95}}, &polynomial{[]byte{54, 167, 42, 150, 201}}},
		{&polynomial{[]byte{26}}, &polynomial{[]byte{2}}, &polynomial{[]byte{24}}},
		{&polynomial{[]byte{188, 25, 255}}, &polynomial{[]byte{187, 30, 220}}, &polynomial{[]byte{7, 7, 35}}},
		{&polynomial{[]byte{141, 246}}, &polynomial{[]byte{71, 17, 111, 195}}, &polynomial{[]byte{71, 17, 226, 53}}},
		{&polynomial{[]byte{149, 191, 128}}, &polynomial{[]byte{241, 192, 229, 179, 203}}, &polynomial{[]byte{241, 192, 112, 12, 75}}},
		{&polynomial{[]byte{19, 30, 39}}, &polynomial{[]byte{73, 73, 106, 38, 87}}, &polynomial{[]byte{73, 73, 121, 56, 112}}},
		{&polynomial{[]byte{58}}, &polynomial{[]byte{18, 92, 26, 183, 134}}, &polynomial{[]byte{18, 92, 26, 183, 188}}},
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
		a, b, res *polynomial
	}{
		{&polynomial{[]byte{1, 2, 3}}, &polynomial{[]byte{1, 2, 3}}, &polynomial{[]byte{1, 0, 4, 0, 5}}},
		{&polynomial{[]byte{228, 230, 113, 50}}, &polynomial{[]byte{0}}, &polynomial{[]byte{0, 0, 0, 0}}},
		{&polynomial{[]byte{185, 177}}, &polynomial{[]byte{144}}, &polynomial{[]byte{157, 105}}},
		{&polynomial{[]byte{195, 254}}, &polynomial{[]byte{145, 217, 174, 196}}, &polynomial{[]byte{102, 253, 63, 97, 76}}},
		{&polynomial{[]byte{59, 184}}, &polynomial{[]byte{81, 88, 86, 108, 220, 144}}, &polynomial{[]byte{202, 226, 172, 165, 216, 4, 13}}},
		{&polynomial{[]byte{240, 22}}, &polynomial{[]byte{84, 112, 12, 210}}, &polynomial{[]byte{138, 202, 90, 201, 119}}},
		{&polynomial{[]byte{90, 34, 26}}, &polynomial{[]byte{47, 8, 66}}, &polynomial{[]byte{254, 61, 75, 252, 250}}},
		{&polynomial{[]byte{77, 31}}, &polynomial{[]byte{227, 188, 192}}, &polynomial{[]byte{97, 141, 118, 168}}},
		{&polynomial{[]byte{32, 178, 92, 76}}, &polynomial{[]byte{137, 176, 103, 232, 123, 49}}, &polynomial{[]byte{240, 227, 220, 53, 1, 167, 212, 31, 141}}},
		{&polynomial{[]byte{215, 81, 110, 45, 100}}, &polynomial{[]byte{141, 59, 44, 248, 189}}, &polynomial{[]byte{129, 201, 61, 60, 172, 76, 38, 145, 140}}},
		{&polynomial{[]byte{160, 33}}, &polynomial{[]byte{40, 164, 149}}, &polynomial{[]byte{208, 156, 139, 194}}},
		{&polynomial{[]byte{125, 44, 133, 160, 160, 134}}, &polynomial{[]byte{180}}, &polynomial{[]byte{29, 32, 82, 15, 15, 147}}},
		{&polynomial{[]byte{153, 241, 255}}, &polynomial{[]byte{195, 190, 0}}, &polynomial{[]byte{48, 196, 89, 44, 0}}},
		{&polynomial{[]byte{26, 205, 94, 155, 159, 201}}, &polynomial{[]byte{216, 224, 38}}, &polynomial{[]byte{34, 53, 181, 211, 242, 155, 39, 148}}},
		{&polynomial{[]byte{67, 135, 38, 245}}, &polynomial{[]byte{220}}, &polynomial{[]byte{96, 28, 112, 99}}},
		{&polynomial{[]byte{67, 247}}, &polynomial{[]byte{88, 193, 147, 61}}, &polynomial{[]byte{107, 127, 113, 101, 167}}},
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
		a, b, res *polynomial
	}{
		{&polynomial{[]byte{119, 226}}, &polynomial{[]byte{145, 250, 254}}, &polynomial{[]byte{119, 226}}},
		{&polynomial{[]byte{216}}, &polynomial{[]byte{253, 123, 170, 154, 44}}, &polynomial{[]byte{216}}},
		{&polynomial{[]byte{75, 115}}, &polynomial{[]byte{5, 15, 19}}, &polynomial{[]byte{75, 115}}},
		{&polynomial{[]byte{64, 190, 60, 139}}, &polynomial{[]byte{50, 140, 69}}, &polynomial{[]byte{172, 132}}},
		{&polynomial{[]byte{152, 10}}, &polynomial{[]byte{229, 36, 246, 207}}, &polynomial{[]byte{152, 10}}},
		{&polynomial{[]byte{236, 109, 114, 176, 103}}, &polynomial{[]byte{64, 121}}, &polynomial{[]byte{173}}},
		{&polynomial{[]byte{134, 18}}, &polynomial{[]byte{136, 76, 48}}, &polynomial{[]byte{134, 18}}},
		{&polynomial{[]byte{95, 142, 23, 166, 84}}, &polynomial{[]byte{113, 250}}, &polynomial{[]byte{234}}},
		{&polynomial{[]byte{107, 240, 13, 36}}, &polynomial{[]byte{254}}, &polynomial{[]byte{}}},
		{&polynomial{[]byte{143}}, &polynomial{[]byte{21, 92, 120, 251, 63, 33}}, &polynomial{[]byte{143}}},
		{&polynomial{[]byte{126, 192, 189, 52, 167, 108}}, &polynomial{[]byte{175}}, &polynomial{[]byte{}}},
		{&polynomial{[]byte{23, 220, 153}}, &polynomial{[]byte{221, 176, 96, 23}}, &polynomial{[]byte{23, 220, 153}}},
		{&polynomial{[]byte{60, 241, 52, 178, 12}}, &polynomial{[]byte{127, 198, 84, 204}}, &polynomial{[]byte{116, 242, 199}}},
		{&polynomial{[]byte{153, 242, 127, 166, 48}}, &polynomial{[]byte{92, 195, 15, 241, 42, 16}}, &polynomial{[]byte{153, 242, 127, 166, 48}}},
		{&polynomial{[]byte{15}}, &polynomial{[]byte{108, 87, 19, 240, 205, 115}}, &polynomial{[]byte{15}}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v / %v", test.a, test.b)
		t.Run(name, func(t *testing.T) {
			if res := test.a.Modulo(test.b); !bytes.Equal(res.Coefficients, test.res.Coefficients) {
				t.Errorf("Expected %v, got %v", test.res, res)
			}
		})
	}
}

func TestPolynomialNormalize(t *testing.T) {
	tests := []struct {
		p, res *polynomial
	}{
		{&polynomial{[]byte{0, 0, 0, 0, 1}}, &polynomial{[]byte{1}}},
		{&polynomial{[]byte{0, 0, 0, 0, 0, 1}}, &polynomial{[]byte{1}}},
		{&polynomial{[]byte{0, 0, 0, 0, 0, 1, 0}}, &polynomial{[]byte{1, 0}}},
		{&polynomial{[]byte{0, 0, 1, 0, 0, 0, 0, 0}}, &polynomial{[]byte{1, 0, 0, 0, 0, 0}}},
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
		res     *polynomial
	}{
		{1, &polynomial{[]byte{1}}},
		{2, &polynomial{[]byte{1, 0}}},
		{3, &polynomial{[]byte{1, 1}}},
		{4, &polynomial{[]byte{1, 0, 0}}},
		{5, &polynomial{[]byte{1, 0, 1}}},
		{6, &polynomial{[]byte{1, 1, 0}}},
		{7, &polynomial{[]byte{1, 1, 1}}},
		{15, &polynomial{[]byte{1, 1, 1, 1}}},
	}

	for _, test := range tests {
		if res := convertByteToPolynomial(test.version).Normalize(); !bytes.Equal(res.Coefficients, test.res.Coefficients) {
			t.Errorf("Expected %v, got %v", test.res, res)
		}
	}
}
