package getdata

type PreReport struct {
	Entry       ProteinIssue
	Polymers    []PolymerIssue
	NonPolymers []NonPolymerIssue
}

type ProteinIssue struct {
	ID           string             `json:"rcsb_id"`
	AccessInfo   RcsbAccessionInfo  `json:"rcsb_accession_info"`
	ArticleInfo  ArticleAccessInfo  `json:"rcsb_primary_citation"`
	ExptlInfo    []ExptlAccessInfo  `json:"exptl"`
	EntitiesInfo EntitiesAccessInfo `json:"rcsb_entry_container_identifiers"`
}

type RcsbAccessionInfo struct {
	DepositDate string `json:"deposit_date"`
}

type ArticleAccessInfo struct {
	DOI   string `json:"pdbx_database_id_doi"`
	Title string `json:"title"`
}

type ExptlAccessInfo struct {
	Method string `json:"method"`
}

type EntitiesAccessInfo struct {
	PolymerID    []string `json:"polymer_entity_ids"`
	NonPolymerID []string `json:"non_polymer_entity_ids"`
}

type NonPolymerIssue struct {
	Entity NameEntityNonPolymerAccession `json:"pdbx_entity_nonpoly"`
	Data   DataEntityNonPolymerAccession `json:"rcsb_nonpolymer_entity"`
}

type NameEntityNonPolymerAccession struct {
	Name   string `json:"name"`
	CompID string `json:"comp_id"`
}

type DataEntityNonPolymerAccession struct {
	FormulaWeight     float32 `json:"formula_weight"`
	Description       string  `json:"pdbx_description"`
	NumberOfMolecules int     `json:"pdbx_number_of_molecules"`
}

type PolymerIssue struct {
	EntityGeneralInfo    EntityGeneralInfoAccess      `json:"rcsb_polymer_entity"`
	EntityPoly           EntityPolyAccess             `json:"entity_poly"`
	EntityPolySourceHost []EntityPolySourceHostAccess `json:"entity_src_gen"`
}

type EntityPolyAccess struct {
	Type     string `json:"type"`
	Length   int    `json:"rcsb_sample_sequence_length"`
	Sequence string `json:"pdbx_seq_one_letter_code"`
}

type EntityPolySourceHostAccess struct {
	Source string `json:"pdbx_gene_src_scientific_name"`
	Host   string `json:"pdbx_host_org_scientific_name"`
}

type EntityGeneralInfoAccess struct {
	FormulaWeight float32 `json:"formula_weight"`
	Description   string  `json:"pdbx_description"`
	Number        int     `json:"pdbx_number_of_molecules"`
}
