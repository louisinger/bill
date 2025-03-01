package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/louisinger/bill/pkg/bill"
)

type BillApp struct {
	app    fyne.App
	window fyne.Window

	// Form fields
	billNumber     *widget.Entry
	companyName    *widget.Entry
	address        *widget.Entry
	vatNumber      *widget.Entry
	toCompanyName  *widget.Entry
	toAddress      *widget.Entry
	toVatNumber    *widget.Entry
	bitcoinAddress *widget.Entry
	currency       *widget.Entry

	// Default values
	defaultCompanyName    *widget.Entry
	defaultAddress        *widget.Entry
	defaultVatNumber      *widget.Entry
	defaultBitcoinAddress *widget.Entry
	defaultCurrency       *widget.Entry

	// Items table
	items      []bill.BillItem
	itemList   *widget.Table
	addButton  *widget.Button
	totalLabel *widget.Label
}

func NewBillApp() *BillApp {
	a := app.New()
	w := a.NewWindow("Bill")

	ba := &BillApp{
		app:    a,
		window: w,
		items:  make([]bill.BillItem, 0),
	}

	// Initialize default value fields
	ba.defaultCompanyName = widget.NewEntry()
	ba.defaultCompanyName.SetPlaceHolder("Default Company Name")

	ba.defaultAddress = widget.NewMultiLineEntry()
	ba.defaultAddress.SetPlaceHolder("Default Company Address")
	ba.defaultAddress.Resize(fyne.NewSize(300, 100))

	ba.defaultVatNumber = widget.NewEntry()
	ba.defaultVatNumber.SetPlaceHolder("Default VAT Number")

	ba.defaultBitcoinAddress = widget.NewEntry()
	ba.defaultBitcoinAddress.SetPlaceHolder("Default Bitcoin Address")

	ba.defaultCurrency = widget.NewEntry()
	ba.defaultCurrency.SetPlaceHolder("€")

	// Create UI
	ba.createUI()

	// Load default parameters after UI is created
	if err := ba.loadDefaultParameters(); err != nil {
		dialog.ShowError(err, ba.window)
	}

	return ba
}

func (ba *BillApp) createUI() {
	// Create form fields with styling
	ba.billNumber = &widget.Entry{
		PlaceHolder: "INV-2024-001",
		TextStyle:   fyne.TextStyle{Monospace: true},
	}
	ba.billNumber.ExtendBaseWidget(ba.billNumber)
	ba.billNumber.Resize(fyne.NewSize(200, 35))

	ba.companyName = widget.NewEntry()
	ba.companyName.SetPlaceHolder("Your Company Name")

	ba.address = widget.NewMultiLineEntry()
	ba.address.SetPlaceHolder("Your Company Address")
	ba.address.Resize(fyne.NewSize(300, 100))

	ba.vatNumber = widget.NewEntry()
	ba.vatNumber.SetPlaceHolder("Your VAT Number")

	ba.toCompanyName = widget.NewEntry()
	ba.toCompanyName.SetPlaceHolder("Client Company Name")

	ba.toAddress = widget.NewMultiLineEntry()
	ba.toAddress.SetPlaceHolder("Client Address")
	ba.toAddress.Resize(fyne.NewSize(300, 100))

	ba.toVatNumber = widget.NewEntry()
	ba.toVatNumber.SetPlaceHolder("Client VAT Number")

	ba.bitcoinAddress = widget.NewEntry()
	ba.bitcoinAddress.SetPlaceHolder("Bitcoin Address")

	ba.currency = widget.NewEntry()
	ba.currency.SetText("€")
	ba.currency.Resize(fyne.NewSize(50, 35))

	// Create items table with direct reference to ba.items
	ba.itemList = widget.NewTable(
		func() (int, int) { return len(ba.items), 5 },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Wide Content"),
				widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
			)
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			box := cell.(*fyne.Container)
			label := box.Objects[0].(*widget.Label)
			button := box.Objects[1].(*widget.Button)

			// Hide button by default
			button.Hide()

			if id.Row == -1 {
				// Header row
				label.TextStyle = fyne.TextStyle{Bold: true}
				switch id.Col {
				case 0:
					label.SetText("Description")
				case 1:
					label.SetText("Quantity")
				case 2:
					label.SetText("Unit Price")
				case 3:
					label.SetText("Total")
				case 4:
					label.SetText("Actions")
				}
				return
			}

			if id.Row >= len(ba.items) {
				label.SetText("")
				button.Hide()
				button.OnTapped = nil
				return
			}

			item := ba.items[id.Row]
			switch id.Col {
			case 0:
				label.SetText(item.Description)
			case 1:
				label.SetText(fmt.Sprintf("%d", item.Quantity))
			case 2:
				label.SetText(fmt.Sprintf("%.2f %s", item.UnitPrice, ba.currency.Text))
			case 3:
				label.SetText(fmt.Sprintf("%.2f %s", item.Total, ba.currency.Text))
			case 4:
				label.SetText("")
				button.Show()
				button.OnTapped = func() {
					ba.deleteItem(id.Row)
				}
				button.Importance = widget.DangerImportance
			}
		},
	)

	// Set column widths
	ba.itemList.SetColumnWidth(0, 400) // Description - wider
	ba.itemList.SetColumnWidth(1, 100) // Quantity
	ba.itemList.SetColumnWidth(2, 120) // Unit Price
	ba.itemList.SetColumnWidth(3, 120) // Total
	ba.itemList.SetColumnWidth(4, 80)  // Actions - slightly wider

	// Create total label
	ba.totalLabel = widget.NewLabelWithStyle("Total: 0.00 "+ba.currency.Text, fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	// Create add item button with icon and styling
	ba.addButton = widget.NewButtonWithIcon("Add Item", theme.ContentAddIcon(), ba.showAddItemDialog)
	ba.addButton.Importance = widget.HighImportance

	// Create header with app title and settings
	header := createHeaderWithSettings("Bill", ba.showSettingsDialog)

	// Create form layout with sections and improved spacing
	invoiceDetails := widget.NewCard("", "", container.NewPadded(
		container.NewVBox(
			widget.NewLabelWithStyle("Invoice Details", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			container.NewPadded(
				container.NewGridWithColumns(2,
					container.NewVBox(
						widget.NewLabelWithStyle("Bill Number", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						ba.billNumber,
					),
					container.NewVBox(
						widget.NewLabelWithStyle("Currency", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						ba.currency,
					),
				),
			),
		),
	))

	companyDetails := createFormCard("Your Company",
		widget.NewFormItem("Company Name", ba.companyName),
		widget.NewFormItem("Address", ba.address),
		widget.NewFormItem("VAT Number", ba.vatNumber),
	)

	clientDetails := createFormCard("Client Details",
		widget.NewFormItem("Company Name", ba.toCompanyName),
		widget.NewFormItem("Address", ba.toAddress),
		widget.NewFormItem("VAT Number", ba.toVatNumber),
	)

	paymentDetails := createFormCard("Payment",
		widget.NewFormItem("Bitcoin Address", ba.bitcoinAddress),
	)

	itemsCard := widget.NewCard("", "", container.NewVBox(
		widget.NewLabelWithStyle("Items", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		container.NewPadded(
			container.NewVBox(
				ba.itemList,
				container.NewHBox(
					ba.addButton,
					layout.NewSpacer(),
					ba.totalLabel,
				),
			),
		),
	))

	// Create main layout with improved spacing
	mainContent := container.NewVBox(
		header,
		container.NewPadded(
			container.NewVBox(
				invoiceDetails,
				container.NewGridWithColumns(2,
					companyDetails,
					clientDetails,
				),
				paymentDetails,
				itemsCard,
			),
		),
	)

	// Create scroll container
	scroll := container.NewVScroll(mainContent)

	// Create footer with generate button
	footer := createFooterWithGenerateButton(ba.generatePDF)

	// Create main container with padding
	content := container.NewBorder(
		nil,
		container.NewPadded(footer),
		nil,
		nil,
		scroll,
	)

	ba.window.SetContent(content)
	ba.window.Resize(fyne.NewSize(1200, 900))
}

func (ba *BillApp) Run() {
	ba.window.ShowAndRun()
}
