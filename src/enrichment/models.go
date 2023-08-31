package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

type PatientAddress struct {
	City       string `json:"city"`
	Address1   string `json:"address1"`
	PostalCode string `json:"postalCode"`
	State      string `json:"state"`
}

type PatientEvent struct {
	Version       string `json:"version"`
	Source        string `json:"source"`
	EventType     string `json:"eventType"`
	CorrelationId string `json:"correlationId"`
	Details       struct {
		Command string       `json:"command"`
		Body    fhir.Patient `json:"entity"`
	} `json:"details"`
}

func buildPatientEvent(r *events.DynamoDBEventRecord) (*PatientEvent, error) {
	pe := &PatientEvent{
		Version:       "1.0",
		Source:        "PatientTable",
		EventType:     "PatientChange",
		CorrelationId: r.EventID,
		Details: struct {
			Command string       `json:"command"`
			Body    fhir.Patient `json:"entity"`
		}{
			Command: "PUT",
		},
	}
	fhirPatient := fhir.Patient{}
	humanName := fhir.HumanName{}

	for k, v := range r.Change.NewImage {
		if k == "address" {
			valueMap := v.Map()
			address := fhir.Address{}

			a := valueMap["address1"]
			b := valueMap["city"]
			c := valueMap["state"]
			d := valueMap["postalCode"]
			e := fhir.AddressUseHome

			city := b.String()
			addr := a.String()
			state := c.String()
			postal := d.String()
			address.Line = append(address.Line, addr)
			address.City = &city
			address.State = &state
			address.PostalCode = &postal
			address.Use = &e
			fhirPatient.Address = append(fhirPatient.Address, address)
		}

		if k == "birthDate" {
			birthDate := v.String()
			if birthDate == "" {
				return nil, fmt.Errorf("(birthDate) value is empty")
			}
			fhirPatient.BirthDate = &birthDate
		}

		if k == "id" {
			id := v.String()
			if id == "" {
				return nil, fmt.Errorf("(id) value is empty")
			}
			fhirPatient.Id = &id
		}

		if k == "firstName" {
			firstName := v.String()
			if firstName == "" {
				return nil, fmt.Errorf("(firstName) value is empty")
			}

			humanName.Given = []string{
				firstName,
			}
		}

		if k == "lastName" {
			lastName := v.String()
			if lastName == "" {
				return nil, fmt.Errorf("(lastName) value is empty")
			}

			humanName.Family = &lastName
		}
	}

	fhirPatient.Name = []fhir.HumanName{humanName}
	pe.Details.Body = fhirPatient
	return pe, nil
}
