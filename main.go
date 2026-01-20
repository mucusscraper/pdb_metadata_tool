package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	cleaninput "github.com/mucusscraper/pdb_metadata_tool/internal/clean_input"
	"github.com/mucusscraper/pdb_metadata_tool/internal/database"
	getdata "github.com/mucusscraper/pdb_metadata_tool/internal/get_data"
	reportgenerator "github.com/mucusscraper/pdb_metadata_tool/internal/report_generator"
)

const issueURL = "https://data.rcsb.org/rest/v1/core/entry/"

type State struct {
	db *database.Queries
}

type Config struct {
	DbURL string `json:"db_url"`
}

func main() {
	file, _ := os.Open(".pdbmdt.config.json")
	defer file.Close()
	var config Config
	json.NewDecoder(file).Decode(&config)
	dbURL := config.DbURL
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
			if ids[0] == "group" && len(ids) >= 3 {
				for _, id := range ids[2:] {
					group_params := database.InsertGroupParams{
						UserGroup: ids[1],
						RcsbID:    strings.ToUpper(id),
					}
					_, err = state.db.InsertGroup(context.Background(), group_params)
					if err != nil {
						fmt.Printf("Error inserting group\n")
						fmt.Printf("%v", err)
						continue
					}
					fmt.Printf("Added %v to group %v\n", id, ids[1])
				}
			} else if ids[0] == "show" && len(ids) >= 2 {
				for _, id := range ids[1:] {
					entry, err := dbQueries.GetEntry(context.Background(), strings.ToUpper(id))
					if err != nil {
						fmt.Printf("Error getting entry\n")
						fmt.Printf("%v", err)
						continue
					}
					fmt.Printf("\n")
					fmt.Printf("#######################ENTRY-ALREADY-LOADED#######################\n")
					fmt.Printf("RCSB ID: %v\n", entry.RcsbID)
					fmt.Printf("DEPOSIT DATE: %v\n", entry.DepositDate)
					fmt.Printf("DOI: %v\n", entry.Doi)
					fmt.Printf("PAPER TITLE: %v\n", entry.PaperTitle)
					fmt.Printf("METHOD: %v\n\n", entry.Method)
				}
			} else if ids[0] == "group_show" && len(ids) >= 2 {
				for _, id := range ids[1:] {
					entries, err := dbQueries.GetEntryByUserGroup(context.Background(), id)
					if err != nil {
						fmt.Printf("Error getting entry by user group\n")
						fmt.Printf("%v\n", err)
						continue
					}
					for _, entry := range entries {
						fmt.Printf("\n")
						fmt.Printf("#######################%v#######################\n", id)
						fmt.Printf("RCSB ID: %v\n", entry.RcsbID)
						fmt.Printf("DEPOSIT DATE: %v\n", entry.DepositDate)
						fmt.Printf("DOI: %v\n", entry.Doi)
						fmt.Printf("PAPER TITLE: %v\n", entry.PaperTitle)
						fmt.Printf("METHOD: %v\n\n", entry.Method)
					}
				}
			} else if ids[0] == "poly_show" && len(ids) >= 2 {
				for _, id := range ids[1:] {
					polys, err := dbQueries.GetPolys(context.Background(), strings.ToUpper(id))
					if err != nil {
						fmt.Printf("Error getting poly by rcsb_id\n")
						fmt.Printf("%v\n", err)
						continue
					}
					fmt.Printf("%v\n", id)
					fmt.Printf("%v\n", len(polys))
					for _, poly := range polys {
						fmt.Printf("\n")
						fmt.Printf("#######################%v-POLYS#######################\n", strings.ToUpper(id))
						fmt.Printf("POLYMER DESCRIPTION: %v\n", poly.Poldescription)
						fmt.Printf("POLYMER TYPE: %v\n", poly.Poltype)
						fmt.Printf("POLYMER SEQUENCE: %v\n", poly.Polsequence)
						fmt.Printf("POLYMER LENGTH: %v\n", poly.Pollength)
						fmt.Printf("FORMULA WEIGHT: %v\n", poly.Formulaweight)
						fmt.Printf("SOURCE: %v\n", poly.Source)
						fmt.Printf("HOST: %v\n", poly.Host)
						fmt.Printf("NUMBER OF MOLECULES: %v\n", poly.NumberOfMolecules)
					}
				}
			} else if ids[0] == "non_poly_show" && len(ids) >= 2 {
				for _, id := range ids[1:] {
					non_polys, err := dbQueries.GetNonPolys(context.Background(), strings.ToUpper(id))
					if err != nil {
						fmt.Printf("Error getting non_poly by rcsb_id\n")
						fmt.Printf("%v\n", err)
						continue
					}
					for _, non_poly := range non_polys {
						fmt.Printf("\n")
						fmt.Printf("#######################%v-NON-POLYS#######################\n", strings.ToUpper(id))
						fmt.Printf("NAME: %v\n", non_poly.Nonpolname)
						fmt.Printf("COMPD ID: %v\n", non_poly.CompID)
						fmt.Printf("DESCRIPTION: %v\n", non_poly.Nonpoldescription)
						fmt.Printf("FORMULA WEIGHT: %v\n", non_poly.FormulaWeight)
						fmt.Printf("NUMBER OF MOLECULES: %v\n", non_poly.NumberOfMolecules)
					}
				}
			} else if ids[0] == "remove_group" && len(ids) >= 2 {
				for _, id := range ids[1:] {
					_, err := dbQueries.RemoveGroup(context.Background(), strings.ToUpper(id))
					if err != nil {
						fmt.Printf("Error removing user group\n")
						fmt.Printf("%v\n", err)
						continue
					}
					fmt.Printf("Removed group from %v\n", id)
				}
			} else if ids[0] == "group_report" && len(ids) == 2 {
				entries, err := dbQueries.GetEntryByUserGroup(context.Background(), ids[1])
				if err != nil {
					fmt.Printf("Error getting entry by user group\n")
					fmt.Printf("%v\n", err)
					continue
				}
				var preReports []getdata.PreReport
				for _, entry := range entries {
					ACSI := getdata.ArticleAccessInfo{
						DOI:   entry.Doi,
						Title: entry.PaperTitle,
					}
					RCSBAI := getdata.RcsbAccessionInfo{
						DepositDate: entry.DepositDate,
					}
					var EXPTLAI []getdata.ExptlAccessInfo
					EXPTLAI = append(EXPTLAI, getdata.ExptlAccessInfo{
						Method: entry.Method,
					})
					EI := getdata.EntitiesAccessInfo{}
					dataEntry := getdata.ProteinIssue{
						ID:           entry.RcsbID,
						ArticleInfo:  ACSI,
						AccessInfo:   RCSBAI,
						ExptlInfo:    EXPTLAI,
						EntitiesInfo: EI,
					}
					polys, err := dbQueries.GetPolys(context.Background(), entry.RcsbID)
					if err != nil {
						fmt.Printf("Error getting polys\n")
						fmt.Printf("%v\n", err)
						continue
					}
					var poylmers_data []getdata.PolymerIssue
					for _, poly := range polys {
						EGIA := getdata.EntityGeneralInfoAccess{
							FormulaWeight: poly.Formulaweight,
							Description:   poly.Poldescription,
							Number:        int(poly.NumberOfMolecules),
						}
						EPA := getdata.EntityPolyAccess{
							Type:     poly.Poltype,
							Length:   int(poly.Pollength),
							Sequence: poly.Polsequence,
						}
						var EPSHA []getdata.EntityPolySourceHostAccess
						NEPSHA := getdata.EntityPolySourceHostAccess{
							Source: poly.Source,
							Host:   poly.Host,
						}
						EPSHA = append(EPSHA, NEPSHA)
						NewPolymerIssue := getdata.PolymerIssue{
							EntityGeneralInfo:    EGIA,
							EntityPoly:           EPA,
							EntityPolySourceHost: EPSHA,
						}
						poylmers_data = append(poylmers_data, NewPolymerIssue)
					}
					non_polys, err := dbQueries.GetNonPolys(context.Background(), entry.RcsbID)
					if err != nil {
						fmt.Printf("Error getting non-polys\n")
						fmt.Printf("%v\n", err)
						continue
					}
					var non_poly_data []getdata.NonPolymerIssue
					for _, non_poly := range non_polys {
						new_entity := getdata.NameEntityNonPolymerAccession{
							Name:   non_poly.Nonpolname,
							CompID: non_poly.CompID,
						}
						new_data := getdata.DataEntityNonPolymerAccession{
							FormulaWeight:     non_poly.FormulaWeight,
							Description:       non_poly.Nonpoldescription,
							NumberOfMolecules: int(non_poly.NumberOfMolecules),
						}
						new_non_poly_data := getdata.NonPolymerIssue{
							Entity: new_entity,
							Data:   new_data,
						}
						non_poly_data = append(non_poly_data, new_non_poly_data)
					}
					preReport := getdata.PreReport{
						Entry:       dataEntry,
						Polymers:    poylmers_data,
						NonPolymers: non_poly_data,
					}
					preReports = append(preReports, preReport)
				}
				report := reportgenerator.Report{
					PreReport: preReports,
					Grouped:   true,
				}
				filename := ids[1]
				err = reportgenerator.GenerateHTML(filename, report)
				if err != nil {
					fmt.Printf("Error generating HTML: %v\n", err)
					continue
				}
				fmt.Printf("%v HTML report generated!\n", filename)
			} else if ids[0] == "report" && len(ids) >= 3 {
				var entries []database.Entry
				for _, id := range ids[2:] {
					entry, err := state.db.GetEntry(context.Background(), strings.ToUpper(id))
					if err != nil {
						fmt.Printf("Error getting entry\n")
						fmt.Printf("%v\n", err)
						continue
					}
					entries = append(entries, entry)
				}
				var preReports []getdata.PreReport
				for _, entry := range entries {
					ACSI := getdata.ArticleAccessInfo{
						DOI:   entry.Doi,
						Title: entry.PaperTitle,
					}
					RCSBAI := getdata.RcsbAccessionInfo{
						DepositDate: entry.DepositDate,
					}
					var EXPTLAI []getdata.ExptlAccessInfo
					EXPTLAI = append(EXPTLAI, getdata.ExptlAccessInfo{
						Method: entry.Method,
					})
					EI := getdata.EntitiesAccessInfo{}
					dataEntry := getdata.ProteinIssue{
						ID:           entry.RcsbID,
						ArticleInfo:  ACSI,
						AccessInfo:   RCSBAI,
						ExptlInfo:    EXPTLAI,
						EntitiesInfo: EI,
					}
					polys, err := dbQueries.GetPolys(context.Background(), entry.RcsbID)
					if err != nil {
						fmt.Printf("Error getting polys\n")
						fmt.Printf("%v\n", err)
						continue
					}
					var poylmers_data []getdata.PolymerIssue
					for _, poly := range polys {
						EGIA := getdata.EntityGeneralInfoAccess{
							FormulaWeight: poly.Formulaweight,
							Description:   poly.Poldescription,
							Number:        int(poly.NumberOfMolecules),
						}
						EPA := getdata.EntityPolyAccess{
							Type:     poly.Poltype,
							Length:   int(poly.Pollength),
							Sequence: poly.Polsequence,
						}
						var EPSHA []getdata.EntityPolySourceHostAccess
						NEPSHA := getdata.EntityPolySourceHostAccess{
							Source: poly.Source,
							Host:   poly.Host,
						}
						EPSHA = append(EPSHA, NEPSHA)
						NewPolymerIssue := getdata.PolymerIssue{
							EntityGeneralInfo:    EGIA,
							EntityPoly:           EPA,
							EntityPolySourceHost: EPSHA,
						}
						poylmers_data = append(poylmers_data, NewPolymerIssue)
					}
					non_polys, err := dbQueries.GetNonPolys(context.Background(), entry.RcsbID)
					if err != nil {
						fmt.Printf("Error getting non-polys\n")
						fmt.Printf("%v\n", err)
						continue
					}
					var non_poly_data []getdata.NonPolymerIssue
					for _, non_poly := range non_polys {
						new_entity := getdata.NameEntityNonPolymerAccession{
							Name:   non_poly.Nonpolname,
							CompID: non_poly.CompID,
						}
						new_data := getdata.DataEntityNonPolymerAccession{
							FormulaWeight:     non_poly.FormulaWeight,
							Description:       non_poly.Nonpoldescription,
							NumberOfMolecules: int(non_poly.NumberOfMolecules),
						}
						new_non_poly_data := getdata.NonPolymerIssue{
							Entity: new_entity,
							Data:   new_data,
						}
						non_poly_data = append(non_poly_data, new_non_poly_data)
					}
					preReport := getdata.PreReport{
						Entry:       dataEntry,
						Polymers:    poylmers_data,
						NonPolymers: non_poly_data,
					}
					preReports = append(preReports, preReport)
				}
				report := reportgenerator.Report{
					PreReport: preReports,
					Grouped:   false,
				}
				filename := ids[1]
				err = reportgenerator.GenerateHTML(filename, report)
				if err != nil {
					fmt.Printf("Error generating HTML: %v\n", err)
					continue
				}
				fmt.Printf("%v HTML report generated!\n", filename)
			} else if ids[0] == "upload" && len(ids) >= 2 {
				for _, id := range ids[1:] {
					entry, err := state.db.GetEntry(context.Background(), strings.ToUpper(id))
					// PDB was not uploaded before
					if err == nil {
						fmt.Printf("\n")
						fmt.Printf("#######################ENTRY-ALREADY-LOADED#######################\n")
						fmt.Printf("RCSB ID: %v\n", entry.RcsbID)
						fmt.Printf("DEPOSIT DATE: %v\n", entry.DepositDate)
						fmt.Printf("DOI: %v\n", entry.Doi)
						fmt.Printf("PAPER TITLE: %v\n", entry.PaperTitle)
						fmt.Printf("METHOD: %v\n", entry.Method)
						fmt.Printf("GROUP: %v\n\n", entry.UserGroup)
					} else {
						new_id := uuid.New()
						url := fmt.Sprintf("%v%v", issueURL, strings.ToUpper(id))
						res, err, list_of_urls_polymer, list_of_urls_non_polymers := getdata.GetIssueDataEntry(url, strings.ToUpper(id))
						if err != nil {
							fmt.Printf("Couldn't get data from the PDB ID provided.\n")
						}
						fmt.Printf("\n")
						fmt.Printf("#######################GENERAL-DATA#######################\n")
						fmt.Printf("RCSB ID: %v\n", res.ID)
						fmt.Printf("DEPOSIT DATE: %v\n", res.AccessInfo.DepositDate)
						fmt.Printf("DOI: %v\n", res.ArticleInfo.DOI)
						fmt.Printf("PAPER TITLE: %v\n", res.ArticleInfo.Title)
						fmt.Printf("METHOD: %v\n\n", res.ExptlInfo[0].Method)
						var entry_params database.CreateEntryParams
						if len(ids) == 2 {
							entry_params = database.CreateEntryParams{
								ID:          new_id,
								CreatedAt:   time.Now(),
								UpdatedAt:   time.Now(),
								RcsbID:      strings.ToUpper(res.ID),
								DepositDate: res.AccessInfo.DepositDate,
								Doi:         res.ArticleInfo.DOI,
								PaperTitle:  res.ArticleInfo.Title,
								Method:      res.ExptlInfo[0].Method,
							}
						} else {
							entry_params = database.CreateEntryParams{
								ID:          new_id,
								CreatedAt:   time.Now(),
								UpdatedAt:   time.Now(),
								RcsbID:      strings.ToUpper(res.ID),
								DepositDate: res.AccessInfo.DepositDate,
								Doi:         res.ArticleInfo.DOI,
								PaperTitle:  res.ArticleInfo.Title,
								Method:      res.ExptlInfo[0].Method,
							}
						}
						/*
							if len(ids) == 2 {
								entry_params = database.CreateEntryParams{
									ID:          new_id,
									CreatedAt:   time.Now(),
									UpdatedAt:   time.Now(),
									RcsbID:      strings.ToUpper(res.ID),
									DepositDate: res.AccessInfo.DepositDate,
									Doi:         res.ArticleInfo.DOI,
									PaperTitle:  res.ArticleInfo.Title,
									Method:      res.ExptlInfo[0].Method,
									UserGroup:   ids[1],
								}

							}*/
						_, err = state.db.CreateEntry(context.Background(), entry_params)
						if err != nil {
							fmt.Printf("Error creating entry\n")
							fmt.Printf("%v\n", err)
							continue
						}
						fmt.Printf("\n")
						fmt.Printf("#######################POLYMERS#######################\n")
						for _, polymer_url := range list_of_urls_polymer {
							Polymer, err := getdata.GetDataForPolymers(polymer_url)
							if err != nil {
								fmt.Printf("Error getting data\n")
								fmt.Printf("%v\n", err)
								continue
							}
							fmt.Printf("DESCRIPTION : %v\n", Polymer.EntityGeneralInfo.Description)
							fmt.Printf("TYPE: %v\n", Polymer.EntityPoly.Type)
							fmt.Printf("SEQUENCE: %v\n", Polymer.EntityPoly.Sequence)
							fmt.Printf("LENGTH: %v\n", Polymer.EntityPoly.Length)
							fmt.Printf("FORMULA WEIGHT: %v\n", Polymer.EntityGeneralInfo.FormulaWeight)
							fmt.Printf("SOURCE: %v\n", Polymer.EntityPolySourceHost[0].Source)
							fmt.Printf("HOST: %v\n", Polymer.EntityPolySourceHost[0].Host)
							fmt.Printf("NUMBER OF MOLECULES: %v\n\n", Polymer.EntityGeneralInfo.Number)
							poly_params := database.CreatePolyParams{
								ID:                uuid.New(),
								EntryID:           new_id,
								Poldescription:    Polymer.EntityGeneralInfo.Description,
								Poltype:           Polymer.EntityPoly.Type,
								Polsequence:       Polymer.EntityPoly.Sequence,
								Pollength:         int32(Polymer.EntityPoly.Length),
								Formulaweight:     Polymer.EntityGeneralInfo.FormulaWeight,
								Source:            Polymer.EntityPolySourceHost[0].Source,
								Host:              Polymer.EntityPolySourceHost[0].Host,
								NumberOfMolecules: int32(Polymer.EntityGeneralInfo.Number),
								CreatedAt:         time.Now(),
							}
							_, err = state.db.CreatePoly(context.Background(), poly_params)
							if err != nil {
								fmt.Printf("Error creating polymer row in database\n")
								fmt.Printf("%v\n", err)
								continue
							}
						}
						fmt.Printf("\n")
						fmt.Printf("#######################NON-POLYMERS#######################\n")
						for _, non_polymer_url := range list_of_urls_non_polymers {
							NonPolymer, err := getdata.GetDataForNonPolymers(non_polymer_url)
							if err != nil {
								fmt.Printf("Error getting data\n")
								fmt.Printf("%v\n", err)
								continue
							}
							fmt.Printf("NAME: %v\n", NonPolymer.Entity.Name)
							fmt.Printf("COMP ID: %v\n", NonPolymer.Entity.CompID)
							fmt.Printf("DESCRIPTION: %v\n", NonPolymer.Data.Description)
							fmt.Printf("FORMULA WEIGHT: %v\n", NonPolymer.Data.FormulaWeight)
							fmt.Printf("NUMBER OF MOLECULES: %v\n\n", NonPolymer.Data.NumberOfMolecules)
							non_poly_params := database.CreateNonPolyParams{
								ID:                uuid.New(),
								EntryID:           new_id,
								Nonpolname:        NonPolymer.Entity.Name,
								CompID:            NonPolymer.Entity.CompID,
								Nonpoldescription: NonPolymer.Data.Description,
								FormulaWeight:     NonPolymer.Data.FormulaWeight,
								NumberOfMolecules: int32(NonPolymer.Data.NumberOfMolecules),
								CreatedAt:         time.Now(),
							}
							_, err = state.db.CreateNonPoly(context.Background(), non_poly_params)
							if err != nil {
								fmt.Printf("Error creating non-polymer row in database\n")
								fmt.Printf("%v\n", err)
								continue
							}
						}
					}
				}
			} else if ids[0] == "help" && len(ids) == 1 {
				fmt.Printf("#######################HELP#######################\n\n")
				fmt.Printf("upload {pdb_id} {pdb_id} .... - uploads at least one PDB using its ID to the database\n\n")
				fmt.Printf("show {pdb_id} {pdb_id} .... - shows info about at least one PDB-ID already uploaded to the database\n\n")
				fmt.Printf("poly_show {pdb_id} {pdb_id} .... - shows info about polymers of at least one PDB-ID already uploaded to the database\n\n")
				fmt.Printf("non_poly_show {pdb_id} {pdb_id} .... - shows info about non-polymers of at least one PDB-ID already uploaded to the database\n\n")
				fmt.Printf("group {group_name} {pdb_id} {pdb_id} .... - inserts at least one PDB-ID already uploaded to the database in the user-specified group\n\n")
				fmt.Printf("group_show {group_name} - shows info about the entries of an existing group\n\n")
				fmt.Printf("remove_group {pdb_id} {pdb_id} .... - removes the existing group of at least one PDB-ID already uploaded to the database\n\n")
				fmt.Printf("report {report_file_name} {pdb_id} {pdb_id} .... - creates a report file with the specified name containing info about the PDB-ID entries already uploaded to the database\n\n")
				fmt.Printf("group_report {group_name} - creates a report file of the PDB-ID entries of the specified group\n\n")
			} else {
				fmt.Printf("Command not found!\n")
			}
		}
	}
}
