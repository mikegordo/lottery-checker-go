package main

import (
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type Numbers struct {
	Set []int
	Mb  int
}

type Builder struct {
	Numbers Numbers
	random  *rand.Rand
}

func (b *Builder) Initialize() {
	randSource := rand.NewSource(time.Now().UnixNano())
	b.random = rand.New(randSource)
}

func (b *Builder) Populate() {
	b.Numbers.Set = make([]int, 5)

	tmp := b.random.Perm(75) // [0, 75)
	b.Numbers.Set[0] = tmp[0] + 1
	b.Numbers.Set[1] = tmp[1] + 1
	b.Numbers.Set[2] = tmp[2] + 1
	b.Numbers.Set[3] = tmp[3] + 1
	b.Numbers.Set[4] = tmp[4] + 1

	sort.Ints(b.Numbers.Set)

	b.Numbers.Mb = b.random.Intn(15) + 1
}

func (b *Builder) GetNumbersString() string {
	output := ""

	for _, v := range b.Numbers.Set {
		output += strconv.Itoa(v) + " "
	}

	output += strconv.Itoa(b.Numbers.Mb)
	return output
}
