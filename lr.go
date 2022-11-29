package main

import (
	"fmt"
	"sync"
	"time"
)

const n int = 20000000

func main() {
	a := []float64{}
	var wg0 sync.WaitGroup
	wg0.Add(1)
	go func() {
		defer wg0.Done()
		for i := 0; i < n-2; i++ {
			a = append(a, 1)
		}
		a = append(a, 0)
	}()
	wg0.Add(1)
	b := []float64{0}
	go func() {
		defer wg0.Done()
		for i := 1; i < n-1; i++ {
			b = append(b, 1)
		}
	}()
	wg0.Add(1)
	c := []float64{1}
	go func() {
		wg0.Done()
		for i := 1; i < n-1; i++ {
			c = append(c, 2.0-5.0/(float64(n-1)*float64(n-1)))
		}
		c = append(c, 1)
	}()
	wg0.Add(1)
	f := []float64{5}
	go func() {
		defer wg0.Done()
		for i := 1; i < n-1; i++ {
			f = append(f, 0)
		}
		f = append(f, -1)
	}()
	wg0.Wait()

	t1 := time.Now()

	alf := []float64{}
	ps := [n - 1]float64{}

	a1 := make([]float64, len(a))
	copy(a1, a)
	c1 := make([]float64, len(c))
	copy(c1, c)
	b1 := make([]float64, len(b))
	copy(b1, b)
	f1 := make([]float64, len(f))
	copy(f1, f)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		alf = alfa(a, c, b)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		ps = psi(a1, c1, b1)
	}()
	wg.Wait()

	bet := []float64{}
	et := [n - 1]float64{}

	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		bet = beta(a, c, alf, f)
	}()
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		et = eta(b, c1, ps, f1)
	}()
	wg1.Wait()

	x1 := [n / 2]float64{}
	x2 := [n / 2]float64{}

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		x2 = solve_ab(f, alf, bet, a, c)
	}()
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		x1 = solve_pe(f1, ps, et, b, c1)
	}()
	wg2.Wait()

	x := []float64{}
	x = append(x, x1[:]...)
	x = append(x, x2[:]...)

	t2 := time.Since(t1)
	fmt.Println(t2)
	_, acc := accuracy(a, b, c, f, x)
	fmt.Println(acc)
}

func alfa(a, c, b []float64) []float64 {
	ret := []float64{}
	ret = append(ret, -b[0]/c[0])
	for i := 1; i < n-1; i++ {
		ret = append(ret, -b[i]/(a[i-1]*ret[i-1]+c[i]))
	}
	return ret
}

func beta(a, c, alf, f []float64) []float64 {
	ret := []float64{}
	ret = append(ret, f[0]/c[0])
	for i := 1; i < n-1; i++ {
		ret = append(ret, (f[i]-a[i-1]*ret[i-1])/(a[i-1]*alf[i-1]+c[i]))
	}
	return ret
}

func solve_ab(f []float64, alf []float64, beta []float64, a []float64, c []float64) [n / 2]float64 {
	ret := [n / 2]float64{}
	ln := n / 2
	ret[ln-1] = (f[n-1] - a[n-2]*beta[n-2]) / (a[n-2]*alf[n-2] + c[n-1])
	for i := n - 2; i >= ln; i-- {
		ret[i-ln] = alf[i]*ret[i-ln+1] + beta[i]
	}
	return ret
}

func psi(a, c, b []float64) [n - 1]float64 {
	ret := [n - 1]float64{}
	ret[n-2] = -a[n-2] / c[n-1]
	for i := n - 3; i >= 0; i-- {
		ret[i] = -a[i] / (b[i+1]*ret[i+1] + c[i+1])
	}
	return ret
}

func eta(b []float64, c []float64, psi [n - 1]float64, f []float64) [n - 1]float64 {
	ret := [n - 1]float64{}
	ret[n-2] = f[n-1] / c[n-1]
	for i := n - 3; i >= 0; i-- {
		ret[i] = (f[i+1] - b[i+1]*ret[i+1]) / (b[i+1]*psi[i+1] + c[i+1])
	}
	return ret
}

func solve_pe(f []float64, psi [n - 1]float64, eta [n - 1]float64, b []float64, c []float64) [n / 2]float64 {
	ret := [n / 2]float64{}
	ln := n / 2
	ret[0] = (f[0] - b[0]*eta[0]) / (b[0]*psi[0] + c[0])
	for i := 1; i <= ln-1; i++ {
		ret[i] = psi[i-1]*ret[i-1] + eta[i-1]
	}
	return ret
}

func accuracy(a, b, c, f []float64, x []float64) ([n]float64, float64) {
	ret := [n]float64{}
	ret[0] = c[0]*x[0] + b[0]*x[1] - f[0]
	max_acc := ret[0]
	for i := 1; i < n-2; i++ {
		ret[i] = a[i-1]*x[i-1] + c[i]*x[i] + b[i]*x[i+1] - f[i]
		if ret[i] > max_acc {
			max_acc = ret[i]
		}
	}
	ret[n-1] = a[n-2]*x[n-2] + c[n-1]*x[n-1] - f[n-1]
	return ret, max_acc
}
