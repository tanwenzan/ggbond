package main

func inc10(i *int) {
	*i += 10
}

func double(arr *[]int) {
	for i, v := range *arr {
		(*arr)[i] = v * 2
	}
}
