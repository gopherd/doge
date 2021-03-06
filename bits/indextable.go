package bits

var bitCountArray = [...]int{
	0, 1, 1, 2, 1, 2, 2, 3,
	1, 2, 2, 3, 2, 3, 3, 4,
	1, 2, 2, 3, 2, 3, 3, 4,
	2, 3, 3, 4, 3, 4, 4, 5,
	1, 2, 2, 3, 2, 3, 3, 4,
	2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5,
	3, 4, 4, 5, 4, 5, 5, 6,
	1, 2, 2, 3, 2, 3, 3, 4,
	2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5,
	3, 4, 4, 5, 4, 5, 5, 6,
	2, 3, 3, 4, 3, 4, 4, 5,
	3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6,
	4, 5, 5, 6, 5, 6, 6, 7,
	1, 2, 2, 3, 2, 3, 3, 4,
	2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5,
	3, 4, 4, 5, 4, 5, 5, 6,
	2, 3, 3, 4, 3, 4, 4, 5,
	3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6,
	4, 5, 5, 6, 5, 6, 6, 7,
	2, 3, 3, 4, 3, 4, 4, 5,
	3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6,
	4, 5, 5, 6, 5, 6, 6, 7,
	3, 4, 4, 5, 4, 5, 5, 6,
	4, 5, 5, 6, 5, 6, 6, 7,
	4, 5, 5, 6, 5, 6, 6, 7,
	5, 6, 6, 7, 6, 7, 7, 8,
}

func count(x uint8) int { return bitCountArray[x] }

var normalizerArray8 = [...]uint8{
	0, 1, 1, 3, 1, 3, 3, 7,
	1, 3, 3, 7, 3, 7, 7, 15,
	1, 3, 3, 7, 3, 7, 7, 15,
	3, 7, 7, 15, 7, 15, 15, 31,
	1, 3, 3, 7, 3, 7, 7, 15,
	3, 7, 7, 15, 7, 15, 15, 31,
	3, 7, 7, 15, 7, 15, 15, 31,
	7, 15, 15, 31, 15, 31, 31, 63,
	1, 3, 3, 7, 3, 7, 7, 15,
	3, 7, 7, 15, 7, 15, 15, 31,
	3, 7, 7, 15, 7, 15, 15, 31,
	7, 15, 15, 31, 15, 31, 31, 63,
	3, 7, 7, 15, 7, 15, 15, 31,
	7, 15, 15, 31, 15, 31, 31, 63,
	7, 15, 15, 31, 15, 31, 31, 63,
	15, 31, 31, 63, 31, 63, 63, 127,
	1, 3, 3, 7, 3, 7, 7, 15,
	3, 7, 7, 15, 7, 15, 15, 31,
	3, 7, 7, 15, 7, 15, 15, 31,
	7, 15, 15, 31, 15, 31, 31, 63,
	3, 7, 7, 15, 7, 15, 15, 31,
	7, 15, 15, 31, 15, 31, 31, 63,
	7, 15, 15, 31, 15, 31, 31, 63,
	15, 31, 31, 63, 31, 63, 63, 127,
	3, 7, 7, 15, 7, 15, 15, 31,
	7, 15, 15, 31, 15, 31, 31, 63,
	7, 15, 15, 31, 15, 31, 31, 63,
	15, 31, 31, 63, 31, 63, 63, 127,
	7, 15, 15, 31, 15, 31, 31, 63,
	15, 31, 31, 63, 31, 63, 63, 127,
	15, 31, 31, 63, 31, 63, 63, 127,
	31, 63, 63, 127, 63, 127, 127, 255,
}

var normalizerArray4 = [...]uint8{
	0, 1, 1, 3, 1, 3, 3, 7,
	1, 3, 3, 7, 3, 7, 7, 15,
	16, 17, 17, 19, 17, 19, 19, 23,
	17, 19, 19, 23, 19, 23, 23, 31,
	16, 17, 17, 19, 17, 19, 19, 23,
	17, 19, 19, 23, 19, 23, 23, 31,
	48, 49, 49, 51, 49, 51, 51, 55,
	49, 51, 51, 55, 51, 55, 55, 63,
	16, 17, 17, 19, 17, 19, 19, 23,
	17, 19, 19, 23, 19, 23, 23, 31,
	48, 49, 49, 51, 49, 51, 51, 55,
	49, 51, 51, 55, 51, 55, 55, 63,
	48, 49, 49, 51, 49, 51, 51, 55,
	49, 51, 51, 55, 51, 55, 55, 63,
	112, 113, 113, 115, 113, 115, 115, 119,
	113, 115, 115, 119, 115, 119, 119, 127,
	16, 17, 17, 19, 17, 19, 19, 23,
	17, 19, 19, 23, 19, 23, 23, 31,
	48, 49, 49, 51, 49, 51, 51, 55,
	49, 51, 51, 55, 51, 55, 55, 63,
	48, 49, 49, 51, 49, 51, 51, 55,
	49, 51, 51, 55, 51, 55, 55, 63,
	112, 113, 113, 115, 113, 115, 115, 119,
	113, 115, 115, 119, 115, 119, 119, 127,
	48, 49, 49, 51, 49, 51, 51, 55,
	49, 51, 51, 55, 51, 55, 55, 63,
	112, 113, 113, 115, 113, 115, 115, 119,
	113, 115, 115, 119, 115, 119, 119, 127,
	112, 113, 113, 115, 113, 115, 115, 119,
	113, 115, 115, 119, 115, 119, 119, 127,
	240, 241, 241, 243, 241, 243, 243, 247,
	241, 243, 243, 247, 243, 247, 247, 255,
}
