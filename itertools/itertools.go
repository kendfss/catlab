package itertools // catlab/itertools
// package main

import (
	"fmt"
    "crypto/rand"
    // "math/rand"
)

func Randex(length int) int {
    // bites := make([]byte, len(iterable))
    bites := make([]byte, 1)
    _, err := rand.Read(bites)
    if err!=nil {
        panic(fmt.Sprintf("Couldn't generate random number(s):\n\t%s", err))
    }
    // return rand.Int() % length
    // n, err := int(rand.Int(rand.Read, length))
    return int(bites[0]) % length
}



func Merge(receiver *[]string, giver []string) {
    for _, str := range giver {
        *receiver = append(*receiver, str)
    }
}

func Show(box []string) {
    // print elements from an arbitrary iterable
    for _, element := range box {
        fmt.Println(element)
    }
}

func Filter(f func(interface{})bool, iter []interface{}) []interface{} {
    out := make([]interface{}, 0)
    for _, elem := range iter {
        if f(elem) {
            out = append(out, elem)
        }
    }
    return out
}
