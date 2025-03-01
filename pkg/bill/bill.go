package bill

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

type Bill struct {
	Number         string
	Date           time.Time
	CompanyName    string
	Address        string
	VATNumber      string
	ToCompanyName  string
	ToAddress      string
	ToVATNumber    string
	Items          []BillItem
	Total          float64
	Currency       string
	BitcoinAddress string
}

type BillItem struct {
	Description string
	Quantity    int
	UnitPrice   float64
	Total       float64
}

type BillTemplate struct {
	CompanyName    string         `json:"company_name"`
	Address        string         `json:"address"`
	VATNumber      string         `json:"vat_number"`
	ToCompanyName  string         `json:"to_company_name"`
	ToAddress      string         `json:"to_address"`
	ToVATNumber    string         `json:"to_vat_number"`
	BitcoinAddress string         `json:"bitcoin_address"`
	Currency       string         `json:"currency"`
	Items          []TemplateItem `json:"items"`
}

type TemplateItem struct {
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

func generateQRCode(address string) string {
	qrFile := "bitcoin_qr.png"
	err := qrcode.WriteFile(fmt.Sprintf("bitcoin:%s", address), qrcode.Medium, 256, qrFile)
	if err != nil {
		fmt.Printf("Error generating QR code: %v\n", err)
		return ""
	}
	return qrFile
}

func GeneratePDF(bill Bill, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Enable UTF-8 encoding
	pdf.SetFont("Helvetica", "", 10)
	tr := pdf.UnicodeTranslatorFromDescriptor("") // Create UTF-8 translator

	pdf.AddPage()

	// Add colors
	pdf.SetDrawColor(28, 72, 107)   // Dark blue for lines
	pdf.SetFillColor(240, 248, 255) // Light blue for header backgrounds
	pdf.SetTextColor(28, 72, 107)   // Dark blue for text

	// Header section
	pdf.SetFont("Helvetica", "B", 24)
	pdf.CellFormat(190, 10, "INVOICE", "", 0, "", false, 0, "")
	pdf.Ln(12)

	// Invoice number and Date section - Moved above separator
	pdf.SetTextColor(28, 72, 107)
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(15, 8, "No.", "", 0, "", false, 0, "")
	pdf.SetFont("Helvetica", "I", 10)
	pdf.CellFormat(40, 8, bill.Number, "", 0, "", false, 0, "")
	pdf.SetX(110)
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(15, 8, "Date", "", 0, "", false, 0, "")
	pdf.SetFont("Helvetica", "I", 10)
	pdf.CellFormat(90, 8, "  "+bill.Date.Format("January 2, 2006"), "", 0, "", false, 0, "")
	pdf.Ln(12)

	// Add line under everything
	pdf.SetLineWidth(0.5)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(15)

	// Company details in two columns
	pdf.SetFont("Helvetica", "B", 12)
	leftCol := 10.0
	rightCol := 110.0
	startY := pdf.GetY() // Store the starting Y position

	// Left column - From (Your company)
	pdf.SetFillColor(28, 72, 107)
	pdf.Rect(leftCol, startY, 90, 8, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(leftCol+5, startY+2)
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(80, 4, "FROM", "", 0, "", false, 0, "")

	// Company details section - Increased gap after header
	pdf.SetFillColor(240, 248, 255)
	pdf.Rect(leftCol, startY+10, 90, 40, "F")
	pdf.SetTextColor(28, 72, 107)

	// Company Name - Adjusted Y positions
	pdf.SetXY(leftCol+5, startY+12)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(80, 6, tr(bill.CompanyName), "", 0, "", false, 0, "")

	// Address - Adjusted Y positions
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetXY(leftCol+5, startY+19)
	pdf.MultiCell(80, 5, tr(bill.Address), "", "", false)

	// VAT Number - Adjusted Y positions
	pdf.SetXY(leftCol+5, startY+39)
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(80, 6, tr(bill.VATNumber), "", 0, "", false, 0, "")

	// Right column - TO section
	pdf.SetFillColor(28, 72, 107)
	pdf.Rect(rightCol, startY, 90, 8, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(rightCol+5, startY+2)
	pdf.SetFont("Helvetica", "B", 11)
	pdf.CellFormat(80, 4, "TO", "", 0, "", false, 0, "")

	// TO Details section
	pdf.SetFillColor(240, 248, 255)
	pdf.Rect(rightCol, startY+10, 90, 40, "F")
	pdf.SetTextColor(28, 72, 107)

	// Recipient Company Name
	pdf.SetXY(rightCol+5, startY+12)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(80, 6, tr(bill.ToCompanyName), "", 0, "", false, 0, "")

	// Recipient Address
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetXY(rightCol+5, startY+19)
	pdf.MultiCell(80, 5, tr(bill.ToAddress), "", "", false)

	// Recipient VAT Number
	pdf.SetXY(rightCol+5, startY+39)
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(80, 6, tr(bill.ToVATNumber), "", 0, "", false, 0, "")

	// Move to next section
	pdf.SetXY(10, startY+59)

	// Items table header
	pdf.SetFillColor(28, 72, 107)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 11)

	pdf.Rect(10, pdf.GetY(), 190, 10, "F")
	pdf.CellFormat(90, 10, "  Description", "", 0, "", false, 0, "")
	pdf.CellFormat(30, 10, "Quantity", "", 0, "", false, 0, "")
	pdf.CellFormat(35, 10, "Unit Price", "", 0, "", false, 0, "")
	pdf.CellFormat(35, 10, "Total", "", 0, "", false, 0, "")
	pdf.Ln(10)

	// Table contents
	pdf.SetTextColor(28, 72, 107)
	pdf.SetFont("Helvetica", "", 11)

	alternate := false
	for _, item := range bill.Items {
		if alternate {
			pdf.SetFillColor(240, 248, 255)
			pdf.Rect(10, pdf.GetY(), 190, 10, "F")
		}
		pdf.SetX(10)
		pdf.CellFormat(90, 10, "  "+tr(item.Description), "", 0, "", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", item.Quantity), "", 0, "", false, 0, "")
		pdf.CellFormat(35, 10, tr(fmt.Sprintf("%.2f %s", item.UnitPrice, bill.Currency)), "", 0, "", false, 0, "")
		pdf.CellFormat(35, 10, tr(fmt.Sprintf("%.2f %s", item.Total, bill.Currency)), "", 0, "", false, 0, "")
		pdf.Ln(10)
		alternate = !alternate
	}

	// Total section
	pdf.Ln(15)
	pdf.SetLineWidth(0.2)
	pdf.Line(120, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(8)

	pdf.SetTextColor(28, 72, 107)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.SetX(120)
	pdf.CellFormat(50, 10, "TOTAL", "", 0, "", false, 0, "")
	pdf.SetX(170)
	pdf.CellFormat(30, 10, tr(fmt.Sprintf("%.2f %s", bill.Total, bill.Currency)), "", 0, "", false, 0, "")

	// Bitcoin Payment Section
	pdf.Ln(25)
	pdf.SetFillColor(240, 248, 255)
	pdf.SetTextColor(28, 72, 107)
	pdf.Rect(10, pdf.GetY(), 190, 50, "F")

	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetX(15)
	pdf.CellFormat(180, 10, "Bitcoin Payment Details", "", 0, "", false, 0, "")
	pdf.Ln(10)

	qrFile := generateQRCode(bill.BitcoinAddress)
	if qrFile != "" {
		pdf.Image(qrFile, 15, pdf.GetY(), 30, 30, false, "", 0, "")
		defer os.Remove(qrFile)
	}

	pdf.SetFillColor(255, 255, 255)
	pdf.Rect(50, pdf.GetY(), 140, 8, "F")
	pdf.SetX(50)
	pdf.SetFont("Courier", "", 8)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(140, 8, "  "+bill.BitcoinAddress, "", 0, "", false, 0, "")

	pdf.Ln(10)
	pdf.SetX(50)
	pdf.SetFont("Helvetica", "I", 8)
	pdf.SetTextColor(28, 72, 107)
	pdf.CellFormat(140, 6, "Please scan the QR code or copy the address above to make your payment", "", 0, "", false, 0, "")

	pdf.SetY(275)
	pdf.SetTextColor(128, 128, 128)

	return pdf.OutputFileAndClose(outputPath)
}

func readString(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func readFloat(reader *bufio.Reader, prompt string) float64 {
	for {
		input := readString(reader, prompt)
		value, err := strconv.ParseFloat(input, 64)
		if err == nil {
			return value
		}
		fmt.Println("Please enter a valid number")
	}
}

func readInt(reader *bufio.Reader, prompt string) int {
	for {
		input := readString(reader, prompt)
		value, err := strconv.Atoi(input)
		if err == nil {
			return value
		}
		fmt.Println("Please enter a valid number")
	}
}

func LoadTemplate(path string) (*BillTemplate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var template BillTemplate
	if err := json.Unmarshal(data, &template); err != nil {
		return nil, err
	}

	return &template, nil
}

func CollectBillData(template *BillTemplate) Bill {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== Bill Generator ===")

	bill := Bill{
		Date:     time.Now(),
		Currency: "€",
	}

	if template != nil && template.Currency != "" {
		bill.Currency = template.Currency
	}

	bill.Number = readString(reader, "Bill Number (e.g., INV-2024-001): ")

	fmt.Println("\n--- Your Company Details ---")
	defaultCompany := ""
	if template != nil {
		defaultCompany = template.CompanyName
	}
	if defaultCompany != "" {
		fmt.Printf("Company Name [%s]: ", defaultCompany)
		input := readString(reader, "")
		if input == "" {
			bill.CompanyName = defaultCompany
		} else {
			bill.CompanyName = input
		}
	} else {
		bill.CompanyName = readString(reader, "Company Name: ")
	}

	fmt.Println("Address (press Enter twice when done):")
	if template != nil && template.Address != "" {
		fmt.Printf("Default address:\n%s\nPress Enter to keep, or type new address:\n", template.Address)
		line := readString(reader, "")
		if line == "" {
			bill.Address = template.Address
		} else {
			var addressLines []string
			addressLines = append(addressLines, line)
			for {
				line = readString(reader, "")
				if line == "" {
					break
				}
				addressLines = append(addressLines, line)
			}
			bill.Address = strings.Join(addressLines, "\n")
		}
	} else {
		var addressLines []string
		for {
			line := readString(reader, "")
			if line == "" {
				break
			}
			addressLines = append(addressLines, line)
		}
		bill.Address = strings.Join(addressLines, "\n")
	}

	if template != nil && template.VATNumber != "" {
		fmt.Printf("VAT Number [%s]: ", template.VATNumber)
		input := readString(reader, "")
		if input == "" {
			bill.VATNumber = template.VATNumber
		} else {
			bill.VATNumber = input
		}
	} else {
		bill.VATNumber = readString(reader, "VAT Number: ")
	}

	fmt.Println("\n--- Client Details ---")
	if template != nil && template.ToCompanyName != "" {
		fmt.Printf("Client Company Name [%s]: ", template.ToCompanyName)
		input := readString(reader, "")
		if input == "" {
			bill.ToCompanyName = template.ToCompanyName
		} else {
			bill.ToCompanyName = input
		}
	} else {
		bill.ToCompanyName = readString(reader, "Client Company Name: ")
	}

	fmt.Println("Client Address (press Enter twice when done):")
	if template != nil && template.ToAddress != "" {
		fmt.Printf("Default client address:\n%s\nPress Enter to keep, or type new address:\n", template.ToAddress)
		line := readString(reader, "")
		if line == "" {
			bill.ToAddress = template.ToAddress
		} else {
			var addressLines []string
			addressLines = append(addressLines, line)
			for {
				line = readString(reader, "")
				if line == "" {
					break
				}
				addressLines = append(addressLines, line)
			}
			bill.ToAddress = strings.Join(addressLines, "\n")
		}
	} else {
		var addressLines []string
		for {
			line := readString(reader, "")
			if line == "" {
				break
			}
			addressLines = append(addressLines, line)
		}
		bill.ToAddress = strings.Join(addressLines, "\n")
	}

	if template != nil && template.ToVATNumber != "" {
		fmt.Printf("Client VAT Number [%s]: ", template.ToVATNumber)
		input := readString(reader, "")
		if input == "" {
			bill.ToVATNumber = template.ToVATNumber
		} else {
			bill.ToVATNumber = input
		}
	} else {
		bill.ToVATNumber = readString(reader, "Client VAT Number: ")
	}

	fmt.Println("\n--- Bill Items ---")
	var items []BillItem
	total := 0.0

	if template != nil && len(template.Items) > 0 {
		for _, templateItem := range template.Items {
			fmt.Printf("\nTemplate item:\nDescription: %s\nQuantity: %d\nUnit Price: %.2f €\n",
				templateItem.Description, templateItem.Quantity, templateItem.UnitPrice)
			fmt.Print("Use this item? [Y/n]: ")
			input := readString(reader, "")
			if input == "" || strings.ToLower(input) == "y" {
				itemTotal := float64(templateItem.Quantity) * templateItem.UnitPrice
				items = append(items, BillItem{
					Description: templateItem.Description,
					Quantity:    templateItem.Quantity,
					UnitPrice:   templateItem.UnitPrice,
					Total:       itemTotal,
				})
				total += itemTotal
			}
		}
	}

	for {
		fmt.Println("\nAdd an item (press Enter without description to finish):")
		description := readString(reader, "Description: ")
		if description == "" {
			break
		}

		quantity := readInt(reader, "Quantity: ")
		unitPrice := readFloat(reader, "Unit Price (€): ")
		itemTotal := float64(quantity) * unitPrice

		items = append(items, BillItem{
			Description: description,
			Quantity:    quantity,
			UnitPrice:   unitPrice,
			Total:       itemTotal,
		})

		total += itemTotal
	}
	bill.Items = items
	bill.Total = total

	if template != nil && template.BitcoinAddress != "" {
		fmt.Printf("Bitcoin Address [%s]: ", template.BitcoinAddress)
		input := readString(reader, "")
		if input == "" {
			bill.BitcoinAddress = template.BitcoinAddress
		} else {
			bill.BitcoinAddress = input
		}
	} else {
		bill.BitcoinAddress = readString(reader, "Bitcoin Address: ")
	}

	return bill
}
