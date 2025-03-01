package main

import (
	"fmt"
	"log"
	"os"

	"github.com/louisinger/bill/pkg/bill"
	"github.com/louisinger/bill/pkg/ui"
	"github.com/urfave/cli/v2"
)

const version = "1.0.0"

func main() {
	app := &cli.App{
		Name:  "bill",
		Usage: "Generate PDF invoices with Bitcoin payment support",
		Action: func(c *cli.Context) error {
			// Default action: launch GUI
			billApp := ui.NewBillApp()
			billApp.Run()
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "Generate a new invoice",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "invoice.pdf",
						Usage:   "Output PDF file path",
					},
					&cli.StringFlag{
						Name:    "template",
						Aliases: []string{"t"},
						Usage:   "Path to template JSON file",
					},
				},
				Action: func(c *cli.Context) error {
					var template *bill.BillTemplate
					if templatePath := c.String("template"); templatePath != "" {
						var err error
						template, err = bill.LoadTemplate(templatePath)
						if err != nil {
							return cli.Exit(fmt.Sprintf("Error loading template: %v", err), 1)
						}
						fmt.Println("Template loaded successfully")
					}

					billData := bill.CollectBillData(template)
					outputPath := c.String("output")

					fmt.Printf("Generating bill PDF to %s...\n", outputPath)
					err := bill.GeneratePDF(billData, outputPath)
					if err != nil {
						return cli.Exit(fmt.Sprintf("Error generating PDF: %v", err), 1)
					}

					fmt.Println("Bill PDF generated successfully!")
					fmt.Printf("Total amount: %.2f %s\n", billData.Total, billData.Currency)
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the version",
				Action: func(c *cli.Context) error {
					fmt.Printf("bill version %s\n", version)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
