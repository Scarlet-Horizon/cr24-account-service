package util

import (
	"github.com/google/uuid"
	"main/model"
	rn "math/rand"
	"time"
)

const dateMin = 1514793600
const dateMax = 1668952800
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const minAmount = 100.0
const maxAmount = 10000.0

func RandomAccount() model.Account {
	//pk := "USER#" + uuid.NewString()
	//sk := "ACCOUNT#" + uuid.NewString()

	//max := new(big.Int)
	//max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))

	//n, err := rand.Int(rand.Reader, max)
	//if err != nil {
	//	//error handling
	//}

	//accNumber := rn.Intn(1001)
	//amount := rn.Intn(1001)
	//limit := rn.Intn(100)
	//openDate := randomDate()
	//
	//var closeDate time.Time
	//if rn.Intn(11) <= 5 {
	//	closeDate := randomDate()
	//}
	//ty := randomString()

	return model.Account{
		PK:       "USER#" + uuid.NewString(),
		SK:       "ACCOUNT#" + uuid.NewString(),
		Amount:   randomAmount(),
		Limit:    rn.Intn(100),
		OpenDate: randomDate(),
		//CloseDate: randomDate(),
		Type: randomString(),
	}
}

func randomAmount() float64 {
	return minAmount + rn.Float64()*(maxAmount-minAmount)
}

func randomDate() time.Time {
	//min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	//max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := time.Unix(dateMax-dateMin, 0).Unix()

	sec := rn.Int63n(delta) + time.Unix(dateMin, 0).Unix()
	return time.Unix(sec, 0)
}

func randomString() string {
	return random(rn.Intn(21))
}

func random(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rn.Intn(len(letters))]
	}
	return string(b)
}
