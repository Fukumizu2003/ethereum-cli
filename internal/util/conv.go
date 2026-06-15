package util

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"math/big"
	"slices"
	"strconv"
	"strings"
)

func bits(num int) int {
	ans := 0
	for num != 0 {
		num = num >> 1
		ans++
	}
	return ans
}

func B64Encode(msg []byte) string {
	encoded := base64.StdEncoding.EncodeToString(msg)
	return encoded
}

func B64Decode(msg string) []byte {
	decoded, _ := base64.StdEncoding.DecodeString(msg)
	return decoded
}

func GweiToEth(gwei string) string {
	digits := len(gwei)
	var ans string
	if digits <= 9 {
		zeroNum := 9 - digits
		zeros := strings.Repeat("0", zeroNum)
		ans = "0." + zeros + gwei
	} else {
		gweiByte := []byte(gwei)
		bigByte := gweiByte[:len(gweiByte)-9]
		smallByte := gweiByte[len(gweiByte)-9]
		big := string(bigByte)
		small := string(smallByte)
		ans = big + "." + small
	}
	ansByte := []byte(ans)
	prevLength := len(ansByte)
	for i := prevLength - 1; ansByte[i] == byte('0'); i-- {
		ansByte = ansByte[:len(ansByte)-1]
	}
	if ansByte[len(ansByte)-1] == byte('.') {
		ansByte = append(ansByte, byte('0'))
	}
	ans = string(ansByte)
	return ans
}

func EthToGwei(eth string) int {
	gwei := 0
	if strings.Contains(eth, ".") {
		numl := strings.Split(eth, ".")
		big, _ := strconv.Atoi(numl[0])
		smallStr := numl[1] + strings.Repeat("0", 9-len(numl[1]))
		small, _ := strconv.Atoi(smallStr)
		gwei += big * 1000000000
		gwei += small
	} else {
		am, _ := strconv.Atoi(eth)
		gwei += am * 1000000000
	}
	return gwei
}

func EthToWei(eth string) []byte {
	if strings.Contains(eth, ".") {
		numl := strings.Split(eth, ".")
		numl0 := ""
		if numl[0] != "0" {
			numl0 = numl[0]
		}
		weistr := numl0 + numl[1] + strings.Repeat("0", 18-len(numl[1]))
		wei, _ := new(big.Int).SetString(weistr, 10)
		return wei.Bytes()
	} else {
		weistr := eth + strings.Repeat("0", 18)
		wei, _ := new(big.Int).SetString(weistr, 10)
		return wei.Bytes()
	}
}

func WeiToEth(weistr string) string {
	digits := len(weistr)
	var ans string
	if digits <= 18 {
		zeroNum := 18 - digits
		zeros := strings.Repeat("0", zeroNum)
		ans = "0." + zeros + weistr
	} else {
		weistrByte := []byte(weistr)
		bigByte := weistrByte[:len(weistrByte)-18]
		smallByte := weistrByte[len(weistrByte)-18:]
		big := string(bigByte)
		small := string(smallByte)
		ans = big + "." + small
	}
	ansByte := []byte(ans)
	prevLength := len(ansByte)
	for i := prevLength - 1; ansByte[i] == byte('0'); i-- {
		ansByte = ansByte[:len(ansByte)-1]
	}
	if ansByte[len(ansByte)-1] == byte('.') {
		ansByte = append(ansByte, byte('0'))
	}
	ans = string(ansByte)
	return ans
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func StrToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func IntToCompactsize(i int) ([]byte, error) {
	if i < 253 {
		return []byte{byte(i)}, nil
	} else if i < 65536 {
		head := []byte{0xfd}
		buf := make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, uint16(i))
		return append(head, buf...), nil
	} else if i < 4294967296 {
		head := []byte{0xfe}
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(i))
		return append(head, buf...), nil
	} else {
		return nil, errors.New("Too large to convert to conpact size")
	}
}

func PureHex(orig string) string {
	ans := []byte(orig)
	if string(ans[:2]) == "0x" {
		ans = []byte(orig)[2:]
	}
	if len(ans)%2 != 0 {
		ans = append([]byte{byte('0')}, ans...)
	}
	return string(ans)
}

func BytesToInt(orig []byte) uint64 {
	copy := orig[:]
	slices.Reverse(copy)
	sum := uint64(0)
	scale := uint64(1)
	for _, by := range copy {
		sum += uint64(by) * scale
		scale <<= 8
	}
	return sum
}

func IntToBytes(orig uint64) []byte {
	bytes := []byte{}
	for orig != 0 {
		bytes = append(bytes, byte(orig&0xff))
		orig = orig >> 8
	}
	slices.Reverse(bytes)
	return bytes
}

func DecstrToBigint(amstr string, dec int) *big.Int {
	if strings.Contains(amstr, ".") {
		numl := strings.Split(amstr, ".")
		intstr := numl[0] + numl[1] + strings.Repeat("0", dec-len(numl[1]))
		n := new(big.Int)
		n, _ = n.SetString(intstr, 10)
		return n
	} else {
		intstr := amstr + strings.Repeat("0", dec)
		n := new(big.Int)
		n, _ = n.SetString(intstr, 10)
		return n
	}
}

func IntstrToFloatstr(amstr string, digits int) string {
	var str, f, l string
	if len(amstr) > digits {
		f = amstr[:len(amstr)-digits]
		l = amstr[len(amstr)-digits:]
	} else {
		f = "0"
		l = strings.Repeat("0", digits-len(amstr)) + amstr
	}
	str = f + "." + l
	for str[len(str)-1] == '0' {
		str = str[:len(str)-1]
	}
	if str[len(str)-1] == '.' {
		str = str + "0"
	}
	return str
}

func RLPconv(origin []byte, savezero bool) []byte {
	begin := 0
	if !savezero {
		for _, b := range origin {
			if b == 0 {
				begin++
			} else {
				break
			}
		}
	}
	orig := origin[begin:]
	length := len(orig)
	if length == 0 {
		return []byte{0x80}
	}
	if length == 1 && orig[0] < 0x80 {
		return orig
	}
	if length < 56 {
		pref := []byte{byte(0x80 + len(orig))}
		return append(pref, orig...)
	}
	lengthBytes := IntToBytes(uint64(length))
	bytes := len(lengthBytes)
	pref := append([]byte{byte(0xb7 + bytes)}, lengthBytes...)
	return append(pref, orig...)
}

func RLPlistConv(orig []byte) []byte {
	length := len(orig)
	if length == 0 {
		return []byte{0xc0}
	}
	if length < 56 {
		pref := []byte{byte(0xc0 + length)}
		return append(pref, orig...)
	}
	lengthBytes := IntToBytes(uint64(length))
	bytes := len(lengthBytes)
	pref := append([]byte{byte(0xf7 + bytes)}, lengthBytes...)
	return append(pref, orig...)
}
