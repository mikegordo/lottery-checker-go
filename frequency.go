package main

type Frequency struct {
	cache FrequencyVal
}

type FrequencyVal struct {
	Numbers [76][6]float32
	Mb      [16]float32
}

func (r *Frequency) Analyse(ld LotteryData) FrequencyVal {
	result := FrequencyVal{}

	for _, v := range ld.Data {
		numbers := v.numbers
		for pos, number := range numbers {
			result.Numbers[number][pos+1]++
		}

		result.Mb[v.mega]++
	}

	max := float32(len(ld.Data))

	for key, value := range result.Numbers {
		for pos := range value {
			result.Numbers[key][pos] = result.Numbers[key][pos] / max * 100
		}
	}

	for key := range result.Mb {
		result.Mb[key] = result.Mb[key] / max * 100
	}

	r.cache = result

	return result
}

func (r *Frequency) CheckSet(n Numbers) []int {
	var normal []int

	for pos_, number := range n.Set {
		pos := 1 + pos_
		q := r.cache.Numbers[number]
		right := q[pos]

		for _, qq := range q {
			if qq > right {
				normal = append(normal, pos_)
			}
		}
	}

	return normal
}
