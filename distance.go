package main

/*
Distance analyzer - probability of this number appearing on this position in the set
*/
type Distance struct {
	cache  DistanceVal
	cacheI LotteryData
}

/*
DistanceVal - cache ('knowledge base')
*/
type DistanceVal struct {
	Numbers [76][100]float32
	Mb      [16][100]float32
}

/*
Analyse - builds a set of values ('knowledge base') for future analysis
*/
func (r *Distance) Analyse(ld LotteryData) DistanceVal {
	result := DistanceVal{}

	/* calc distances */
	for i := 1; i <= 75; i++ {
		step := 1
		for step <= 99 {
			for key, value := range ld.Data {
				if key >= len(ld.Data)-step {
					continue
				}
				if !inArray(i, value.numbers) {
					continue
				}
				anotherSet := ld.Data[key+step]
				anotherNumbers := anotherSet.numbers
				if inArray(i, anotherNumbers) {
					result.Numbers[i][step]++
				}
			}
			step++
		}
	}

	for i := 1; i <= 15; i++ {
		step := 1
		for step <= 99 {
			for key, value := range ld.Data {
				if key >= len(ld.Data)-step {
					continue
				}
				if i != value.mega {
					continue
				}
				anotherSet := ld.Data[key+step]
				anotherNumber := anotherSet.mega
				if i == anotherNumber {
					result.Mb[i][step]++
				}
			}
			step++
		}
	}

	max := float32(len(ld.Data))

	for key, value := range result.Numbers {
		for dist := range value {
			result.Numbers[key][dist] = result.Numbers[key][dist] / max
			if result.Numbers[key][dist] == 0 {
				result.Numbers[key][dist] = 1 / 100 / max
			}
		}
	}

	for key, value := range result.Mb {
		for dist := range value {
			result.Mb[key][dist] = result.Mb[key][dist] / max
			if result.Mb[key][dist] == 0 {
				result.Mb[key][dist] = 1 / 100 / max
			}
		}
	}

	r.cache = result
	r.cacheI = ld
	return result
}

/*
CheckSet - analyze a single set
*/
func (r *Distance) CheckSet(n Numbers) ([6]float32, float32, float32) {
	var normal [6]float32
	var normMb float32
	var total float32

	result := DistanceVal{}

	for _, i := range n.Set {
		step := 1
		key := 0
		maxRec := 0
		for step <= 99 {
			maxRec++
			if maxRec > 1000 {
				step++
				continue
			}
			if key > len(r.cacheI.Data)-step {
				continue
			}
			anotherSet := r.cacheI.Data[key+step]
			anotherNumbers := anotherSet.numbers
			if inArray(i, anotherNumbers) {
				result.Numbers[i][step]++
			}
			step++
		}
	}

	i := n.Mb
	step := 1
	key := 0
	maxRec := 0
	for step <= 99 {
		maxRec++
		if maxRec > 1000 {
			step++
			continue
		}
		if key > len(r.cacheI.Data)-step {
			continue
		}
		anotherSet := r.cacheI.Data[key+step]
		anotherNumber := anotherSet.mega
		if i == anotherNumber {
			result.Mb[i][step]++
		}
		step++
	}

	result_ := result
	result = DistanceVal{}

	for _, i := range n.Set {
		result.Numbers[i] = result_.Numbers[i]
	}

	result.Mb[n.Mb] = result_.Mb[n.Mb]

	/* calc probability */
	for key, value := range result.Numbers {
		for dist := range value {
			result.Numbers[key][dist] *= r.cache.Numbers[key][dist]
		}
	}

	for key, value := range result.Mb {
		for dist := range value {
			result.Mb[key][dist] *= r.cache.Mb[key][dist]
		}
	}

	/* calc probability for positions */
	for pos, i := range n.Set {
		normal[pos+1] = 0
		for _, v := range result.Numbers[i] {
			normal[pos+1] += v
		}
	}

	for _, v := range result.Mb[n.Mb] {
		normMb += v
	}

	total = 1.0

	for _, n := range normal {
		if n == 0 {
			continue
		}
		total *= n
	}

	total *= normMb * 1000000

	return normal, normMb, total
}

func inArray(i int, list []int) bool {
	for _, v := range list {
		if v == i {
			return true
		}
	}
	return false
}
