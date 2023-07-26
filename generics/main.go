package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/tour/pic"
	"golang.org/x/tour/tree"
	"golang.org/x/tour/wc"
)

type Number interface {
	int64 | float64
}

func main() {
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}

	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	fmt.Printf("Non-Generic Sums: %v and %v\n",
		SumInts(ints),
		SumFloats(floats))

	fmt.Printf("Generic Sums: %v and %v\n",
		SumIntsOrFloats[string, int64](ints),
		SumIntsOrFloats[string, float64](floats))

	fmt.Printf("Generic Sums, type parameters inferred: %v and %v\n",
		SumIntsOrFloats(ints),
		SumIntsOrFloats(floats))

	fmt.Printf("Generic Sums with Constraint: %v and %v\n",
		SumNumbers(ints),
		SumNumbers(floats))

	v, i := Sqrt(15)
	fmt.Printf("Sqrt(15): %v after %v iterations\n", v, i)

	fmt.Printf("Go runtime: %v", runtime.GOOS)

	pic.Show(Pic)
	wc.Test(WordCount)

	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Printf("%v, ", f())
	}
	fmt.Println()

	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
	fmt.Println()

	var p List[int]
	p.val = 2
	p.Append(12)
	p.Append(24)
	p.Append(36)
	p.Append(48)
	p.Append(60)
	fmt.Println(p.Length())
	fmt.Println(p.Get(3))
	fmt.Println(p.Get(6))

	t1 := tree.New(1)
	t2 := tree.New(1)
	t3 := tree.New(2)
	c1 := make(chan int)
	c2 := make(chan int)
	c3 := make(chan int)
	go Walk(t1, c1)
	go Walk(t2, c2)
	go Walk(t3, c3)
	fmt.Printf("t1 = t2: %v\n", Same(t1, t2))
	fmt.Printf("t1 = t3: %v\n", Same(t1, t3))

	testCrawl()
}

func testCrawl() {
	wg := &sync.WaitGroup{}
	urlMap := SafeMap{v: make(map[string]int)}
	counter := 0

	Crawl(&urlMap, "https://golang.org/", 4, fetcher, wg, counter)

	wg.Wait()
}

func SumNumbers[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func SumInts(m map[string]int64) int64 {
	var s int64
	for _, v := range m {
		s += v
	}
	return s
}

func SumFloats(m map[string]float64) float64 {
	var s float64
	for _, v := range m {
		s += v
	}
	return s
}

func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func Sqrt(x float64) (float64, int) {
	z := x / 2.0
	i := 0
	y := 0.0
	for {
		y = (z*z - x) / (2 * z)
		if math.Abs(y) < 1e-10 {
			return (z - y), i
		}
		z -= y
		i += 1
	}
}

func Pic(dx, dy int) [][]uint8 {
	img := make([][]uint8, dy)
	for y := 0; y < dy; y++ {
		slice := make([]uint8, dx)
		for x := 0; x < dx; x++ {
			slice[x] = uint8((x + y) / 2)
		}
		img[y] = slice
	}
	return img
}

func WordCount(s string) map[string]int {
	counts := make(map[string]int)
	words := strings.Fields(s)
	for _, word := range words {
		_, ok := counts[word]
		if ok {
			counts[word] += 1
		} else {
			counts[word] = 1
		}
	}
	return counts
}

func fibonacci() func() int {
	prev := 0
	curr := 1
	return func() int {
		oldprev := prev
		newcurr := prev + curr
		prev = curr
		curr = newcurr
		return oldprev
	}
}

type rot13Reader struct {
	r io.Reader
}

func (r rot13Reader) Read(b []byte) (int, error) {
	n, err := r.r.Read(b)
	for i := 0; i < n; i++ {
		var v byte
		if b[i] >= 65 && b[i] <= 90 {
			v = b[i] + 13
			if v > 90 {
				v = v - 90 + 64
			}
		} else if b[i] >= 97 && b[i] <= 122 {
			v = b[i] + 13
			if v > 122 {
				v = v - 122 + 96
			}
		} else {
			v = b[i]
		}
		b[i] = v
	}
	return n, err
}

type List[T any] struct {
	next *List[T]
	val  T
}

func (l *List[T]) Append(x T) {
	var next List[T]
	next.val = x
	tail := l.Tail()
	tail.next = &next
}

func (l *List[T]) Head() *List[T] {
	return l
}

func (l *List[T]) Tail() *List[T] {
	temp := l
	for ; temp.next != nil; temp = temp.next {
	}
	return temp
}

func (l *List[T]) Get(i int) *List[T] {
	temp := l
	for x := 0; x < i; x++ {
		if temp.next != nil {
			temp = temp.next
		} else {
			return nil
		}
	}
	return temp
}

func (l *List[T]) Length() int {
	count := 1
	for temp := l; temp.next != nil; temp = temp.next {
		count += 1
	}
	return count
}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if t.Left != nil {
		Walk(t.Left, ch)
	}

	ch <- t.Value

	if t.Right != nil {
		Walk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	c1 := make(chan int)
	c2 := make(chan int)
	go Walk(t1, c1)
	go Walk(t2, c2)
	for {
		v1, ok1 := <-c1
		v2, ok2 := <-c2

		if (!ok1 && ok2) || (ok1 && !ok2) {
			return false
		} else if v1 != v2 {
			return false
		}
		return true
	}
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type SafeMap struct {
	mu sync.Mutex
	v  map[string]int
}

func checkUrlExist(urlMap *SafeMap, url string) bool {
	urlMap.mu.Lock()
	defer urlMap.mu.Unlock()
	return urlMap.v[url] > 0
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(urlMap *SafeMap, url string, depth int, fetcher Fetcher, wg *sync.WaitGroup, counter int) {

	if depth <= 0 {
		return
	}

	if checkUrlExist(urlMap, url) {
		return
	}

	urlMap.mu.Lock()
	urlMap.v[url] = 1
	urlMap.mu.Unlock()

	body, urls, err := fetcher.Fetch(url)
	counter += 1

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q (count: %v)\n", url, body, counter)

	for _, u := range urls {
		wg.Add(1)
		go func(uu string) {
			defer wg.Done()
			Crawl(urlMap, uu, depth-1, fetcher, wg, counter)
		}(u)
	}
	return
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
