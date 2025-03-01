package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/louisinger/bill/pkg/bill"
)

func (ba *BillApp) deleteItem(index int) {
	if index < 0 || index >= len(ba.items) {
		return
	}
	ba.items = append(ba.items[:index], ba.items[index+1:]...)
	ba.updateTotal()
	ba.itemList.Refresh()
}

func (ba *BillApp) updateTotal() {
	total := 0.0
	for _, item := range ba.items {
		total += item.Total
	}
	ba.totalLabel.SetText(fmt.Sprintf("Total: %.2f %s", total, ba.currency.Text))
}

func (ba *BillApp) showAddItemDialog() {
	description := widget.NewMultiLineEntry()
	description.SetPlaceHolder("Item Description")
	description.Resize(fyne.NewSize(500, 150))

	quantity := widget.NewEntry()
	quantity.SetPlaceHolder("Quantity")
	quantity.Resize(fyne.NewSize(200, 35))

	unitPrice := widget.NewEntry()
	unitPrice.SetPlaceHolder("Unit Price")
	unitPrice.Resize(fyne.NewSize(200, 35))

	// Create a custom form with larger spacing
	form := container.NewVBox(
		widget.NewLabelWithStyle("Description", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		description,
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			container.NewVBox(
				widget.NewLabelWithStyle("Quantity", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				quantity,
			),
			container.NewVBox(
				widget.NewLabelWithStyle("Unit Price", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				unitPrice,
			),
		),
	)

	// Create and show a custom dialog
	customDialog := dialog.NewCustomConfirm("Add Item", "Add", "Cancel", form, func(confirm bool) {
		if !confirm {
			return
		}

		if description.Text == "" {
			dialog.ShowError(fmt.Errorf("description is required"), ba.window)
			return
		}

		qty, err := strconv.Atoi(quantity.Text)
		if err != nil || qty <= 0 {
			dialog.ShowError(fmt.Errorf("invalid quantity (must be a positive number)"), ba.window)
			return
		}

		price, err := strconv.ParseFloat(unitPrice.Text, 64)
		if err != nil || price < 0 {
			dialog.ShowError(fmt.Errorf("invalid unit price (must be a non-negative number)"), ba.window)
			return
		}

		total := float64(qty) * price
		ba.items = append(ba.items, bill.BillItem{
			Description: description.Text,
			Quantity:    qty,
			UnitPrice:   price,
			Total:       total,
		})
		ba.updateTotal()
		ba.itemList.Refresh()

		// Update the table's data function to reflect the new item count
		ba.itemList.Length = func() (int, int) {
			return len(ba.items), 5
		}
	}, ba.window)

	customDialog.Resize(fyne.NewSize(600, 400))
	customDialog.Show()
}
