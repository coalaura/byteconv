// Copyright 2026 coalaura. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package byteconv

type decimal struct {
	d     [800]byte // digits, big-endian representation
	nd    int       // number of digits used
	dp    int       // decimal point
	neg   bool      // negative flag
	trunc bool      // discarded nonzero digits beyond d[:nd]
}

// trim trailing zeros from number.
func trim(a *decimal) {
	for a.nd > 0 && a.d[a.nd-1] == '0' {
		a.nd--
	}
	if a.nd == 0 {
		a.dp = 0
	}
}

// Assign v to a.
func (a *decimal) Assign(v uint64) {
	var buf [24]byte

	n := 0
	for v > 0 {
		v1 := v / 10
		v -= 10 * v1
		buf[n] = byte(v + '0')
		n++
		v = v1
	}

	a.nd = 0
	for n--; n >= 0; n-- {
		a.d[a.nd] = buf[n]
		a.nd++
	}
	a.dp = a.nd
	trim(a)
}

const uintSize = 32 << (^uint(0) >> 63)
const maxShift = uintSize - 4

func rightShift(a *decimal, k uint) {
	r := 0
	w := 0

	var n uint
	for ; n>>k == 0; r++ {
		if r >= a.nd {
			if n == 0 {
				a.nd = 0
				return
			}
			for n>>k == 0 {
				n = n * 10
				r++
			}
			break
		}
		c := uint(a.d[r])
		n = n*10 + c - '0'
	}
	a.dp -= r - 1

	var mask uint = (1 << k) - 1

	for ; r < a.nd; r++ {
		c := uint(a.d[r])
		dig := n >> k
		n &= mask
		a.d[w] = byte(dig + '0')
		w++
		n = n*10 + c - '0'
	}

	for n > 0 {
		dig := n >> k
		n &= mask
		if w < len(a.d) {
			a.d[w] = byte(dig + '0')
			w++
		} else if dig > 0 {
			a.trunc = true
		}
		n = n * 10
	}

	a.nd = w
	trim(a)
}

type leftCheat struct {
	delta  int
	cutoff []byte
}

var leftcheats = []leftCheat{
	{0, nil},
	{1, []byte("5")},
	{1, []byte("25")},
	{1, []byte("125")},
	{2, []byte("625")},
	{2, []byte("3125")},
	{2, []byte("15625")},
	{3, []byte("78125")},
	{3, []byte("390625")},
	{3, []byte("1953125")},
	{4, []byte("9765625")},
	{4, []byte("48828125")},
	{4, []byte("244140625")},
	{4, []byte("1220703125")},
	{5, []byte("6103515625")},
	{5, []byte("30517578125")},
	{5, []byte("152587890625")},
	{6, []byte("762939453125")},
	{6, []byte("3814697265625")},
	{6, []byte("19073486328125")},
	{7, []byte("95367431640625")},
	{7, []byte("476837158203125")},
	{7, []byte("2384185791015625")},
	{7, []byte("11920928955078125")},
	{8, []byte("59604644775390625")},
	{8, []byte("298023223876953125")},
	{8, []byte("1490116119384765625")},
	{9, []byte("7450580596923828125")},
	{9, []byte("37252902984619140625")},
	{9, []byte("186264514923095703125")},
	{10, []byte("931322574615478515625")},
	{10, []byte("4656612873077392578125")},
	{10, []byte("23283064365386962890625")},
	{10, []byte("116415321826934814453125")},
	{11, []byte("582076609134674072265625")},
	{11, []byte("2910383045673370361328125")},
	{11, []byte("14551915228366851806640625")},
	{12, []byte("72759576141834259033203125")},
	{12, []byte("363797880709171295166015625")},
	{12, []byte("1818989403545856475830078125")},
	{13, []byte("9094947017729282379150390625")},
	{13, []byte("45474735088646411895751953125")},
	{13, []byte("227373675443232059478759765625")},
	{13, []byte("1136868377216160297393798828125")},
	{14, []byte("5684341886080801486968994140625")},
	{14, []byte("28421709430404007434844970703125")},
	{14, []byte("142108547152020037174224853515625")},
	{15, []byte("710542735760100185871124267578125")},
	{15, []byte("3552713678800500929355621337890625")},
	{15, []byte("17763568394002504646778106689453125")},
	{16, []byte("88817841970012523233890533447265625")},
	{16, []byte("444089209850062616169452667236328125")},
	{16, []byte("2220446049250313080847263336181640625")},
	{16, []byte("11102230246251565404236316680908203125")},
	{17, []byte("55511151231257827021181583404541015625")},
	{17, []byte("277555756156289135105907917022705078125")},
	{17, []byte("1387778780781445675529539585113525390625")},
	{18, []byte("6938893903907228377647697925567626953125")},
	{18, []byte("34694469519536141888238489627838134765625")},
	{18, []byte("173472347597680709441192448139190673828125")},
	{19, []byte("867361737988403547205962240695953369140625")},
}

func prefixIsLessThan(b []byte, s []byte) bool {
	for i := 0; i < len(s); i++ {
		if i >= len(b) {
			return true
		}
		if b[i] != s[i] {
			return b[i] < s[i]
		}
	}
	return false
}

func leftShift(a *decimal, k uint) {
	delta := leftcheats[k].delta
	if prefixIsLessThan(a.d[0:a.nd], leftcheats[k].cutoff) {
		delta--
	}

	r := a.nd
	w := a.nd + delta

	var n uint
	for r--; r >= 0; r-- {
		n += (uint(a.d[r]) - '0') << k
		quo := n / 10
		rem := n - 10*quo
		w--
		if w < len(a.d) {
			a.d[w] = byte(rem + '0')
		} else if rem != 0 {
			a.trunc = true
		}
		n = quo
	}

	for n > 0 {
		quo := n / 10
		rem := n - 10*quo
		w--
		if w < len(a.d) {
			a.d[w] = byte(rem + '0')
		} else if rem != 0 {
			a.trunc = true
		}
		n = quo
	}

	a.nd += delta
	if a.nd >= len(a.d) {
		a.nd = len(a.d)
	}
	a.dp += delta
	trim(a)
}

func (a *decimal) Shift(k int) {
	switch {
	case a.nd == 0:
	case k > 0:
		for k > maxShift {
			leftShift(a, maxShift)
			k -= maxShift
		}
		leftShift(a, uint(k))
	case k < 0:
		for k < -maxShift {
			rightShift(a, maxShift)
			k += maxShift
		}
		rightShift(a, uint(-k))
	}
}

func shouldRoundUp(a *decimal, nd int) bool {
	if nd < 0 || nd >= a.nd {
		return false
	}
	if a.d[nd] == '5' && nd+1 == a.nd {
		if a.trunc {
			return true
		}
		return nd > 0 && (a.d[nd-1]-'0')%2 != 0
	}
	return a.d[nd] >= '5'
}

func (a *decimal) Round(nd int) {
	if nd < 0 || nd >= a.nd {
		return
	}
	if shouldRoundUp(a, nd) {
		a.RoundUp(nd)
	} else {
		a.RoundDown(nd)
	}
}

func (a *decimal) RoundDown(nd int) {
	if nd < 0 || nd >= a.nd {
		return
	}
	a.nd = nd
	trim(a)
}

func (a *decimal) RoundUp(nd int) {
	if nd < 0 || nd >= a.nd {
		return
	}

	for i := nd - 1; i >= 0; i-- {
		c := a.d[i]
		if c < '9' {
			a.d[i]++
			a.nd = i + 1
			return
		}
	}

	a.d[0] = '1'
	a.nd = 1
	a.dp++
}

func (a *decimal) RoundedInteger() uint64 {
	if a.dp > 20 {
		return 0xFFFFFFFFFFFFFFFF
	}
	var i int
	n := uint64(0)
	for i = 0; i < a.dp && i < a.nd; i++ {
		n = n*10 + uint64(a.d[i]-'0')
	}
	for ; i < a.dp; i++ {
		n *= 10
	}
	if shouldRoundUp(a, a.dp) {
		n++
	}
	return n
}
