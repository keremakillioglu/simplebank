package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghjijklmnoprstuvwxyz"

// init function will be called automaticall when the package is first used
func init() {
	// this ensures that our data will be different whenever we run the code
	rand.Seed(time.Now().UnixNano())
	// if we dont call rand.seed rand will behave as the seed val= 1, and generated val would be the same for everyone

}

//RandomInt generates a random int btw min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // returns a rand int btw 0->max-min
}

//RandomString generates a random string length of n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		// intn returns a rand int from 0 to k
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

//RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

//RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

//RandomCurrency generates a random currency from a provided list
func RandomCurrency() string {
	currencies := []string{EUR, USD, TRY}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

//RandomEmail generates a random email
func RandomEmail() string {

	return fmt.Sprintf("%s@email.com", RandomString(6))
}
