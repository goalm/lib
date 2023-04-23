package utils

import "strconv"

type Date struct {
	Year  int
	Month int
}

func (this Date) DateValue() (res int) {
	res = this.Year*100 + this.Month
	return res
}

func (this Date) DateStr() (res string) {
	if this.Month >= 10 {
		res = strconv.Itoa(this.Month) + "/" + strconv.Itoa(this.Year)
	} else {
		res = "0" + strconv.Itoa(this.Month) + "/" + strconv.Itoa(this.Year)
	}

	return res
}

func (this Date) CalendarMth(t int) (res int) {
	num := (this.Month + t) % 12
	if num == 0 {
		res = 12
	} else {
		res = num
	}
	return
}

func (this Date) CalendarYr(t int) (res int) {
	if this.Month == 12 {
		res = this.Year + (t+11)/12
	} else {
		res = this.Year + (this.Month+t-1)/12
	}
	return
}

func (this Date) T(Year, Month int) int {
	res := (Year-this.Year)*12 + Month - this.Month
	return res
}

func (this Date) CalendarDate(i int) Date {
	res := Date{this.CalendarYr(i), this.CalendarMth(i)}
	return res
}

func Dur(x Date, y Date) int {
	res := (y.Year-x.Year)*12 + y.Month - x.Month
	return res
}
