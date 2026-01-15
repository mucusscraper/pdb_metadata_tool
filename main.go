package main

import (
	"bufio"
	"fmt"
	"os"

	cleaninput "github.com/mucusscraper/pdb_metadata_tool/internal/clean_input"
	getdata "github.com/mucusscraper/pdb_metadata_tool/internal/get_data"
)

const issueURL = "https://data.rcsb.org/rest/v1/core/entry/"

func main() {
	fmt.Printf("Welcome to PDBMetaDataTool!\n\n")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("PDBMDT > ")
		if scanner.Scan() {
			all_string := scanner.Text()
			ids := cleaninput.CleanInput(all_string)
			for _, id := range ids {
				url := fmt.Sprintf("%v%v", issueURL, id)
				res, err, list_of_urls_polymer, list_of_urls_non_polymers := getdata.GetIssueDataEntry(url)
				if err != nil {
					fmt.Printf("Error getting data\n")
					return
				}
				fmt.Printf("#######################GENERAL DATA#######################\n")
				fmt.Printf("RCSB ID: %v\n", res.ID)
				fmt.Printf("DEPOSIT DATE: %v\n", res.AccessInfo.DepositDate)
				fmt.Printf("DOI: %v\n", res.ArticleInfo.DOI)
				fmt.Printf("PAPER TITLE: %v\n", res.ArticleInfo.Title)
				for _, method := range res.ExptlInfo {
					fmt.Printf("METHOD: %v\n\n", method)
				}
				/*
					for _, polymer_id := range res.EntitiesInfo.PolymerID {
						fmt.Printf("%v\n", polymer_id)
					}
					for _, non_polymer_id := range res.EntitiesInfo.NonPolymerID {
						fmt.Printf("%v\n", non_polymer_id)
					}
				*/
				fmt.Printf("#######################POLYMERS#######################\n")
				for _, polymer_url := range list_of_urls_polymer {
					// fmt.Printf("%v\n", polymer_url)
					Polymer, err := getdata.GetDataForPolymers(polymer_url)
					if err != nil {
						fmt.Printf("Error getting data")
						return
					}
					fmt.Printf("DESCRIPTION : %v\n", Polymer.EntityGeneralInfo.Description)
					fmt.Printf("TYPE: %v\n", Polymer.EntityPoly.Type)
					fmt.Printf("SEQUENCE: %v\n", Polymer.EntityPoly.Sequence)
					fmt.Printf("LENGTH: %v\n", Polymer.EntityPoly.Length)
					fmt.Printf("FORMULA WEIGHT: %v\n", Polymer.EntityGeneralInfo.FormulaWeight)
					fmt.Printf("SOURCE: %v\n", Polymer.EntityPolySourceHost[0].Source)
					fmt.Printf("HOST: %v\n", Polymer.EntityPolySourceHost[0].Host)
					fmt.Printf("NUMBER OF MOLECULES: %v\n\n", Polymer.EntityGeneralInfo.Number)
				}
				fmt.Printf("#######################NON-POLYMERS#######################\n")
				for _, non_polymer_url := range list_of_urls_non_polymers {
					// fmt.Printf("%v\n", non_polymer_url)
					NonPolymer, err := getdata.GetDataForNonPolymers(non_polymer_url)
					if err != nil {
						fmt.Printf("Error getting data")
						return
					}
					fmt.Printf("NAME: %v\n", NonPolymer.Entity.Name)
					fmt.Printf("COMP ID: %v\n", NonPolymer.Entity.CompID)
					fmt.Printf("DESCRIPTION: %v\n", NonPolymer.Data.Description)
					fmt.Printf("FORMULA WEIGHT: %v\n", NonPolymer.Data.FormulaWeight)
					fmt.Printf("NUMBER OF MOLECULES: %v\n\n", NonPolymer.Data.NumberOfMolecules)
				}
			}
		}
	}
}
