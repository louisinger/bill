package ui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/louisinger/bill/pkg/bill"
)

func (ba *BillApp) generatePDF() {
	// Validate required fields
	if ba.billNumber.Text == "" {
		dialog.ShowError(fmt.Errorf("bill number is required"), ba.window)
		return
	}
	if ba.companyName.Text == "" {
		dialog.ShowError(fmt.Errorf("your company name is required"), ba.window)
		return
	}
	if ba.toCompanyName.Text == "" {
		dialog.ShowError(fmt.Errorf("client company name is required"), ba.window)
		return
	}
	if len(ba.items) == 0 {
		dialog.ShowError(fmt.Errorf("at least one item is required"), ba.window)
		return
	}

	// Create a sanitized filename from the invoice number
	defaultFileName := strings.ReplaceAll(ba.billNumber.Text, "/", "-") + ".pdf"

	dialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, ba.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		// Create bill data
		b := bill.Bill{
			Number:         ba.billNumber.Text,
			Date:           time.Now(),
			CompanyName:    ba.companyName.Text,
			Address:        ba.address.Text,
			VATNumber:      ba.vatNumber.Text,
			ToCompanyName:  ba.toCompanyName.Text,
			ToAddress:      ba.toAddress.Text,
			ToVATNumber:    ba.toVatNumber.Text,
			Items:          ba.items,
			Currency:       ba.currency.Text,
			BitcoinAddress: ba.bitcoinAddress.Text,
		}

		// Calculate total
		total := 0.0
		for _, item := range ba.items {
			total += item.Total
		}
		b.Total = total

		// Generate PDF
		outputPath := writer.URI().Path()
		if filepath.Ext(outputPath) != ".pdf" {
			outputPath += ".pdf"
		}

		err = bill.GeneratePDF(b, outputPath)
		if err != nil {
			dialog.ShowError(err, ba.window)
			return
		}

		notification := fyne.NewNotification("Success", "PDF generated successfully")
		ba.app.SendNotification(notification)
	}, ba.window)

	dialog.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
	dialog.SetFileName(defaultFileName)
	dialog.Show()
}
