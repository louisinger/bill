package ui

import (
	"encoding/json"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DefaultParameters struct {
	// Company details
	CompanyName    string `json:"company_name"`
	Address        string `json:"address"`
	VATNumber      string `json:"vat_number"`
	BitcoinAddress string `json:"bitcoin_address"`
	Currency       string `json:"currency"`

	// Client details
	ToCompanyName string `json:"to_company_name"`
	ToAddress     string `json:"to_address"`
	ToVATNumber   string `json:"to_vat_number"`

	// Invoice details
	BillNumber string `json:"bill_number"`
}

func (ba *BillApp) showSettingsDialog() {
	// Create client address multi-line entry
	defaultToAddress := widget.NewMultiLineEntry()
	defaultToAddress.SetPlaceHolder("Default Client Address")
	defaultToAddress.Resize(fyne.NewSize(300, 100))

	// Create form content
	form := container.NewVBox(
		widget.NewCard("Default Company Details", "", widget.NewForm(
			widget.NewFormItem("Company Name", ba.defaultCompanyName),
			widget.NewFormItem("Address", ba.defaultAddress),
			widget.NewFormItem("VAT Number", ba.defaultVatNumber),
		)),
		widget.NewCard("Default Client Details", "", widget.NewForm(
			widget.NewFormItem("Client Company Name", widget.NewEntry()),
			widget.NewFormItem("Client Address", defaultToAddress),
			widget.NewFormItem("Client VAT Number", widget.NewEntry()),
		)),
		widget.NewCard("Default Invoice Details", "", widget.NewForm(
			widget.NewFormItem("Bill Number", widget.NewEntry()),
			widget.NewFormItem("Bitcoin Address", ba.defaultBitcoinAddress),
			widget.NewFormItem("Currency", ba.defaultCurrency),
		)),
	)

	// Store references to new default fields
	defaultToCompanyName := form.Objects[1].(*widget.Card).Content.(*widget.Form).Items[0].Widget.(*widget.Entry)
	defaultToVatNumber := form.Objects[1].(*widget.Card).Content.(*widget.Form).Items[2].Widget.(*widget.Entry)
	defaultBillNumber := form.Objects[2].(*widget.Card).Content.(*widget.Form).Items[0].Widget.(*widget.Entry)

	// Set placeholders
	defaultToCompanyName.SetPlaceHolder("Default Client Company Name")
	defaultToVatNumber.SetPlaceHolder("Default Client VAT Number")
	defaultBillNumber.SetPlaceHolder("INV-{YYYY}-{NNN}")

	// Create buttons
	saveDefaultsButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		params := DefaultParameters{
			CompanyName:    ba.defaultCompanyName.Text,
			Address:        ba.defaultAddress.Text,
			VATNumber:      ba.defaultVatNumber.Text,
			BitcoinAddress: ba.defaultBitcoinAddress.Text,
			Currency:       ba.defaultCurrency.Text,
			ToCompanyName:  defaultToCompanyName.Text,
			ToAddress:      defaultToAddress.Text,
			ToVATNumber:    defaultToVatNumber.Text,
			BillNumber:     defaultBillNumber.Text,
		}

		if err := ba.saveDefaultParametersWithData(params); err != nil {
			dialog.ShowError(err, ba.window)
			return
		}
		notification := fyne.NewNotification("Success", "Default values saved")
		ba.app.SendNotification(notification)
	})
	saveDefaultsButton.Importance = widget.HighImportance

	applyDefaultsButton := widget.NewButtonWithIcon("Apply to Current", theme.ConfirmIcon(), func() {
		ba.companyName.SetText(ba.defaultCompanyName.Text)
		ba.address.SetText(ba.defaultAddress.Text)
		ba.vatNumber.SetText(ba.defaultVatNumber.Text)
		ba.bitcoinAddress.SetText(ba.defaultBitcoinAddress.Text)
		ba.currency.SetText(ba.defaultCurrency.Text)
		ba.toCompanyName.SetText(defaultToCompanyName.Text)
		ba.toAddress.SetText(defaultToAddress.Text)
		ba.toVatNumber.SetText(defaultToVatNumber.Text)
		ba.billNumber.SetText(defaultBillNumber.Text)
		notification := fyne.NewNotification("Success", "Default values applied")
		ba.app.SendNotification(notification)
	})
	applyDefaultsButton.Importance = widget.WarningImportance

	// Add buttons to form
	form.Add(container.NewHBox(
		layout.NewSpacer(),
		applyDefaultsButton,
		saveDefaultsButton,
	))

	// Create and show dialog
	w := ba.app.NewWindow("Settings")
	w.SetContent(container.NewPadded(form))
	w.Resize(fyne.NewSize(600, 700))
	w.Show()

	// Load existing values into the new fields
	configDir, _ := os.UserConfigDir()
	data, err := os.ReadFile(filepath.Join(configDir, "invoice-generator", "defaults.json"))
	if err == nil {
		var params DefaultParameters
		if json.Unmarshal(data, &params) == nil {
			defaultToCompanyName.SetText(params.ToCompanyName)
			defaultToAddress.SetText(params.ToAddress)
			defaultToVatNumber.SetText(params.ToVATNumber)
			defaultBillNumber.SetText(params.BillNumber)
		}
	}
}

func (ba *BillApp) saveDefaultParametersWithData(params DefaultParameters) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(configDir, "invoice-generator")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(appDir, "defaults.json"), data, 0644)
}

func (ba *BillApp) loadDefaultParameters() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filepath.Join(configDir, "invoice-generator", "defaults.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, not an error
		}
		return err
	}

	var params DefaultParameters
	if err := json.Unmarshal(data, &params); err != nil {
		return err
	}

	// Set values in default fields
	ba.defaultCompanyName.SetText(params.CompanyName)
	ba.defaultAddress.SetText(params.Address)
	ba.defaultVatNumber.SetText(params.VATNumber)
	ba.defaultBitcoinAddress.SetText(params.BitcoinAddress)
	ba.defaultCurrency.SetText(params.Currency)

	// Always apply default values to form fields
	ba.companyName.SetText(params.CompanyName)
	ba.address.SetText(params.Address)
	ba.vatNumber.SetText(params.VATNumber)
	ba.bitcoinAddress.SetText(params.BitcoinAddress)
	ba.currency.SetText(params.Currency)
	ba.toCompanyName.SetText(params.ToCompanyName)
	ba.toAddress.SetText(params.ToAddress)
	ba.toVatNumber.SetText(params.ToVATNumber)
	ba.billNumber.SetText(params.BillNumber)

	return nil
}
