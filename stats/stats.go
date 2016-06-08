package stats

import (
//    "fmt"
    "math"
    "sort"
)


func Sum(numbers []float64) (total float64) {
    for _, x := range numbers {
        total += x
    }
    return total
}

func Median(numbers []float64) float64 {
    middle := len(numbers) / 2
    result := numbers[middle]
    if len(numbers)%2 == 0 {
        result = (result + numbers[middle-1]) / 2
    }
    return result
}

func Mode(numbers []float64) (modes []float64) {
    frequencies := make(map[float64]int, len(numbers))
    highestFrequency := 0
    for _, x := range numbers {
        frequencies[x]++
        if frequencies[x] > highestFrequency {
            highestFrequency = frequencies[x]
        }
    }
    for x, frequency := range frequencies {
        if frequency == highestFrequency {
            modes = append(modes, x)
        }
    }
    if highestFrequency == 1 || len(modes) == len(numbers) {
        modes = modes[:0] // Or: modes = []float64{}
    }
    sort.Float64s(modes)
    return modes
}

func StdDev(numbers []float64, mean float64) float64 {
    total := 0.0
    for _, number := range numbers {
        total += math.Pow(number-mean, 2)
    }
    variance := total / float64(len(numbers)-1)
    return math.Sqrt(variance)
}

// Return low-pass filter output samples, given input samples,
 // 0.0 <= filterConstant <= 1.0
 
/*
func lowpass(unfiltered[]float64, filterConstant float64) (filtered []float64) {
   filtered = append(filtered, unfiltered[0])
   for i := 1; i < len(unfiltered); i++ {
       newfiltered := filterConstant * unfiltered[i] + (1-filterConstant) * filtered[i-1]
       filtered = append(filtered, newfiltered)
   }
   return filtered
}
*/