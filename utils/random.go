package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var cLastTokens = [...]string{
	"BAR", "OUGHT", "ABLE", "PRI", "PRES",
	"ESE", "ANTI", "CALLY", "ATION", "EING"}

// It's used for the non-uniform random generator.
var cLoad int

// cCustomerID is the value of C for the customer id generator. 2.1.6.
var cCustomerID int

// cCustomerID is the value of C for the item id generator. 2.1.6.
var cItemID int

func init() {
	rand.Seed(time.Now().UnixNano())
	cLoad = rand.Intn(256)
	cItemID = rand.Intn(1024)
	cCustomerID = rand.Intn(8192)
}

func GetRandomIntRange(minnum, maxnum int, r *rand.Rand) int {
	return r.Intn(maxnum-minnum+1) + minnum
}
func GetCIDRandom(r *rand.Rand) int {
	return ((r.Intn(1024) | (r.Intn(3000) + 1) + cCustomerID) % 3000) + 1
}
func GetOlquantitiesRandom(r *rand.Rand, Olcnt int) []int {
	arr:=make([]int,Olcnt)
	for i := 0; i < Olcnt; i++ {
		arr[i] = r.Intn(10) + 1
	}
	return arr
}
func GetOliidsRandom(r *rand.Rand, Olcnt int) []int {
	arr:=make([]int,Olcnt)
	for i := 0; i < Olcnt; i++ {
		arr[i] = ((r.Intn(8190) | (r.Intn(100000) + 1) + cItemID) % 100000) + 1
	}
	return arr
}
func GetOlsupplywidsRandom(r *rand.Rand, Olcnt int, wid int) ([]int, int) {
	arr:=make([]int,Olcnt)
	alllocal := 1
	for i := 0; i < Olcnt; i++ {
		if r.Intn(100) == 0 {
			if wid != 1 {
				arr[i] = wid - 1
			} else {
				arr[i] = wid + 1
			}

			alllocal = 0

		} else {
			arr[i] = wid
		}

	}
	return arr, alllocal
}
func GetCwdidRandom(r *rand.Rand, wid, did int) (int, int) {
	if r.Intn(100) < 85 {
		return wid, did
	} else {
		if wid == 1 {
			return wid + 1, r.Intn(10) + 1
		} else {
			return wid - 1, r.Intn(10) + 1
		}
	}
}
func GetCidlastRandom(r *rand.Rand, wid, did int) (int, string) {
	if r.Intn(100) < 60 {
		return -12345, GetCLastRandom(r)
	} else {
		return GetCIDRandom(r),""
	}
}
func randCLastSyllables(n int) string {
	result := ""
	for i := 0; i < 3; i++ {
		result = cLastTokens[n%10] + result
		n /= 10
	}
	return result
}

// See 4.3.2.3.
func GetCLastRandom(r *rand.Rand) string {
	return randCLastSyllables(((r.Intn(256) | r.Intn(1000)) + cLoad) % 1000)
}

// Return a non-uniform random item ID. See 2.1.6.
func GetItemIDRandom(r *rand.Rand) int {
	return ((r.Intn(8190) | (r.Intn(100000) + 1) + cItemID) % 100000) + 1
}

func GetRandomStrings(len int64) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rs := ""
	var i int64
	for i = 0; i < len; i++ {
		r := rand.Intn(62)
		if r == 0 {
			r = 1
		}
		rs += str[r-1 : r]
	}
	return rs
}
func GetRandomintegers(len int64) int64 {
	var (
		in string
		i  int64
	)
	for i = 0; i < len; i++ {
		in += fmt.Sprintf("%d", rand.Intn(10))
	}
	n, _ := strconv.ParseInt(in, 10, 64)
	return n
}
