package main

/*
Range analyzer - checks if number is on its normal position
*/
type Range struct {
	cache []RangeVal
}

/*
RangeVal - single cache value
*/
type RangeVal struct {
	Min int
	Max int
}

/*
Analyse - builds a set of values ('knowledge base') for future analysis
*/
func (r *Range) Analyse(ld LotteryData) []RangeVal {
	result := make([]RangeVal, 6)

	result[0].Min = 76
	result[1].Min = 76
	result[2].Min = 76
	result[3].Min = 76
	result[4].Min = 76
	result[5].Min = 76

	for _, v := range ld.Data {
		numbers := v.numbers
		for pos, number := range numbers {
			if result[pos+1].Min > number {
				result[pos+1].Min = number
			}
			if result[pos+1].Max < number {
				result[pos+1].Max = number
			}
		}
	}

	r.cache = result

	return result
}

/*
CheckSet - analyze a single set
*/
func (r *Range) CheckSet(n Numbers) []int {
	var normal []int

	for pos_, number := range n.Set {
		pos := 1 + pos_
		if r.cache[pos].Min > number || r.cache[pos].Max < number {
			normal = append(normal, pos_)
		}
	}

	return normal
}
