package islice

// best way's you can use "github.com/emirpasic/gods/stacks/linkedliststack"
// but it's you mush push element first
// also you can use simple way down code

/*
Pop from queue
x, a = a[0], a[1:]

Pop from stack
x, a = a[len(a)-1], a[:len(a)-1]

Push
a = append(a, x)

demo:
	// 弹出一个
	ups	:= []userPackagePb.UserPackageModel{}

	for i:=0; i < 10; i++ {
		x := userPackagePb.UserPackageModel{}
		if len(ups) > 0 {
			x, ups = ups[0], ups[1:]
		}
	}
*/
