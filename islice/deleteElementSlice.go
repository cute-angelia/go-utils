package islice

// 不介意顺序

//a := []string{"A", "B", "C", "D", "E"}
//i := 2
//
//// Remove the element at index i from a.
//a[i] = a[len(a)-1] // Copy last element to index i.
//a[len(a)-1] = ""   // Erase last element (write zero value).
//a = a[:len(a)-1]   // Truncate slice.
//
//fmt.Println(a) // [A B E D]


// 介意顺序

//a := []string{"A", "B", "C", "D", "E"}
//i := 2
//
//// Remove the element at index i from a.
//copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index.
//a[len(a)-1] = ""     // Erase last element (write zero value).
//a = a[:len(a)-1]     // Truncate slice.
//
//fmt.Println(a) // [A B D E]