package raja_ongkir

type GetDeliveryRequest struct {
	Origin      string `form:"origin"`
	Destination string `form:"destination"`
	Weight      int    `form:"weight"`
	Courier     string `form:"courier"`
}
