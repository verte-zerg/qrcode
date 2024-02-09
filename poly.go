package qrcode

func GFMul(a, b byte) byte {
	if a == 0 || b == 0 {
		return 0
	}
	return Exp[(Log[a]+Log[b])%255]
}

func GFDiv(a, b byte) byte {
	if a == 0 || b == 0 {
		// we ignore impossible devision by 0, because it's not possible
		// as we divide only by polynoms without first zero coefficient
		return 0
	}
	return Exp[(Log[a]+Log[b]*254)%255]
}

func ConvertByteToPolynomial(b byte) *Polynomial {
	return &Polynomial{[]byte{
		b >> 7 & 1,
		b >> 6 & 1,
		b >> 5 & 1,
		b >> 4 & 1,
		b >> 3 & 1,
		b >> 2 & 1,
		b >> 1 & 1,
		b & 1,
	}}
}

type Polynomial struct {
	Coefficients []byte
}

func (p *Polynomial) Copy() *Polynomial {
	coefficients := make([]byte, len(p.Coefficients))
	copy(coefficients, p.Coefficients)
	return &Polynomial{coefficients}
}

func (p *Polynomial) Add(other *Polynomial) *Polynomial {
	var big, small *Polynomial = p, other
	if len(other.Coefficients) > len(p.Coefficients) {
		big, small = other, p
	}

	sumPoly := big.Copy()
	for i := 0; i < len(small.Coefficients); i++ {
		sumPoly.Coefficients[len(sumPoly.Coefficients)-1-i] ^= small.Coefficients[len(small.Coefficients)-1-i]
	}
	return sumPoly
}

func (p *Polynomial) Normalize() *Polynomial {
	normPoly := p.Copy()
	for len(normPoly.Coefficients) > 0 && normPoly.Coefficients[0] == 0 {
		normPoly.Coefficients = normPoly.Coefficients[1:]
	}
	return normPoly
}

func (p *Polynomial) IncreaseDegree(degree int) *Polynomial {
	coefficients := make([]byte, len(p.Coefficients)+degree)
	copy(coefficients, p.Coefficients)
	return &Polynomial{coefficients}
}

func (p *Polynomial) Multiply(other *Polynomial) *Polynomial {
	coefficients := make([]byte, len(p.Coefficients)+len(other.Coefficients)-1)
	for i, a := range p.Coefficients {
		for j, b := range other.Coefficients {
			coefficients[i+j] ^= GFMul(a, b)
		}
	}
	return &Polynomial{coefficients}
}

func (p *Polynomial) Divide(other *Polynomial) *Polynomial {
	var steps int = len(p.Coefficients) - len(other.Coefficients) + 1
	var rest *Polynomial = p.Copy()
	for i := 0; i < steps; i++ {
		if rest.Coefficients[0] != 0 {
			factor := GFDiv(rest.Coefficients[0], other.Coefficients[0])
			factorPoly := &Polynomial{[]byte{factor}}
			sub := factorPoly.Multiply(other)
			sub = sub.IncreaseDegree(len(rest.Coefficients) - len(sub.Coefficients))
			rest = rest.Add(sub)
		}
		rest.Coefficients = rest.Coefficients[1:]
	}

	return rest
}

var Log [256]uint16 = [256]uint16{0, 0, 1, 25, 2, 50, 26, 198, 3, 223, 51, 238, 27, 104, 199, 75, 4, 100, 224, 14, 52, 141, 239, 129, 28, 193, 105, 248, 200, 8, 76, 113, 5, 138, 101, 47, 225, 36, 15, 33, 53, 147, 142, 218, 240, 18, 130, 69, 29, 181, 194, 125, 106, 39, 249, 185, 201, 154, 9, 120, 77, 228, 114, 166, 6, 191, 139, 98, 102, 221, 48, 253, 226, 152, 37, 179, 16, 145, 34, 136, 54, 208, 148, 206, 143, 150, 219, 189, 241, 210, 19, 92, 131, 56, 70, 64, 30, 66, 182, 163, 195, 72, 126, 110, 107, 58, 40, 84, 250, 133, 186, 61, 202, 94, 155, 159, 10, 21, 121, 43, 78, 212, 229, 172, 115, 243, 167, 87, 7, 112, 192, 247, 140, 128, 99, 13, 103, 74, 222, 237, 49, 197, 254, 24, 227, 165, 153, 119, 38, 184, 180, 124, 17, 68, 146, 217, 35, 32, 137, 46, 55, 63, 209, 91, 149, 188, 207, 205, 144, 135, 151, 178, 220, 252, 190, 97, 242, 86, 211, 171, 20, 42, 93, 158, 132, 60, 57, 83, 71, 109, 65, 162, 31, 45, 67, 216, 183, 123, 164, 118, 196, 23, 73, 236, 127, 12, 111, 246, 108, 161, 59, 82, 41, 157, 85, 170, 251, 96, 134, 177, 187, 204, 62, 90, 203, 89, 95, 176, 156, 169, 160, 81, 11, 245, 22, 235, 122, 117, 44, 215, 79, 174, 213, 233, 230, 231, 173, 232, 116, 214, 244, 234, 168, 80, 88, 175}
var Exp [256]byte = [256]byte{1, 2, 4, 8, 16, 32, 64, 128, 29, 58, 116, 232, 205, 135, 19, 38, 76, 152, 45, 90, 180, 117, 234, 201, 143, 3, 6, 12, 24, 48, 96, 192, 157, 39, 78, 156, 37, 74, 148, 53, 106, 212, 181, 119, 238, 193, 159, 35, 70, 140, 5, 10, 20, 40, 80, 160, 93, 186, 105, 210, 185, 111, 222, 161, 95, 190, 97, 194, 153, 47, 94, 188, 101, 202, 137, 15, 30, 60, 120, 240, 253, 231, 211, 187, 107, 214, 177, 127, 254, 225, 223, 163, 91, 182, 113, 226, 217, 175, 67, 134, 17, 34, 68, 136, 13, 26, 52, 104, 208, 189, 103, 206, 129, 31, 62, 124, 248, 237, 199, 147, 59, 118, 236, 197, 151, 51, 102, 204, 133, 23, 46, 92, 184, 109, 218, 169, 79, 158, 33, 66, 132, 21, 42, 84, 168, 77, 154, 41, 82, 164, 85, 170, 73, 146, 57, 114, 228, 213, 183, 115, 230, 209, 191, 99, 198, 145, 63, 126, 252, 229, 215, 179, 123, 246, 241, 255, 227, 219, 171, 75, 150, 49, 98, 196, 149, 55, 110, 220, 165, 87, 174, 65, 130, 25, 50, 100, 200, 141, 7, 14, 28, 56, 112, 224, 221, 167, 83, 166, 81, 162, 89, 178, 121, 242, 249, 239, 195, 155, 43, 86, 172, 69, 138, 9, 18, 36, 72, 144, 61, 122, 244, 245, 247, 243, 251, 235, 203, 139, 11, 22, 44, 88, 176, 125, 250, 233, 207, 131, 27, 54, 108, 216, 173, 71, 142, 0}
var VersionPolynomial *Polynomial = &Polynomial{[]byte{1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1}}
var GeneratorPolynomials [69]*Polynomial = [69]*Polynomial{
	{[]byte{1}},
	{[]byte{1, 1}},
	{[]byte{1, 3, 2}},
	{[]byte{1, 7, 14, 8}},
	{[]byte{1, 15, 54, 120, 64}},
	{[]byte{1, 31, 198, 63, 147, 116}},
	{[]byte{1, 63, 1, 218, 32, 227, 38}},
	{[]byte{1, 127, 122, 154, 164, 11, 68, 117}},
	{[]byte{1, 255, 11, 81, 54, 239, 173, 200, 24}},
	{[]byte{1, 226, 207, 158, 245, 235, 164, 232, 197, 37}},
	{[]byte{1, 216, 194, 159, 111, 199, 94, 95, 113, 157, 193}},
	{[]byte{1, 172, 130, 163, 50, 123, 219, 162, 248, 144, 116, 160}},
	{[]byte{1, 68, 119, 67, 118, 220, 31, 7, 84, 92, 127, 213, 97}},
	{[]byte{1, 137, 73, 227, 17, 177, 17, 52, 13, 46, 43, 83, 132, 120}},
	{[]byte{1, 14, 54, 114, 70, 174, 151, 43, 158, 195, 127, 166, 210, 234, 163}},
	{[]byte{1, 29, 196, 111, 163, 112, 74, 10, 105, 105, 139, 132, 151, 32, 134, 26}},
	{[]byte{1, 59, 13, 104, 189, 68, 209, 30, 8, 163, 65, 41, 229, 98, 50, 36, 59}},
	{[]byte{1, 119, 66, 83, 120, 119, 22, 197, 83, 249, 41, 143, 134, 85, 53, 125, 99, 79}},
	{[]byte{1, 239, 251, 183, 113, 149, 175, 199, 215, 240, 220, 73, 82, 173, 75, 32, 67, 217, 146}},
	{[]byte{1, 194, 8, 26, 146, 20, 223, 187, 152, 85, 115, 238, 133, 146, 109, 173, 138, 33, 172, 179}},
	{[]byte{1, 152, 185, 240, 5, 111, 99, 6, 220, 112, 150, 69, 36, 187, 22, 228, 198, 121, 121, 165, 174}},
	{[]byte{1, 44, 243, 13, 131, 49, 132, 194, 67, 214, 28, 89, 124, 82, 158, 244, 37, 236, 142, 82, 255, 89}},
	{[]byte{1, 89, 179, 131, 176, 182, 244, 19, 189, 69, 40, 28, 137, 29, 123, 67, 253, 86, 218, 230, 26, 145, 245}},
	{[]byte{1, 179, 68, 154, 163, 140, 136, 190, 152, 25, 85, 19, 3, 196, 27, 113, 198, 18, 130, 2, 120, 93, 41, 71}},
	{[]byte{1, 122, 118, 169, 70, 178, 237, 216, 102, 115, 150, 229, 73, 130, 72, 61, 43, 206, 1, 237, 247, 127, 217, 144, 117}},
	{[]byte{1, 245, 49, 228, 53, 215, 6, 205, 210, 38, 82, 56, 80, 97, 139, 81, 134, 126, 168, 98, 226, 125, 23, 171, 173, 193}},
	{[]byte{1, 246, 51, 183, 4, 136, 98, 199, 152, 77, 56, 206, 24, 145, 40, 209, 117, 233, 42, 135, 68, 70, 144, 146, 77, 43, 94}},
	{[]byte{1, 240, 61, 29, 145, 144, 117, 150, 48, 58, 139, 94, 134, 193, 105, 33, 169, 202, 102, 123, 113, 195, 25, 213, 6, 152, 164, 217}},
	{[]byte{1, 252, 9, 28, 13, 18, 251, 208, 150, 103, 174, 100, 41, 167, 12, 247, 56, 117, 119, 233, 127, 181, 100, 121, 147, 176, 74, 58, 197}},
	{[]byte{1, 228, 193, 196, 48, 170, 86, 80, 217, 54, 143, 79, 32, 88, 255, 87, 24, 15, 251, 85, 82, 201, 58, 112, 191, 153, 108, 132, 143, 170}},
	{[]byte{1, 212, 246, 77, 73, 195, 192, 75, 98, 5, 70, 103, 177, 22, 217, 138, 51, 181, 246, 72, 25, 18, 46, 228, 74, 216, 195, 11, 106, 130, 150}},
	{[]byte{1, 180, 74, 173, 182, 161, 15, 36, 192, 124, 187, 31, 53, 238, 202, 236, 158, 199, 147, 168, 27, 27, 160, 2, 36, 26, 197, 196, 237, 220, 28, 89}},
	{[]byte{1, 116, 64, 52, 174, 54, 126, 16, 194, 162, 33, 33, 157, 176, 197, 225, 12, 59, 55, 253, 228, 148, 47, 179, 185, 24, 138, 253, 20, 142, 55, 172, 88}},
	{[]byte{1, 233, 245, 160, 143, 188, 120, 30, 231, 36, 121, 246, 74, 239, 159, 147, 122, 233, 126, 102, 101, 49, 113, 145, 89, 67, 51, 115, 149, 229, 247, 55, 245, 45}},
	{[]byte{1, 206, 60, 154, 113, 6, 117, 208, 90, 26, 113, 31, 25, 177, 132, 99, 51, 105, 183, 122, 22, 43, 136, 93, 94, 62, 111, 196, 23, 126, 135, 67, 222, 23, 10}},
	{[]byte{1, 128, 113, 84, 231, 131, 204, 112, 112, 50, 51, 154, 48, 33, 191, 146, 190, 26, 236, 248, 11, 6, 37, 195, 129, 51, 61, 38, 140, 29, 191, 96, 102, 206, 105, 214}},
	{[]byte{1, 28, 196, 67, 76, 123, 192, 207, 251, 185, 73, 124, 1, 126, 73, 31, 27, 11, 104, 45, 161, 43, 74, 127, 89, 26, 219, 59, 137, 118, 200, 237, 216, 31, 243, 96, 59}},
	{[]byte{1, 57, 15, 21, 150, 111, 145, 13, 247, 159, 144, 217, 171, 91, 169, 186, 191, 59, 50, 121, 241, 173, 196, 181, 156, 213, 206, 201, 109, 17, 29, 26, 74, 130, 87, 115, 90, 228}},
	{[]byte{1, 115, 78, 148, 61, 244, 210, 125, 226, 140, 43, 227, 198, 180, 190, 193, 206, 53, 231, 140, 199, 31, 138, 25, 108, 176, 252, 155, 212, 198, 131, 219, 96, 11, 45, 59, 146, 185, 25}},
	{[]byte{1, 231, 195, 241, 35, 28, 85, 210, 228, 225, 84, 225, 63, 196, 215, 73, 117, 145, 219, 31, 184, 251, 32, 57, 153, 151, 255, 200, 213, 54, 243, 187, 143, 146, 88, 102, 37, 248, 90, 245}},
	{[]byte{1, 210, 248, 240, 209, 173, 67, 133, 167, 133, 209, 131, 186, 99, 93, 235, 52, 40, 6, 220, 241, 72, 13, 215, 128, 255, 156, 49, 62, 254, 212, 35, 99, 51, 218, 101, 180, 247, 40, 156, 38}},
	{[]byte{1, 184, 126, 20, 66, 149, 9, 164, 91, 108, 45, 187, 39, 204, 189, 50, 128, 178, 176, 189, 97, 177, 229, 127, 217, 220, 115, 62, 123, 199, 81, 196, 28, 211, 75, 148, 53, 78, 176, 42, 41, 160}},
	{[]byte{1, 108, 136, 69, 244, 3, 45, 158, 245, 1, 8, 105, 176, 69, 65, 103, 107, 244, 29, 165, 52, 217, 41, 38, 92, 66, 78, 34, 9, 53, 34, 242, 14, 139, 142, 56, 197, 179, 191, 50, 237, 5, 217}},
	{[]byte{1, 217, 194, 8, 233, 155, 239, 39, 190, 44, 189, 168, 161, 117, 92, 148, 34, 146, 133, 160, 192, 139, 8, 113, 230, 180, 127, 60, 93, 65, 197, 166, 15, 211, 1, 236, 184, 34, 77, 239, 38, 118, 130, 33}},
	{[]byte{1, 174, 128, 111, 118, 188, 207, 47, 160, 252, 165, 225, 125, 65, 3, 101, 197, 58, 77, 19, 131, 2, 11, 238, 120, 84, 222, 18, 102, 199, 62, 153, 99, 20, 50, 155, 41, 221, 229, 74, 46, 31, 68, 202, 49}},
	{[]byte{1, 64, 123, 101, 108, 45, 179, 179, 191, 122, 220, 22, 16, 220, 232, 74, 61, 68, 101, 68, 234, 39, 202, 226, 134, 184, 2, 38, 225, 16, 129, 46, 226, 178, 235, 144, 105, 156, 254, 184, 201, 238, 145, 80, 220, 36}},
	{[]byte{1, 129, 113, 254, 129, 71, 18, 112, 124, 220, 134, 225, 32, 80, 31, 23, 238, 105, 76, 169, 195, 229, 178, 37, 2, 16, 217, 185, 88, 202, 13, 251, 29, 54, 233, 147, 241, 20, 3, 213, 18, 119, 112, 9, 90, 211, 38}},
	{[]byte{1, 30, 198, 122, 91, 240, 252, 86, 103, 13, 117, 172, 137, 90, 14, 100, 17, 182, 65, 119, 242, 101, 93, 33, 209, 51, 220, 147, 108, 87, 158, 174, 30, 102, 131, 182, 96, 184, 64, 105, 242, 81, 145, 18, 73, 109, 163, 111}},
	{[]byte{1, 61, 3, 200, 46, 178, 154, 185, 143, 216, 223, 53, 68, 44, 111, 171, 161, 159, 197, 124, 45, 69, 206, 169, 230, 98, 167, 104, 83, 226, 85, 59, 149, 163, 117, 131, 228, 132, 11, 65, 232, 113, 144, 107, 5, 99, 53, 78, 208}},
	{[]byte{1, 123, 118, 2, 212, 25, 138, 139, 95, 189, 49, 20, 59, 121, 72, 22, 233, 81, 180, 207, 78, 36, 221, 218, 34, 51, 83, 47, 33, 108, 1, 60, 105, 84, 55, 172, 142, 89, 174, 129, 254, 163, 186, 223, 189, 32, 135, 49, 3, 228}},
	{[]byte{1, 247, 51, 213, 209, 198, 58, 199, 159, 162, 134, 224, 25, 156, 8, 162, 206, 100, 176, 224, 36, 159, 135, 157, 230, 102, 162, 46, 230, 176, 239, 176, 15, 60, 181, 87, 157, 31, 190, 151, 47, 61, 62, 235, 255, 151, 215, 239, 247, 109, 167}},
	{[]byte{1, 242, 63, 42, 119, 116, 195, 21, 99, 123, 150, 68, 94, 225, 222, 138, 222, 181, 89, 170, 99, 43, 94, 60, 53, 63, 65, 62, 112, 233, 165, 196, 69, 15, 121, 12, 139, 204, 221, 235, 222, 174, 247, 45, 159, 179, 38, 67, 131, 97, 99, 1}},
	{[]byte{1, 248, 5, 177, 110, 5, 172, 216, 225, 130, 159, 177, 204, 151, 90, 149, 243, 170, 239, 234, 19, 210, 77, 74, 176, 224, 218, 142, 225, 174, 113, 210, 190, 151, 31, 17, 243, 235, 118, 234, 30, 177, 175, 53, 176, 28, 172, 34, 39, 22, 142, 248, 10}},
	{[]byte{1, 236, 249, 245, 79, 14, 232, 64, 167, 151, 101, 242, 237, 220, 185, 41, 56, 202, 15, 39, 154, 179, 131, 199, 81, 213, 219, 224, 235, 187, 193, 72, 112, 122, 252, 128, 186, 139, 235, 28, 151, 52, 142, 145, 19, 41, 1, 186, 181, 192, 171, 242, 246, 136}},
	{[]byte{1, 196, 6, 56, 127, 89, 69, 31, 117, 159, 190, 193, 5, 11, 149, 54, 36, 68, 105, 162, 43, 189, 145, 6, 226, 149, 130, 20, 233, 156, 142, 11, 255, 123, 240, 197, 3, 236, 119, 59, 208, 239, 253, 133, 56, 235, 29, 146, 210, 34, 192, 7, 30, 192, 228}},
	{[]byte{1, 148, 141, 197, 126, 76, 127, 171, 11, 144, 175, 82, 131, 6, 223, 61, 98, 203, 141, 35, 54, 37, 242, 80, 31, 49, 137, 219, 221, 114, 111, 35, 181, 1, 184, 168, 216, 28, 148, 148, 33, 80, 238, 95, 90, 234, 83, 76, 116, 61, 178, 209, 179, 238, 50, 89}},
	{[]byte{1, 52, 59, 104, 213, 198, 195, 129, 248, 4, 163, 27, 99, 37, 56, 112, 122, 64, 168, 142, 114, 169, 81, 215, 162, 205, 66, 204, 42, 98, 54, 219, 241, 174, 24, 116, 214, 22, 149, 34, 151, 73, 83, 217, 201, 99, 111, 12, 200, 131, 170, 57, 112, 166, 180, 111, 116}},
	{[]byte{1, 105, 132, 139, 182, 52, 111, 25, 4, 127, 202, 50, 239, 115, 99, 116, 114, 32, 118, 146, 210, 27, 16, 237, 234, 185, 219, 168, 238, 101, 61, 222, 2, 90, 215, 31, 183, 3, 255, 14, 66, 223, 2, 89, 128, 147, 57, 225, 115, 46, 236, 159, 41, 174, 169, 113, 153, 97}},
	{[]byte{1, 211, 248, 6, 131, 97, 12, 222, 104, 173, 98, 28, 55, 235, 160, 216, 176, 89, 168, 57, 139, 227, 21, 130, 27, 73, 54, 83, 214, 71, 42, 190, 145, 51, 201, 143, 96, 236, 44, 249, 64, 23, 43, 48, 77, 204, 218, 83, 233, 237, 48, 212, 161, 115, 42, 243, 51, 82, 197}},
	{[]byte{1, 186, 124, 247, 232, 100, 155, 8, 83, 60, 194, 48, 63, 150, 52, 199, 224, 152, 47, 73, 242, 137, 238, 140, 119, 67, 111, 71, 236, 19, 119, 162, 84, 58, 13, 104, 179, 18, 186, 142, 216, 72, 247, 69, 50, 44, 237, 209, 211, 171, 207, 171, 39, 5, 177, 239, 22, 150, 150, 49}},
	{[]byte{1, 104, 132, 6, 205, 58, 21, 125, 141, 72, 141, 86, 193, 178, 34, 86, 59, 24, 49, 204, 64, 17, 131, 4, 167, 7, 186, 124, 86, 34, 189, 230, 211, 74, 148, 11, 140, 230, 162, 118, 177, 232, 151, 96, 49, 107, 3, 50, 127, 190, 68, 174, 172, 94, 12, 162, 76, 225, 128, 39, 44}},
	{[]byte{1, 209, 250, 26, 124, 95, 58, 69, 203, 76, 77, 82, 78, 168, 110, 135, 180, 142, 207, 148, 156, 112, 101, 16, 121, 115, 178, 145, 169, 173, 108, 3, 127, 96, 59, 72, 251, 91, 14, 101, 128, 114, 43, 245, 238, 51, 171, 228, 241, 151, 119, 17, 192, 93, 34, 221, 95, 255, 36, 229, 154, 193}},
	{[]byte{1, 190, 112, 31, 67, 188, 9, 27, 199, 249, 113, 1, 236, 74, 201, 4, 61, 105, 118, 128, 26, 169, 120, 125, 199, 94, 30, 9, 225, 101, 5, 94, 206, 50, 152, 121, 102, 49, 156, 69, 237, 235, 232, 122, 164, 41, 197, 242, 106, 124, 64, 28, 17, 6, 207, 98, 43, 204, 239, 37, 110, 103, 52}},
	{[]byte{1, 96, 188, 37, 188, 90, 100, 123, 103, 111, 4, 229, 50, 223, 79, 210, 98, 245, 77, 68, 53, 215, 245, 249, 194, 200, 166, 40, 129, 207, 223, 223, 118, 196, 154, 137, 60, 148, 225, 234, 245, 160, 93, 176, 129, 155, 103, 197, 222, 56, 155, 133, 145, 185, 49, 74, 209, 207, 184, 207, 45, 124, 79, 252}},
	{[]byte{1, 193, 10, 255, 58, 128, 183, 115, 140, 153, 147, 91, 197, 219, 221, 220, 142, 28, 120, 21, 164, 147, 6, 204, 40, 230, 182, 14, 121, 48, 143, 77, 228, 81, 85, 43, 162, 16, 195, 163, 35, 149, 154, 35, 132, 100, 100, 51, 176, 11, 161, 134, 208, 132, 244, 176, 192, 221, 232, 171, 125, 155, 228, 242, 245}},
	{[]byte{1, 158, 183, 131, 44, 74, 11, 249, 133, 134, 43, 159, 33, 7, 133, 91, 86, 189, 66, 63, 3, 97, 194, 19, 105, 11, 164, 219, 100, 69, 57, 179, 70, 253, 205, 210, 174, 61, 90, 160, 81, 91, 129, 122, 74, 185, 116, 35, 231, 1, 130, 74, 179, 255, 41, 133, 202, 63, 111, 164, 205, 143, 226, 94, 31, 106}},
	{[]byte{1, 32, 199, 138, 150, 79, 79, 191, 10, 159, 237, 135, 239, 231, 152, 66, 131, 141, 179, 226, 246, 190, 158, 171, 153, 206, 226, 34, 212, 101, 249, 229, 141, 226, 128, 238, 57, 60, 206, 203, 106, 118, 84, 161, 127, 253, 71, 44, 102, 155, 60, 78, 247, 52, 5, 252, 211, 30, 154, 194, 52, 179, 3, 184, 182, 193, 26}},
	{[]byte{1, 65, 123, 31, 177, 128, 63, 207, 55, 114, 108, 67, 31, 225, 177, 249, 36, 228, 174, 105, 39, 168, 194, 75, 99, 20, 57, 243, 170, 13, 216, 230, 102, 255, 81, 36, 94, 144, 154, 16, 73, 66, 136, 3, 104, 111, 221, 115, 108, 25, 36, 26, 230, 67, 126, 4, 40, 76, 176, 187, 89, 200, 136, 27, 177, 178, 212, 179}},
	{[]byte{1, 131, 115, 9, 39, 18, 182, 60, 94, 223, 230, 157, 142, 119, 85, 107, 34, 174, 167, 109, 20, 185, 112, 145, 172, 224, 170, 182, 107, 38, 107, 71, 246, 230, 225, 144, 20, 14, 175, 226, 245, 20, 219, 212, 51, 158, 88, 63, 36, 199, 4, 80, 157, 211, 239, 255, 7, 119, 11, 235, 12, 34, 149, 204, 8, 32, 29, 99, 11}},
}