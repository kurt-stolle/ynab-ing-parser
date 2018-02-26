package main

import (
	"path/filepath"

	"os"

	"encoding/csv"
	"log"

	"io"

	"strings"

	"github.com/spf13/cobra"
)

// Input and output filenames
var input, output string
var verbose bool

// Define the CLI using spf13's Cobra library
var cli = &cobra.Command{
	Use:   "ynab-ing-parser -i input.csv -o output.csv",
	Short: "Convert ING Bank CSV-exports to a YNAB importable",
	Long: `YNAB ING Parser converts the CSV-exports generated
	by the Dutch ING Bank into CSV based format that can be
	imported by the YNAB budgeting software.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		var err error

		// Check whether the input and output file names make sense
		if input[len(input)-4:] != ".csv" {
			log.Fatal("Input file must have .csv extension")
		} else if output[len(output)-4:] != ".csv" {
			log.Fatal("Output file must have .csv extension")
		}

		// Resolve the absolute file paths
		if input, err = filepath.Abs(input); err != nil {
			log.Fatal(err)
		} else if output, err = filepath.Abs(output); err != nil {
			log.Fatal(err)
		} else if _, err := os.Stat(output); !os.IsNotExist(err) {
			if err == nil {
				log.Fatalf("File at %s already exists!\n", output)
			}

			log.Fatal(err)
		}

		// Open the input file
		fin, err := os.Open(input)
		if err != nil {
			log.Fatal(err)
		}
		defer fin.Close()

		// Open the output file
		fout, err := os.Create(output)
		if err != nil {
			log.Fatal(err)
		}
		defer fout.Close()

		// Logging purposes
		if verbose {
			log.Printf("Parsing %s\n", fin.Name())
		}

		// Parse csv file line-by-line
		r := csv.NewReader(fin)
		r.Read() // Skip the header
		w := csv.NewWriter(fout)
		w.Write([]string{"Date", "Payee", "Memo", "Outflow", "Inflow"})
		for {
			// Read the record
			record, err := r.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			// Create a new transaction
			var transaction = []string{"", "", "", "", ""}

			// Fill out the required fields
			transaction[0] = formatDate(record[0])            // Date
			transaction[1] = formatPayee(record[1])           // Payee
			transaction[2] = formatMemo(record[4], record[8]) // Memo

			// Determine direction of cash flow
			if dir := strings.Trim(strings.ToLower(record[5]), " \n\r\t"); dir == "af" {
				transaction[3] = formatAmount(record[6]) // Outflow

				if verbose {
					log.Printf("Parsed transaction of %s EUR to %s\n", transaction[3], transaction[1])
				}
			} else if dir == "bij" {
				transaction[4] = formatAmount(record[6]) // Inflow

				if verbose {
					log.Printf("Parsed transaction of %s EUR from %s\n", transaction[4], transaction[1])
				}
			} else {
				log.Fatalf("Unknown direction of cash flow: %s\n", dir)
			}

			w.Write(transaction)
		}

		// Verboseness
		if verbose {
			log.Print("Writing out buffer to file\n")
		}

		// Write the buffer
		w.Flush()

		if err := w.Error(); err != nil {
			os.Remove(output)
			log.Fatal(err)
		}
	},
}

func init() {
	// Handle required flags for input an output file
	cli.Flags().StringVarP(&input, "input", "i", "", "Input CSV file (required)")
	cli.MarkFlagRequired("input")
	cli.Flags().StringVarP(&output, "output", "o", "", "Output CSV file (required)")
	cli.MarkFlagRequired("output")
	cli.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

func main() {
	// Execute the CLI
	if err := cli.Execute(); err != nil {
		panic(err)
	}
}
