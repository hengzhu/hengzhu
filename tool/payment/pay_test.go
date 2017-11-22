package payment

import (
	"testing"
)

func TestBbnpay(t *testing.T) {
	c := BbnPayConfig{
		Key:   "2b80219bf2592dc7dc72751de4fffa6e",
		AppId: "1032017051111958",
	}
	p := NewBbnPay(c)
	i := BbnPayPlaceOrder{
		Money:     1,
		GoodsId:   153,
		NotifyUrl: "123",
		PcorderId: "order_1",
		PcuserId:  "1",
		GoodsName: "test",
	}
	payInfo, err := p.PlaceOrder(&i)
	t.Log(err, payInfo)
}
