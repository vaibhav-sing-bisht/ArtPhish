package providers

import (
	"encoding/csv"
	"io"
	"net/http"
	"net/mail"
	"regexp"

	"github.com/gophish/gophish/models"
)

var (
	firstNameRegex = regexp.MustCompile(`(?i)first[\s_-]*name`)
	lastNameRegex  = regexp.MustCompile(`(?i)last[\s_-]*name`)
	emailRegex     = regexp.MustCompile(`(?i)email`)
	positionRegex  = regexp.MustCompile(`(?i)position`)
)

// ParseCSV contains the logic to parse the user provided csv file containing Target entries
func ParseCSV(r *http.Request) ([]models.Target, error) {

	mr, err := r.MultipartReader()
	ts := []models.Target{}
	if err != nil {
		return ts, err
	}
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		// Skip the "submit" part
		if part.FileName() == "" {
			continue
		}
		defer part.Close()
		reader := csv.NewReader(part)
		reader.TrimLeadingSpace = true
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		fi := -1
		li := -1
		ei := -1
		pi := -1
		fn := ""
		ln := ""
		ea := ""
		ps := ""
		for i, v := range record {
			switch {
			case firstNameRegex.MatchString(v):
				fi = i
			case lastNameRegex.MatchString(v):
				li = i
			case emailRegex.MatchString(v):
				ei = i
			case positionRegex.MatchString(v):
				pi = i
			}
		}
		if fi == -1 && li == -1 && ei == -1 && pi == -1 {
			continue
		}
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if fi != -1 && len(record) > fi {
				fn = record[fi]
			}
			if li != -1 && len(record) > li {
				ln = record[li]
			}
			if ei != -1 && len(record) > ei {
				csvEmail, err := mail.ParseAddress(record[ei])
				if err != nil {
					continue
				}
				ea = csvEmail.Address
			}
			if pi != -1 && len(record) > pi {
				ps = record[pi]
			}
			t := models.Target{
				BaseRecipient: models.BaseRecipient{
					FirstName: fn,
					LastName:  ln,
					Email:     ea,
					Position:  ps,
				},
			}
			ts = append(ts, t)
		}
	}
	return ts, nil
}
