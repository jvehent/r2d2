package main

import (
	"fmt"
	"time"
)

func stardateCalc() string {
	cur := time.Now().UTC()
	c := cur.Year()/100 - 19
	m := int(cur.Month())
	if m >= 10 {
		m = 0
	}
	d := cur.YearDay()
	h := cur.Hour()
	ret := fmt.Sprintf("current stardate: %v%v%v.%v", c, int(m), d, h)
	return ret
}
