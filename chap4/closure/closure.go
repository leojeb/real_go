package main

func a() (s []func()) {
	for i := 0; i < 2; i++ {
		println(i)
		s = append(s, func() {
			println(&i, i)
		})
	}
	return s

}

func b() (s []func()) {
	for i := 0; i < 2; i++ {
		var x = i
		s = append(s, func() {
			println(&x, x)
		})
	}
	return s
}
func main() {
	for _, f := range a() {
		f()
	}
	for _, f := range b() {
		f()
	}
	//var i = 1
	//chap1.Myprint("i", i)
}
