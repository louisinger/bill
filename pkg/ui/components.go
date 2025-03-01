package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/louisinger/bill/pkg/bill"
)

func createHeaderWithSettings(title string, onSettings func()) *fyne.Container {
	settingsButton := widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), onSettings)
	settingsButton.Importance = widget.WarningImportance

	header := container.NewHBox(
		container.NewHBox(
			widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{
				Bold:      true,
				Monospace: true,
			}),
		),
		layout.NewSpacer(),
		settingsButton,
	)
	header.Resize(fyne.NewSize(1000, 60))
	return header
}

func createFormCard(title string, items ...*widget.FormItem) *widget.Card {
	return widget.NewCard("", "", container.NewVBox(
		widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		container.NewPadded(
			widget.NewForm(items...),
		),
	))
}

func createItemsTable(items []bill.BillItem, currency string, onDelete func(int)) *widget.Table {
	table := widget.NewTable(
		func() (int, int) { return len(items), 5 },
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

			if id.Row >= len(items) {
				label.SetText("")
				button.Hide()
				button.OnTapped = nil
				return
			}

			item := items[id.Row]
			switch id.Col {
			case 0:
				label.SetText(item.Description)
			case 1:
				label.SetText(fmt.Sprintf("%d", item.Quantity))
			case 2:
				label.SetText(fmt.Sprintf("%.2f %s", item.UnitPrice, currency))
			case 3:
				label.SetText(fmt.Sprintf("%.2f %s", item.Total, currency))
			case 4:
				label.SetText("")
				button.Show()
				button.OnTapped = func() {
					onDelete(id.Row)
				}
				button.Importance = widget.DangerImportance
			}
		},
	)

	// Set column widths
	table.SetColumnWidth(0, 400) // Description - wider
	table.SetColumnWidth(1, 100) // Quantity
	table.SetColumnWidth(2, 120) // Unit Price
	table.SetColumnWidth(3, 120) // Total
	table.SetColumnWidth(4, 80)  // Actions - slightly wider

	return table
}

func createFooterWithGenerateButton(onGenerate func()) *fyne.Container {
	generateButton := widget.NewButtonWithIcon("Generate PDF", theme.DocumentCreateIcon(), onGenerate)
	generateButton.Importance = widget.HighImportance
	generateButton.Resize(fyne.NewSize(200, 40))

	return container.NewHBox(
		layout.NewSpacer(),
		generateButton,
	)
}
