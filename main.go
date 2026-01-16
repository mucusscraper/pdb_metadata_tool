package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	cleaninput "github.com/mucusscraper/pdb_metadata_tool/internal/clean_input"
	"github.com/mucusscraper/pdb_metadata_tool/internal/database"
	getdata "github.com/mucusscraper/pdb_metadata_tool/internal/get_data"
)

const issueURL = "https://data.rcsb.org/rest/v1/core/entry/"

type State struct {
	db *database.Queries
}

func main() {
	dbURL := "postgres://postgres:postgres@localhost:5432/pdbmdt"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return
	}
	dbQueries := database.New(db)
	state := State{
		db: dbQueries,
	}
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
				fmt.Printf("METHOD: %v\n", res.ExptlInfo[0].Method)
				if len(ids) == 1 {
					entry_params := database.CreateEntryParams{
						ID:          uuid.New(),
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						RcsbID:      res.ID,
						DepositDate: res.AccessInfo.DepositDate,
						Doi:         res.ArticleInfo.DOI,
						PaperTitle:  res.ArticleInfo.Title,
						Method:      res.ExptlInfo[0].Method,
					}
				}
				if len(ids) == 2 {
					entry_params := database.CreateEntryParams{
						ID:          uuid.New(),
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						RcsbID:      res.ID,
						DepositDate: res.AccessInfo.DepositDate,
						Doi:         res.ArticleInfo.DOI,
						PaperTitle:  res.ArticleInfo.Title,
						Method:      res.ExptlInfo[0].Method,
						UserGroup:   ids[1],
					}
				} else {
					entry_params := database.CreateEntryParams{
						ID:          uuid.New(),
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						RcsbID:      res.ID,
						DepositDate: res.AccessInfo.DepositDate,
						Doi:         res.ArticleInfo.DOI,
						PaperTitle:  res.ArticleInfo.Title,
						Method:      res.ExptlInfo[0].Method,
					}
				}
				_, err = state.db.CreateEntry(context.Background(), entry_params)
				if err != nil {
					return
				}
				fmt.Printf("#######################POLYMERS#######################\n")
				for _, polymer_url := range list_of_urls_polymer {
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
