## License
MIT License © 2026 Daniel Lacerda Oliveira

# PDBMDT - PDB MetaData Tool

PDBMDT/PDBMetaDataTool is an interactive Command-Line-Interface (CLI) tool built in GO that automates metadata annotation of your favorite PDB entries. It also provides a simple and fast way to generate organized HTML report files based on these entries.

## 🤔 Motivation

This package was built due to address two major issues I faced while working on projects with 3D molecular structures: 
1. The difficulty of visualizing structural metadata in HPC systems such as computing clusters, where GUI is often unavailable;  
2. The need for a simple and practical way to store key information from PDB entries in a simple file like a HTML. 
PDBMDT offers an effective solution to both problems, helping structural bioinformatician quickly organize and document essential metadata without relying on graphical tools.

## 🚀 Quick Start

Before using PDBMDT, it's necessary to install two dependencies:

- [GO 1.20+](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)

After installing the software above, navigate to the folder where you would like to install PDBMDT and run:
```bash
# Install PDBMDT
go install github.com/mucusscraper/pdb_metadata_tool@v0.1.0
```
then create a configuration file with the name `.pdbmdtconfig.json` with the following content:
```bash
{
  "db_url": "postgres://{username}:{password}@localhost:5432/pdbmdt?sslmode=disable"
}
# {username} should be the user name in PostgreSQL and {password} should be its corresponding password.
```
Ensure that $HOME/go/bin is in your PATH.

An executable should be created and PDBMDT is ready to use:
```bash
pdb_metadata_tool
```

## 📖 Usage and examples

The help command provides details about all available commands:
```bash
upload {pdb_id} {pdb_id} .... - uploads at least one PDB using its ID to the database
show {pdb_id} {pdb_id} .... - shows info about at least one PDB-ID already uploaded to the database
poly_show {pdb_id} {pdb_id} .... - shows info about polymers of at least one PDB-ID already uploaded to the database
non_poly_show {pdb_id} {pdb_id} .... - shows info about non-polymers of at least one PDB-ID already uploaded to the database
group {group_name} {pdb_id} {pdb_id} .... - inserts at least one PDB-ID already uploaded to the database in the user-specified group
group_show {group_name} - shows info about the entries of an existing group
remove_group {pdb_id} {pdb_id} .... - removes the existing group of at least one PDB-ID already uploaded to the database
report {report_file_name} {pdb_id} {pdb_id} .... - creates a report file with the specified name containing info about the PDB-ID entries already uploaded to the database
group_report {group_name} - creates a report file of the PDB-ID entries of the specified group
exit - exits the program
```
The main workflow to use PDBMDT is to upload related PDB entries, organize them into the same group and generate a report:

```bash
upload 7urv 7urx 6al5 7jic
group cd19 7urv 7urx 6al5 7jic
group_report cd19
```
This will generate an HTML file in the reports folder containing metadata from all entries in the CD19 group.

You can also generate reports without using groups:
```bash
upload 7o4y 7o52
report cool_pdbs 7o4y 7o52
```

## 🤝 Contributing

If you'd like to contribute, feel free to fork the repo and submit pull requests with new features or improvements! 
Also, you are welcome to open issues or contact me so we can discuss about ideas and suggestions.