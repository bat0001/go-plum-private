package stylefile

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"
	// "strconv"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"go-PlumAIO/src/bot/utils"
)

type apiUsersMeResponse struct {
	BagID string `json:"bagId"`
}

type apiCommerceV1Products struct {
	BreadCrumbs []struct {
		Text string `json:"text"`
		Link string `json:"link"`
	} `json:"breadCrumbs"`
	ImageGroups []struct {
		Images []struct {
			URL string `json:"url"`
		} `json:"images"`
	} `json:"imageGroups"`
	Sizes []struct {
		SizeID          string `jsons:"sizeId"`
		SizeDescription string `json:"sizeDescription"`
		Scale           string `json:"scale"`
		Variants        []struct {
			MerchantID     int    `json:"merchantId"`
			FormattedPrice string `json:"formattedPrice"`
		} `json:"variants"`
	} `json:"sizes"`
}

type apiCommerceV1BagsPayload struct {
	MerchantID       int    `json:"merchantId"`
	ProductID        string `json:"productId"`
	Quantity         int    `json:"quantity"`
	Scale            string `json:"scale"`
	Size             string `json:"size"`
	CustomAttributes string `json:"customAttributes"`
}

type apiCommerceV1BagsResponse struct {
	BagSummary struct {
		GrandTotal float64 `json:"grandTotal"`
	} `json:"BagSummary"`
}

type apiCheckoutV1OrdersPayload struct {
	BagID            string `json:"bagId"`
	GuestUserEmail   string `json:"guestUserEmail"`
	UsePaymentIntent bool   `json:"usePaymentIntent"`
}

type apiCheckoutV1OrdersResponse struct {
	ID int `json:"id"`
}

type apiCheckoutV1OrdersResponse2 struct {
	CheckoutOrder struct {
		GrandTotal      float64 `json:"grandTotal"`
		PaymentIntentID string  `json:"paymentIntentId"`
	}
	ShippingOptions []struct {
		Price            float64 `json:"price"`
		FormattedPrice   string  `json:"formattedPrice"`
		ShippingCostType int     `json:"shippingCostType"`
		ShippingService  struct {
			Description              string  `json:"description"`
			ID                       int     `json:"id"`
			Name                     string  `json:"name"`
			Type                     string  `json:"type"`
			MinEstimatedDeliveryHour float64 `json:"minEstimatedDeliveryHour"`
			MaxEstimatedDeliveryHour float64 `json:"maxEstimatedDeliveryHour"`
		} `json:"shippingService"`
	} `json:"shippingOptions"`
}

type apiCheckoutV1OrderChargesResponse struct {
	RedirectURL string `json:"redirectUrl"`
}

func (t *Task) StylefileStart() {
	var err error

	if err = utils.CheckProfile(t.Email, t.FirstName, t.LastName, t.PhoneNumber, t.Address1, t.City, t.Postcode, t.Country); err != nil {
		t.HandleError(err)
		return
	}

	t.Email, t.FirstName, t.LastName, t.PhoneNumber, t.Address1 = utils.WrapProfile(t.Email, t.FirstName, t.LastName, t.PhoneNumber, t.Address1)

	t.Client, err = t.NewClient()

	if t.HandleError(err) {
		return
	}

	t.SetupClient()

	t.StylefileLogin()
}

func (t *Task) Create_StylefileSession() {

	t.Warn("Creating session")

	for {

		u := "https://www.stylefile.fr/"

		req, err := http.NewRequest("GET", u, nil)

		if err != nil {
			t.Error("Error creating session - 0 - %v", err.Error())
			continue
		}

		resp, err := t.Client.Do(req)

		if err != nil {
			t.Error("Error creating session - 1 - %v", err.Error())
			t.Rotate()
			continue
		}

		// b, err := ioutil.ReadAll(resp.Body)
		// resp.Body.Close()

		// if err != nil {
		// 	t.Error("Error creating session - 2 - %v", err.Error())
		// 	t.Sleep()
		// 	continue
		// }

		switch resp.StatusCode {
		case 200:

			t.Info("Success creating session - [%v]", resp.StatusCode)
			t.StylefileLogin()
			return

		case 400:
			t.Error("Error creating session - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		case 403:
			t.Error("Error creating session - Access Denied [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 429:
			t.Error("Error creating session - Rate limited [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 500, 501, 502, 503:
			t.Error("Error creating session - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		default:
			t.Error("Error creating session - Unhandled error [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		}
	}
}

func (t * Task) findToken() string {

	t.Warn("Researching token ...")
	url := "https://www.stylefile.fr/login"
	
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		t.Error("Error finding login token - 0 - %v", err.Error())

	}

	resp, err := t.Client.Do(req)

	if err != nil {
		fmt.Println("Error, reading the response", err)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}

	s := doc.Find("[name='csrf_token']")
	token, _ := s.Attr("value")
	return token

}

func (t *Task) StylefileLogin() {

	for {
		t.Warn("Login Started ...")

		token := t.findToken()
		t.Info("Token find ...")

		var data = strings.NewReader("dwfrm_login_username_d0kzedhdvzoa=lmaojungle@gmail.com&dwfrm_login_password_d0wrwezydcau=Baptiste2003&dwfrm_login_login=Inscrire&csrf_token=" + token)
		req, err := http.NewRequest("POST", "https://www.stylefile.fr/loginform", data)
		
		if err != nil {
			t.Error("Error Login - 1 - %v", err.Error())
			continue
		}

		req.Header.Set("authority", "www.stylefile.fr")
		req.Header.Set("cache-control", "max-age=0")
		req.Header.Set("upgrade-insecure-requests", "1")
		req.Header.Set("origin", "https://www.stylefile.fr")
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
		req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Set("sec-gpc", "1")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("sec-fetch-mode", "navigate")
		req.Header.Set("sec-fetch-user", "?1")
		req.Header.Set("sec-fetch-dest", "document")
		req.Header.Set("referer", "https://www.stylefile.fr/loginform")
		req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
		
		resp, err := t.Client.Do(req)

		if err != nil {
			t.Error("Error Login - 2 - %v", err.Error())
			t.Rotate()
			continue
		}

		switch resp.StatusCode {
		case 200:
			t.Info("Success Login - [%v]", resp.StatusCode)
			t.AddToCart()
			return

		case 400:
			t.Error("Error Login - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		case 403:
			t.Error("Error Login - Access Denied [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 429:
			t.Error("Error Login - Rate limited [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 500, 501, 502, 503:
			t.Error("Error Login - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		default:
			t.Error("Error Login - Unhandled error [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		}
	return
	}
}

func (t *Task) AddToCart() {

	t.Warn("Carting %v", t.ProductName)

	for {

		// url :=  t.URL
		var data = strings.NewReader(`Quantity=1&cartAction=add&pid=KK15D0000900250&css-tabs=on`)
		req, err := http.NewRequest("POST", "https://www.stylefile.fr/on/demandware.store/Sites-STF-FR-Site/fr_FR/Cart-AddProduct?format=ajax", data)
	
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("authority", "www.stylefile.fr")
		req.Header.Set("accept", "*/*")
		req.Header.Set("x-requested-with", "XMLHttpRequest")
		req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
		req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("sec-gpc", "1")
		req.Header.Set("origin", "https://www.stylefile.fr")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("referer", "https://www.stylefile.fr/lacoste-sweat-a-capuche-bleu-boys/P.KK15D00009002.html")
		req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
		resp, err := t.Client.Do(req)
		
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()
		
		if err != nil {
			fmt.Println(err)
		}

		if t.Method == "PP" {
			fmt.Println("PP PAYMENT")
		}
		if t.Method == "CC" {
			fmt.Println("CC PAYMENT")
		}
		return

		// switch resp.StatusCode {
		// case 200:
		// 	body := new(apiCommerceV1BagsResponse)

		// 	if err := json.Unmarshal(b, &body); err != nil {
		// 		t.Error("Error adding to cart - 4 - %v", err.Error())
		// 		t.SleepAndRotate()
		// 		continue
		// 	}

		// 	if body.BagSummary.GrandTotal > 0 {
		// 		t.Info("%v added to cart", t.ProductName)
		// 		t.SubmitGuest()
		// 		return
		// 	} else {
		// 		t.Error("Product OOS")
		// 		// t.CheckStock()
		// 		return
		// 	}

		// case 400:
		// 	t.Error("Error adding to cart - Product OOS [%v]", resp.StatusCode)
		// 	t.Sleep()
		// 	continue
		// case 403:
		// 	t.Error("Error adding to cart - Access Denied [%v]", resp.StatusCode)
		// 	t.SleepAndRotate()
		// 	continue
		// case 429:
		// 	t.Error("Error adding to cart - Rate limited [%v]", resp.StatusCode)
		// 	t.SleepAndRotate()
		// 	continue
		// case 500, 501, 502, 503:
		// 	t.Error("Error adding to cart - Site error [%v]", resp.StatusCode)
		// 	t.Sleep()
		// 	continue
		// default:
		// 	t.Error("Error adding to cart - Unhandled error [%v]", resp.StatusCode)
		// 	t.SleepAndRotate()
		// 	continue
		// }
	}
}

// func (t *Task) SubmitShipping() {

// 	t.Warn("Setting shipping")

// 	for {

// 		u := fmt.Sprintf("https://www.ambushdesign.com/api/checkout/v1/orders/%v", t.OrderID)

// 		p := map[string]interface{}{
// 			"shippingAddress": map[string]interface{}{
// 				"firstName": t.FirstName,
// 				"lastName":  t.LastName,
// 				"country": map[string]string{
// 					"name": utils.GetFullCountry(t.Country),
// 					"id":   t.CountryID,
// 				},
// 				"addressLine1": t.Address1,
// 				"addressLine2": t.Address2,
// 				"addressLine3": "",
// 				"city": map[string]string{
// 					"name": t.City,
// 				},
// 				"state": map[string]string{
// 					"name": t.State,
// 				},
// 				"zipCode":   t.Postcode,
// 				"phone":     t.PhoneNumber,
// 				"vatNumber": "",
// 			},
// 			"billingAddress": map[string]interface{}{
// 				"firstName": t.FirstName,
// 				"lastName":  t.LastName,
// 				"country": map[string]string{
// 					"name": utils.GetFullCountry(t.Country),
// 					"id":   t.CountryID,
// 				},
// 				"addressLine1": t.Address1,
// 				"addressLine2": t.Address2,
// 				"addressLine3": "",
// 				"city": map[string]string{
// 					"name": t.City,
// 				},
// 				"state": map[string]string{
// 					"name": t.State,
// 				},
// 				"zipCode":   t.Postcode,
// 				"phone":     t.PhoneNumber,
// 				"vatNumber": "",
// 			},
// 		}

// 		payload, err := json.Marshal(&p)

// 		if err != nil {
// 			t.Error("Error setting shipping - 0 - %v", err.Error())
// 			t.Sleep()
// 			continue
// 		}

// 		req, err := http.NewRequest("PATCH", u, bytes.NewReader(payload))

// 		if err != nil {
// 			t.Error("Error setting shipping - 1 - %v", err.Error())
// 			continue
// 		}

// 		req.Header.Set("Accept", "application/json, text/plain, */*")
// 		req.Header.Set("FF-Country", t.Country)
// 		req.Header.Set("FF-Currency", t.Currency)
// 		req.Header.Set("Accept-Language", "en-US")
// 		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
// 		req.Header.Set("Content-Type", "application/json")
// 		req.Header.Set("Origin", "https://www.ambushdesign.com")
// 		req.Header.Set("Sec-Fetch-Site", "same-origin")
// 		req.Header.Set("Sec-Fetch-Mode", "cors")
// 		req.Header.Set("Sec-Fetch-Dest", "empty")
// 		req.Header.Set("Referer", "https://www.ambushdesign.com/")

// 		resp, err := t.Client.Do(req)

// 		if err != nil {
// 			t.Error("Error setting shipping - 2 - %v", err.Error())
// 			t.Rotate()
// 			continue
// 		}

// 		b, err := ioutil.ReadAll(resp.Body)
// 		resp.Body.Close()

// 		if err != nil {
// 			t.Error("Error setting shipping - 3 - %v", err.Error())
// 			t.Sleep()
// 			continue
// 		}

// 		switch resp.StatusCode {
// 		case 200:
// 			body := new(apiCheckoutV1OrdersResponse2)

// 			if err := json.Unmarshal(b, &body); err != nil {
// 				t.Error("Error setting shipping - 4 - %v", err.Error())
// 				t.SleepAndRotate()
// 				continue
// 			}

// 			t.ShippingPrice = int(body.ShippingOptions[0].Price)
// 			t.ShippingFormattedPrice = body.ShippingOptions[0].FormattedPrice
// 			t.ShippingCostType = body.ShippingOptions[0].ShippingCostType
// 			t.ShippingDescription = body.ShippingOptions[0].ShippingService.Description
// 			t.ShippingID = body.ShippingOptions[0].ShippingService.ID
// 			t.ShippingName = body.ShippingOptions[0].ShippingService.Name
// 			t.ShippingType = body.ShippingOptions[0].ShippingService.Type
// 			t.MinEstimatedDeliveryHour = int(body.ShippingOptions[0].ShippingService.MinEstimatedDeliveryHour)
// 			t.MaxEstimatedDeliveryHour = int(body.ShippingOptions[0].ShippingService.MaxEstimatedDeliveryHour)
// 			t.GrandTotal = int(body.CheckoutOrder.GrandTotal)
// 			t.PaymentIntentID = body.CheckoutOrder.PaymentIntentID

// 			t.Info("Shipping set")
// 			t.SubmitDelivery()
// 			return

// 		case 400:
// 			t.Error("Error setting shipping - Site error [%v]", resp.StatusCode)
// 			t.Sleep()
// 			continue
// 		case 403:
// 			t.Error("Error setting shipping - Access Denied [%v]", resp.StatusCode)
// 			t.SleepAndRotate()
// 			continue
// 		case 429:
// 			t.Error("Error setting shipping - Rate limited [%v]", resp.StatusCode)
// 			t.SleepAndRotate()
// 			continue
// 		case 500, 501, 502, 503:
// 			t.Error("Error setting shipping - Site error [%v]", resp.StatusCode)
// 			t.Sleep()
// 			continue
// 		default:
// 			t.Error("Error setting shipping - Unhandled error [%v]", resp.StatusCode)
// 			t.SleepAndRotate()
// 			continue
// 		}
// 	}
// }

// func (t *Task) SubmitDelivery() {

// 	t.Warn("Setting delivery")

// 	for {

// 		u := fmt.Sprintf("https://www.ambushdesign.com/api/checkout/v1/orders/%v", t.OrderID)

// 		p := map[string]interface{}{
// 			"shippingOption": map[string]interface{}{
// 				"discount":         0,
// 				"merchants":        []int{t.MerchantID},
// 				"price":            t.ShippingPrice,
// 				"formattedPrice":   t.ShippingFormattedPrice,
// 				"shippingCostType": t.ShippingCostType,
// 				"shippingService": map[string]interface{}{
// 					"description":              t.ShippingDescription,
// 					"id":                       t.ShippingID,
// 					"name":                     t.ShippingName,
// 					"type":                     t.ShippingType,
// 					"minEstimatedDeliveryHour": t.MinEstimatedDeliveryHour,
// 					"maxEstimatedDeliveryHour": t.MaxEstimatedDeliveryHour,
// 					"trackingCodes":            []string{},
// 				},
// 				"shippingWithoutCapped": 0,
// 				"baseFlatRate":          0,
// 			},
// 		}

// 		payload, err := json.Marshal(&p)

// 		if err != nil {
// 			t.Error("Error setting delivery - 0 - %v", err.Error())
// 			t.Sleep()
// 			continue
// 		}

// 		req, err := http.NewRequest("PATCH", u, bytes.NewReader(payload))

// 		if err != nil {
// 			t.Error("Error setting delivery - 1 - %v", err.Error())
// 			continue
// 		}

// 		req.Header.Set("Accept", "application/json, text/plain, */*")
// 		req.Header.Set("FF-Country", t.Country)
// 		req.Header.Set("FF-Currency", t.Currency)
// 		req.Header.Set("Accept-Language", "en-US")
// 		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
// 		req.Header.Set("Content-Type", "application/json")
// 		req.Header.Set("Origin", "https://www.ambushdesign.com")
// 		req.Header.Set("Sec-Fetch-Site", "same-origin")
// 		req.Header.Set("Sec-Fetch-Mode", "cors")
// 		req.Header.Set("Sec-Fetch-Dest", "empty")
// 		req.Header.Set("Referer", "https://www.ambushdesign.com/")

// 		resp, err := t.Client.Do(req)

// 		if err != nil {
// 			t.Error("Error setting delivery - 2 - %v", err.Error())
// 			t.Rotate()
// 			continue
// 		}

// 		resp.Body.Close()

// 		switch resp.StatusCode {
// 		case 200:
// 			t.Info("Delivery set")
// 			t.SubmitPayment()
// 			return
// 		case 400:
// 			t.Error("Error setting delivery - Site error [%v]", resp.StatusCode)
// 			t.Sleep()
// 			continue
// 		case 403:
// 			t.Error("Error setting delivery - Access Denied [%v]", resp.StatusCode)
// 			t.SleepAndRotate()
// 			continue
// 		case 429:
// 			t.Error("Error setting delivery - Rate limited [%v]", resp.StatusCode)
// 			t.SleepAndRotate()
// 			continue
// 		case 500, 501, 502, 503:
// 			t.Error("Error setting delivery - Site error [%v]", resp.StatusCode)
// 			t.Sleep()
// 			continue
// 		default:
// 			t.Error("Error setting delivery - Unhandled error [%v]", resp.StatusCode)
// 			t.SleepAndRotate()
// 			continue
// 		}
// 	}
// }

// func (t *Task) SubmitPayment() {

// 	t.Warn("Setting payment")

// 	ctx := utils.GetCookie(&t.CookieJar, "ctx")
// 	t.CTX, _ = utils.Extract(ctx, "%3a", "%2c")
// 	i, _ := strconv.Atoi(t.CTX)

// 	for {

// 		u := fmt.Sprintf("https://www.ambushdesign.com/api/payment/v1/intents/%v/instruments", t.PaymentIntentID)

// 		p := map[string]interface{}{
// 			"method":      "PayPal",
// 			"option":      "PayPalExp",
// 			"createToken": false,
// 			"payer": map[string]interface{}{
// 				"id":        i,
// 				"firstName": t.FirstName,
// 				"lastName":  t.LastName,
// 				"email":     t.Email,
// 				"birthDate": nil,
// 				"address": map[string]interface{}{
// 					"city": map[string]interface{}{
// 						"countryId": t.CountryID,
// 						"id":        0,
// 						"name":      t.City,
// 					},
// 					"country": map[string]interface{}{
// 						"alpha2Code":  t.Country,
// 						"alpha3Code":  t.Country,
// 						"culture":     "it-IT",
// 						"id":          t.CountryID,
// 						"name":        utils.GetFullCountry(t.Country),
// 						"nativeName":  utils.GetFullCountry(t.Country),
// 						"region":      "Europe",
// 						"regionId":    0,
// 						"continentId": 3,
// 					},
// 					"id":       "00000000-0000-0000-0000-000000000000",
// 					"lastName": t.LastName,
// 					"state": map[string]interface{}{
// 						"countryId": 0,
// 						"id":        0,
// 						"code":      t.State,
// 						"name":      t.State,
// 					},
// 					"userId":                   0,
// 					"isDefaultBillingAddress":  false,
// 					"isDefaultShippingAddress": false,
// 					"addressLine1":             t.Address1,
// 					"addressLine2":             t.Address2,
// 					"firstName":                t.FirstName,
// 					"phone":                    t.PhoneNumber,
// 					"vatNumber":                "",
// 					"zipCode":                  t.Postcode,
// 				},
// 			},
// 			"amounts": []map[string]interface{}{{
// 				"value": t.GrandTotal,
// 			}},
// 			"data": map[string]interface{}{},
// 		}

// 		payload, err := json.Marshal(&p)

// 		if err != nil {
// 			t.Error("Error setting payment - 0 - %v", err.Error())
// 			t.Sleep()
// 			continue
// 		}

// 		req, err := http.NewRequest("POST", u, bytes.NewReader(payload))

// 		if err != nil {
// 			t.Error("Error setting payment - 1 - %v", err.Error())
// 			continue
// 		}

// 		req.Header.Set("Accept", "application/json, text/plain, */*")
// 		req.Header.Set("FF-Country", t.Country)
// 		req.Header.Set("FF-Currency", t.Currency)
// 		req.Header.Set("Accept-Language", "en-US")
// 		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
// 		req.Header.Set("Content-Type", "application/json")
// 		req.Header.Set("Origin", "https://www.ambushdesign.com")
// 		req.Header.Set("Sec-Fetch-Site", "same-origin")
// 		req.Header.Set("Sec-Fetch-Mode", "cors")
// 		req.Header.Set("Sec-Fetch-Dest", "empty")
// 		req.Header.Set("Referer", "https://www.ambushdesign.com/")

// 		resp, err := t.Client.Do(req)

// 		if err != nil {
// 			t.Error("Error setting payment - 2 - %v", err.Error())
// 			t.Rotate()
// 			continue
// 		}

// 		resp.Body.Close()

// 		switch resp.StatusCode {
// 		case 201:
// 			t.Info("Payment set")
// 			t.CheckCharge()
// 			return
// 		case 400:
// 			t.Error("Error setting payment - Site error [%v]", resp.StatusCode)
// 			t.Sleep()
// 			continue
// 		case 403:
// 			t.Error("Error setting payment - Access Denied [%v]", resp.StatusCode)
// 			t.SleepAndRotate()
// 			continue
// 		case 429:
// 			t.Error("Error setting payment - Rate limited [%v]", resp.StatusCode)
// 			t.SleepAndRotate()
// 			continue
// 		case 500, 501, 502, 503:
// 			t.Error("Error setting payment - Site error [%v]", resp.StatusCode)
// 			t.Sleep()
// 			continue
// 		default:
// 			t.Error("Error setting payment - Unhandled error [%v]", resp.StatusCode)
// 			t.SleepAndRotate()
// 			continue
// 		}
// 	}
// }



// func (t *Task) SubmitPayPal() {

// 	t.Warn("Submitting PayPal")

// 	for {

// 		req, err := http.NewRequest("GET", t.RedirectURL, nil)

// 		if err != nil {
// 			t.Error("Error submitting PayPal - 1 - %v", err.Error())
// 			continue
// 		}

// 		req.Header.Set("Accept", "application/json, text/plain, */*")
// 		req.Header.Set("FF-Country", t.Country)
// 		req.Header.Set("FF-Currency", t.Currency)
// 		req.Header.Set("Accept-Language", "en-US")
// 		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
// 		req.Header.Set("Content-Type", "application/json")
// 		req.Header.Set("Origin", "https://www.ambushdesign.com")
// 		req.Header.Set("Sec-Fetch-Site", "same-origin")
// 		req.Header.Set("Sec-Fetch-Mode", "cors")
// 		req.Header.Set("Sec-Fetch-Dest", "empty")
// 		req.Header.Set("Referer", "https://www.ambushdesign.com/")

// 		t.SetAllowRedirects(false)

// 		resp, err := t.Client.Do(req)

// 		if err != nil {
// 			t.Error("Error submitting PayPal - 2 - %v", err.Error())
// 			t.Rotate()
// 			continue
// 		}

// 		switch resp.StatusCode {
// 		case 302:

// 			if strings.Contains(resp.Header.Get("Location"), "paypal") {

// 				t.PayPalURL = resp.Header.Get("Location")
// 				return
// 			}

// 			return
// 		case 400:
// 			t.Error("Error submitting PayPal - Site error [%v]", resp.StatusCode)
// 			t.Sleep()
// 			continue
// 		case 403:
// 			t.Error("Error submitting PayPal - Access Denied [%v]", resp.StatusCode)
// 			t.SleepAndRotate()
// 			continue
// 		case 429:
// 			t.Error("Error submitting PayPal - Rate limited [%v]", resp.StatusCode)
// 			t.SleepAndRotate()
// 			continue
// 		case 500, 501, 502, 503:
// 			t.Error("Error submitting PayPal - Site error [%v]", resp.StatusCode)
// 			t.Sleep()
// 			continue
// 		default:
// 			t.Error("Error submitting PayPal - Unhandled error [%v]", resp.StatusCode)
// 			t.SleepAndRotate()
// 			continue
// 		}
// 	}
// }
