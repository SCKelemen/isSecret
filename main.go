package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/seiflotfy/cuckoofilter"
)

type Environment int

const (
	INVALID Environment = iota
	PRODUCTION
	BLUESTEEL
	LOADTEST
	INTEGRATION
	DEVELOPMENT
)

func main() {

	filter := cuckoo.NewFilter(1000)

	secrets := []string{"password", "pass", "password1", "secret"}
	for _, secret := range secrets {
		filter.Insert(protecc(secret))
	}

	for i, secret := range secrets {

		b := filter.Lookup(protecc(secret))

		fmt.Printf("%d\t%s\t->\t%v\n", i, secret, b)
	}

	standards := []string{}
	for i, secret := range secrets {
		val := sha256.Sum256([]byte(secret))
		standards = append(standards, string(val[0:32]))

		b := filter.Lookup(val[0:32])
		fmt.Printf("%d\t%s\t!>>\t%v\n", i, secret, b)
	}

	nots := []string{"dog", "cat", "foo", "bar", "foobar"}

	for i, not := range nots {
		b := filter.Lookup(protecc(not))
		fmt.Printf("%d\t%s\t!>\t%v\n", i, not, b)
	}

}

func protecc(value string) []byte {
	sum := sha256.Sum256([]byte(value))
	val := sha256.Sum256(sum[0:24])
	return val[0:32]
}

func initializeEnvironments() map[Environment]*cuckoo.Filter {
	var envs map[Environment]*cuckoo.Filter
	envs = make(map[Environment]*cuckoo.Filter)

	for i := 0; i < int(DEVELOPMENT); i++ {
		envs[Environment(i)] = cuckoo.NewFilter(1000)
	}
	return envs
}
