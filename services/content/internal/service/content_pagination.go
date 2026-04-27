package service

func pageValues(page, pageSize int32) (int, int) {
	p := int(page)
	if p <= 0 {
		p = 1
	}
	ps := int(pageSize)
	if ps <= 0 {
		ps = 20
	}
	if ps > 100 {
		ps = 100
	}
	return p, ps
}
