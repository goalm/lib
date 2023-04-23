package asset

import (
	"github.com/goalm/lib/utils"
	"math"
)

type Bond struct {
	Mp
	Initialized  bool
	Redeemed     bool
	PropHeld     float64
	PrevPropHeld float64
	PrevMv       float64
	PrevAbv      float64
	Mv           float64
	Abv          float64
}

type Mp struct {
	Seg           int     `csv:"SEGMENT"`
	AssetType     string  `csv:"ASSET_TYPE"`
	Economy       string  `csv:"ECONOMY"`
	SpCode        int     `csv:"SP_CODE"`
	Pool          int     `csv:"POOL" `
	Cat           int     `csv:"CATEGORY" `
	RedempYear    int     `csv:"REDEMP_YEAR" `
	RedempMonth   int     `csv:"REDEMP_MONTH" `
	RedempAmt     float64 `csv:"REDEMP_AMT"`
	CouponPc      float64 `csv:"COUPON_PC"`
	CouponFreq    int     `csv:"COUPON_FREQ"`
	AssetScalar   float64 `csv:"ASSET_SCALAR"`
	BasisFlag     int     `csv:"BASIS_FLAG"`
	InitMvUsed    bool    `csv:"INIT_MV_USED"`
	InitMv        float64 `csv:"I_MV"`
	InitMktSpdPc  float64 `csv:"I_MARKET_SPREAD_PC"`
	InitBvUsed    bool    `csv:"INIT_BV_USED"`
	AmortType     int     `csv:"AMORT_TYPE"`
	InitAcciUsed  bool    `csv:"INIT_ACCI_USED"`
	InitAbv       float64 `csv:"I_ABV"`
	InitAcci      float64 `csv:"I_ACCI"`
	AmortRatePc   float64 `csv:"AMORT_RATE_PC"`
	SpreadBand    int     `csv:"SPREAD_BAND"`
	FaceValue     float64 `csv:"FACE_VALUE"`
	NotUsedString string  `csv:"-"`
}

// Bond rolling forward to the next period
func BondRolls(this *Bond, start, end utils.Date) (cf []float64) {
	redempT := this.RedempT(start)
	n := utils.Dur(start, end)
	this.PrevAbv = this.Abv

	cf = make([]float64, n+1)
	for i := 1; i <= n; i++ {
		//Coupon
		if this.CouponFreq != 0 {
			if (redempT-i)%(12/this.CouponFreq) == 0 {
				cf[i] += this.RedempAmt * this.CouponPc / 100 / float64(this.CouponFreq)
			}
		}
		// Redemption
		if i == redempT {
			cf[i] += this.RedempAmt
			this.Abv = 0
			this.Redeemed = true
			break
		}
	}
	if this.Redeemed == false {
		this.Val(end)
	}
	//fmt.Println(redempT)
	return cf
}

func (bond *Bond) Val(valDate utils.Date) {
	bond.Abv = bond.AmortBookValue(valDate)
}

func (bond Mp) RedempT(start utils.Date) int {
	t := (bond.RedempYear-start.Year)*12 + (bond.RedempMonth - start.Month)
	return t
}

func (bond Mp) AmortBookValue(start utils.Date) float64 {
	pv := 0.00
	mat := bond.RedempT(start)
	for i := 1; i <= mat; i++ {
		if bond.CouponFreq != 0 {
			if (mat-i)%(12/bond.CouponFreq) == 0 {
				coupon := bond.RedempAmt * bond.CouponPc / 100 / float64(bond.CouponFreq)
				pv += coupon / math.Pow(1+bond.AmortRatePc/100, float64(i)/12)
			}

		}
	}
	pv += bond.RedempAmt / math.Pow(1+bond.AmortRatePc/100, float64(mat)/12)
	return pv
}
