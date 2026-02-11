package excel

import (
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jung-kurt/gofpdf"
)

// Invoice Generate

func NewPDF() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()
	return pdf
}

func SetHeaderStyle(pdf *gofpdf.Fpdf, config *domain.PDFConfig) {
	if len(config.PrimaryColor) == 3 {
		pdf.SetFillColor(config.PrimaryColor[0], config.PrimaryColor[1], config.PrimaryColor[2])
	} else {
		pdf.SetFillColor(41, 128, 185) // Modern blue
	}
	pdf.SetTextColor(255, 255, 255)
	fontSize := config.HeaderFontSize
	if fontSize == 0 {
		fontSize = 20
	}
	pdf.SetFont("Arial", "B", fontSize)
}

func GetPDFConfig(ShowLogo bool, PrimaryColor, SecondaryColor, AccentColor []int, FontFamily string, HeaderFontSize, BodyFontSize, TableFontSize float64) *domain.PDFConfig {
	return &domain.PDFConfig{
		ShowLogo:       ShowLogo,
		PrimaryColor:   PrimaryColor,
		SecondaryColor: SecondaryColor,
		AccentColor:    AccentColor,
		FontFamily:     FontFamily,
		HeaderFontSize: HeaderFontSize,
		BodyFontSize:   BodyFontSize,
		TableFontSize:  TableFontSize,
		ShowSections:   []string{"header", "table", "reconciliation", "serviceFee", "notes"},
		CustomFields:   []domain.CustomField{},
	}
}
func GetDefaultPDFConfig() *domain.PDFConfig {
	return &domain.PDFConfig{
		ShowLogo:       false,
		PrimaryColor:   []int{41, 128, 185},  // Modern blue
		SecondaryColor: []int{52, 73, 94},    // Dark blue-gray
		AccentColor:    []int{149, 165, 166}, // Light gray
		FontFamily:     "Arial",
		HeaderFontSize: 20,
		BodyFontSize:   11,
		TableFontSize:  10,
		ShowSections:   []string{"header", "table", "reconciliation", "serviceFee", "notes"},
		CustomFields:   []domain.CustomField{},
	}
}

func SetTableRow(pdf *gofpdf.Fpdf, config *domain.PDFConfig, alt bool) {
	if alt {
		pdf.SetFillColor(245, 248, 250)
	} else {
		pdf.SetFillColor(255, 255, 255)
	}
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", config.TableFontSize)
}

func RenderHeader(pdf *gofpdf.Fpdf, d domain.Statement, config *domain.PDFConfig) {
	// Add logo if enabled
	if config.ShowLogo {
		pdf.Image("logo.png", 15, 15, 30, 0, false, "", 0, "")
		pdf.Ln(35)
	}

	SetHeaderStyle(pdf, config)
	pdf.CellFormat(0, 15, d.CompanyName, "", 1, "C", true, 0, "")
	pdf.Ln(10)

	SetNormalText(pdf, config)
	pdf.Cell(100, 8, "Provider: "+d.ProviderName)
	pdf.CellFormat(0, 8, "Date: "+d.Date, "", 1, "R", false, 0, "")

	pdf.Cell(100, 8, "Supplier Code: "+d.SupplierCode)
	pdf.CellFormat(0, 8, "Reference: "+d.Reference, "", 1, "R", false, 0, "")

	if d.GLCode != "" {
		pdf.Cell(100, 8, "GL Code: "+d.GLCode)
		pdf.CellFormat(0, 8, "Total Amount Payable: "+d.TotalPayable, "", 1, "R", false, 0, "")
	} else {
		pdf.CellFormat(0, 8, "Total Amount Payable: "+d.TotalPayable, "", 1, "R", false, 0, "")
	}
	pdf.Ln(12)
}

func SetNormalText(pdf *gofpdf.Fpdf, config *domain.PDFConfig) {
	pdf.SetTextColor(0, 0, 0)
	fontSize := config.BodyFontSize
	if fontSize == 0 {
		fontSize = 11
	}
	pdf.SetFont("Arial", "", fontSize)
}

func SetBoldText(pdf *gofpdf.Fpdf, config *domain.PDFConfig, size float64) {
	pdf.SetTextColor(0, 0, 0)
	if size == 0 {
		size = config.BodyFontSize
		if size == 0 {
			size = 11
		}
	}
	pdf.SetFont("Arial", "B", size)
}

func SetTableHeader(pdf *gofpdf.Fpdf, config *domain.PDFConfig) {
	if len(config.SecondaryColor) == 3 {
		pdf.SetFillColor(config.SecondaryColor[0], config.SecondaryColor[1], config.SecondaryColor[2])
	} else {
		pdf.SetFillColor(52, 73, 94) // Dark blue-gray
	}
	pdf.SetTextColor(255, 255, 255)
	fontSize := config.TableFontSize
	if fontSize == 0 {
		fontSize = 10
	}
	pdf.SetFont("Arial", "B", fontSize)
}

func setTableRow(pdf *gofpdf.Fpdf, config *domain.PDFConfig, alt bool) {
	if alt {
		pdf.SetFillColor(245, 248, 250) // Light gray
	} else {
		pdf.SetFillColor(255, 255, 255)
	}
	pdf.SetTextColor(0, 0, 0)
	fontSize := config.TableFontSize
	if fontSize == 0 {
		fontSize = 9
	}
	pdf.SetFont("Arial", "", fontSize)
}

func setSectionTitle(pdf *gofpdf.Fpdf, config *domain.PDFConfig) {
	if len(config.AccentColor) == 3 {
		pdf.SetFillColor(config.AccentColor[0], config.AccentColor[1], config.AccentColor[2])
	} else {
		pdf.SetFillColor(149, 165, 166) // Light gray
	}
	pdf.SetTextColor(255, 255, 255)
	fontSize := config.BodyFontSize + 2
	if fontSize < 13 {
		fontSize = 13
	}
	pdf.SetFont("Arial", "B", fontSize)
}
