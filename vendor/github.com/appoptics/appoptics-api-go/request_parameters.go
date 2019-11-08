package appoptics

import (
	"net/http"
	"strconv"
)

// PaginationParameters holds pagination values
// https://docs.appoptics.com/api/?shell#request-parameters
type PaginationParameters struct {
	Offset  int
	Length  int
	Orderby string
	Sort    string
}

// AddToRequest mutates the provided http.Request with the PaginationParameters values
// Note that only valid values for Sort are "asc" and "desc" but the client does not enforce this.
func (rp *PaginationParameters) AddToRequest(req *http.Request) {
	if rp == nil {
		return
	}
	values := req.URL.Query()
	if rp.Offset > 0 {
		values.Add("offset", strconv.Itoa(rp.Offset))
	}

	if rp.Orderby != "" {
		values.Add("orderby", rp.Orderby)
	}

	if rp.Length > 0 {
		values.Add("length", strconv.Itoa(rp.Length))
	}

	if rp.Sort != "" {
		values.Add("sort", rp.Sort)
	}

	req.URL.RawQuery = values.Encode()
}
