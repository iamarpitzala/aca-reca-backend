package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

func BindAndValidate(c *gin.Context, rq interface{}) error {
	if err := c.ShouldBind(rq); err != nil {
		return err
	}

	if err := validator.New().Struct(rq); err != nil {
		return err
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()
	c.Set("user_agent", userAgent)
	c.Set("ip_address", ipAddress)

	return nil
}

// EXCL COMMON METHOD ////////////////////////////////////////////////////////

type ExcelStyles struct {
	Header int
	Text   map[string]int
	Number map[string]int
}

func SetupSheet(f *excelize.File, name string) {
	index, _ := f.NewSheet(name)
	f.SetActiveSheet(index)
}

func AlignmentCenter() *excelize.Alignment {
	return &excelize.Alignment{
		Horizontal: "center",
		Vertical:   "center",
	}
}
func FillColor(bg string) excelize.Fill {
	return excelize.Fill{Type: "pattern", Color: []string{bg}, Pattern: 1}
}

func SetAllBorders(color string) []excelize.Border {
	return []excelize.Border{
		{Type: "left", Style: 1, Color: color},
		{Type: "right", Style: 1, Color: color},
		{Type: "top", Style: 1, Color: color},
		{Type: "bottom", Style: 1, Color: color},
	}
}

func CreateStyles(f *excelize.File, currency string) (*ExcelStyles, error) {
	header, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{
			Type: "pattern", Color: []string{"#2F9EAA"}, Pattern: 1,
		},
		Alignment: AlignmentCenter(),
	})
	if err != nil {
		return nil, err
	}

	text := make(map[string]int)
	number := make(map[string]int)

	for _, bg := range []string{"#E6F5F8", "#FFFFFF"} {
		t, _ := f.NewStyle(&excelize.Style{
			Fill:      FillColor(bg),
			Alignment: AlignmentCenter(),
			Border:    SetAllBorders("#6cc3e9"),
		})

		n, _ := f.NewStyle(&excelize.Style{
			Fill:         FillColor(bg),
			Alignment:    AlignmentCenter(),
			Border:       SetAllBorders("#6cc3e9"),
			CustomNumFmt: &currency,
		})

		text[bg] = t
		number[bg] = n
	}

	return &ExcelStyles{
		Header: header,
		Text:   text,
		Number: number,
	}, nil
}

func WriteHeaders(
	f *excelize.File,
	sheet string,
	headers []string,
	headerStyle int,
) {
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	end, _ := excelize.CoordinatesToCellName(len(headers), 1)
	f.SetCellStyle(sheet, "A1", end, headerStyle)
}

func SetColWidths(f *excelize.File, sheet string, widths map[string]float64) {
	for col, w := range widths {
		f.SetColWidth(sheet, col, col, w)
	}
}

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
