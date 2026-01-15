package getdata

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetIssueDataEntry(url string) (ProteinIssue, error, []string, []string) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error here 1\n")
		return ProteinIssue{}, err, nil, nil
	}
	defer res.Body.Close()
	var PDB ProteinIssue
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&PDB)
	if err != nil {
		fmt.Printf("Error here 2\n")
		return ProteinIssue{}, err, nil, nil
	}
	var list_of_urls_polymers []string
	for _, polymer_id := range PDB.EntitiesInfo.PolymerID {
		list_of_urls_polymers = append(list_of_urls_polymers, fmt.Sprintf("https://data.rcsb.org/rest/v1/core/polymer_entity/7o52/%v", polymer_id))
	}
	var list_of_urls_non_polymers []string
	for _, non_polymer_id := range PDB.EntitiesInfo.NonPolymerID {
		list_of_urls_non_polymers = append(list_of_urls_non_polymers, fmt.Sprintf("https://data.rcsb.org/rest/v1/core/nonpolymer_entity/7o52/%v", non_polymer_id))
	}
	return PDB, nil, list_of_urls_polymers, list_of_urls_non_polymers
}

func GetDataForNonPolymers(url string) (NonPolymerIssue, error) {
	res, err := http.Get(url)
	if err != nil {
		return NonPolymerIssue{}, err
	}
	defer res.Body.Close()
	var NonPolymer NonPolymerIssue
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&NonPolymer)
	if err != nil {
		return NonPolymerIssue{}, err
	}
	return NonPolymer, nil
}

func GetDataForPolymers(url string) (PolymerIssue, error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error here 1\n")
		return PolymerIssue{}, err
	}
	defer res.Body.Close()
	var Polymer PolymerIssue
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&Polymer)
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Printf("Error here 2\n")
		return PolymerIssue{}, err
	}
	return Polymer, nil

}
func GetEntitiesURL(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error here 3\n")
		return nil, err
	}
	defer res.Body.Close()
	var PDB ProteinIssue
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&PDB)
	if err != nil {
		fmt.Printf("Error here 3\n")
		return nil, err
	}
	var list_of_urls []string
	for _, entity := range PDB.EntitiesInfo.PolymerID {
		list_of_urls = append(list_of_urls, fmt.Sprintf("curl https://data.rcsb.org/rest/v1/core/polymer_entity/7urx/%v", entity))
	}
	return list_of_urls, nil
}
