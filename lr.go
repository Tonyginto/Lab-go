package main

import (
	"fmt"
	"time"
)

const n int = 8192

func main() {
	a := []float64{}
	for i := 0; i < n-2; i++ {
		a = append(a, 1)
	}
	a = append(a, 0)
	b := []float64{0}
	for i := 1; i < n-1; i++ {
		b = append(b, 1)
	}
	c := []float64{1}
	for i := 1; i < n-1; i++ {
		c = append(c, 2.0-5.0/(float64(n-1)*float64(n-1)))
	}
	c = append(c, 1)
	f := []float64{5}
	for i := 1; i < n-1; i++ {
		f = append(f, 0)
	}
	f = append(f, -1)

	t1 := time.Now()

	out1 := make(chan []float64)
	out2 := make(chan [n - 1]float64)
	go func() {
		out1 <- alfa(a, c, b)
	}()
	go func() {
		out2 <- psi(a, c, b)
	}()
	alf := <-out1
	psi := <-out2

	out3 := make(chan []float64)
	out4 := make(chan [n - 1]float64)
	go func() {
		out3 <- beta(a, c, alf, f)
	}()
	go func() {
		out4 <- eta(b, c, psi, f)
	}()
	bet := <-out3
	eta := <-out4

	out5 := make(chan [n / 2]float64)
	out6 := make(chan [n / 2]float64)
	go func() {
		out5 <- solve_ab(f, alf, bet, a, c)
	}()
	go func() {
		out6 <- solve_pe(f, psi, eta, b, c)
	}()
	x1 := <-out6
	x2 := <-out5
	x := []float64{}
	x = append(x, x1[:]...)
	x = append(x, x2[:]...)

	t2 := time.Since(t1)
	fmt.Println(t2)
	_, acc := accuracy(a, b, c, f, x)
	fmt.Println("accurace = ", acc)
	fmt.Println("x = ", x)

	/*alf := alfa(a, c, b)
	//fmt.Println("alfa = ", alf)
	bet := beta(a, c, alf, f)
	//fmt.Println("beta = ", bet)
	x2 := solve_ab(f, alf, bet, a, c)*/

	/*psi := psi(a, c, b)
	//fmt.Println("psi = ", psi)
	eta := eta(b, c, psi, f)
	//fmt.Println("eta = ", eta)
	x1 := solve_pe(f, psi, eta, b, c)*/

	/*x := []float64{}
	x = append(x, x1[:]...)
	x = append(x, x2[:]...)
	acc := accuracy(a, b, c, f, x)
	fmt.Println("accurace = ", acc)
	fmt.Println("x1 = ", x1)
	fmt.Println("x2 = ", x2)
	fmt.Println("x = ", x)*/
}

func alfa(a []float64, c []float64, b []float64) []float64 {
	ret := []float64{}
	ret = append(ret, -b[0]/c[0])
	for i := 1; i < n-1; i++ {
		ret = append(ret, -b[i]/(a[i-1]*ret[i-1]+c[i]))
	}
	return ret
}

func beta(a []float64, c []float64, alf []float64, f []float64) []float64 {
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

func psi(a []float64, c []float64, b []float64) [n - 1]float64 {
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

func accuracy(a []float64, b []float64, c []float64, f []float64, x []float64) ([n]float64, float64) {
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
