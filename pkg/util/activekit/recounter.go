package activekit

import "io"

type Recounter struct {
	Re  io.Reader
	sum int64
}

func (re *Recounter) Read(p []byte) (int, error) {
	n, err := re.Re.Read(p)
	re.sum += int64(n)
	return n, err
}
func (re *Recounter) Sum() int64 {
	return re.sum
}
