package model

import "net/url"

type Address struct {
	//Userid   string `json:"user_id"`
	//Primary  bool   `json:"primary_ad"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Address3 string `json:"address3"`
	PIN      string `json:"pin"`
	Country  string `json:"country"`
	State    string `json:"state"`
	//Vendorid int    `json:"vendor_id"`
}

func (a *Address) Check() url.Values {
	err := url.Values{}
	if len(a.PIN) != 6 {
		err.Add("Pincode", "invalid pincode")
	}

	if len(a.Address1) < 3 {
		err.Add("Address1", "invalid address")
	}
	if len(a.Address2) < 3 {
		err.Add("Address2", "invalid address")
	}
	if len(a.Address3) < 3 {
		err.Add("Address3", "invalid address")
	}
	if len(a.Country) < 2 {
		err.Add("Country", "invalid Country")
	}
	if len(a.State) < 2 {
		err.Add("State", "invalid State")
	}

	return err

}
