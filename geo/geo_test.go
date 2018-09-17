package geo

import (
	"fmt"
	"testing"
)

func TestLatitudeLongitudeDistance(t *testing.T) {
	fmt.Println(">>>>>>>>", LatitudeLongitudeDistance(23.125178, 113.280637, 13.03659, 101.4925))
	fmt.Println(">>>>>>>>", LatitudeLongitudeDistance(23.125178, 113.280637, 36.204824, 138.252924))
}

func TestIpip(t *testing.T) {

	city, err := datx.NewCity("../data/17monipdb.datx") // For City Level IP Database

	if err != nil {
		t.Log(err)
	}

	t.Log(city.Find("8.8.8.8"))
	fmt.Println(city.Find("128.8.8.8"))
	fmt.Println(city.Find("255.255.255.255"))
	loc, err := city.FindLocation("27.190.252.103")
	if err == nil {
		fmt.Println(string(loc.ToJSON()))
		// Output:
		/*
		   {
		       "Country": "China",
		       "Province": "Hebei",
		       "City": "Tangshan",
		       "Organization": "",
		       "ISP": "ChinaTelecom",
		       "Latitude": "39.635113",
		       "Longitude": "118.175393",
		       "TimeZone": "Asia/Shanghai",
		       "TimeZone2": "UTC+8",
		       "CityCode": "130200",
		       "PhonePrefix": "86",
		       "CountryCode": "CN",
		       "ContinentCode": "AP",
		       "IDC": "",
		       "BaseStation": "",
		       "Anycast": false
		   }
		*/

	}
}
