package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
	"github.com/xuri/excelize/v2"
)

type PayslipHandler struct {
}

func NewPayslipHandler() *PayslipHandler {
	return &PayslipHandler{}
}

// Export Excel Income
// POST /api/v1/payslip/export/income
// @Summary Export Excel Income
// @Description Export Excel Income
// @Tags Payslip
// @Accept json
// @Produce json
// @Param data body []domain.ExportIncome true "Data"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /payslip/export/income [post]
func (h *PayslipHandler) ExportExcelIncome(c *gin.Context) {
	var data []domain.ExportIncome
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}
	fileName, err := IncomeFormate(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Income Excel created successfully",
		"file":    fileName,
	})
}

// Export Excel Expenses
// POST /api/v1/payslip/export/expenses
// @Summary Export Excel Expenses
// @Description Export Excel Expenses
// @Tags Payslip
// @Accept json
// @Produce json
// @Param data body []domain.ExportExpenses true "Data"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /payslip/export/expenses [post]
func (h *PayslipHandler) ExportExcelExpanses(c *gin.Context) {
	var data []domain.ExportExpenses
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}

	fileName, err := ExpensesFormate(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Expenses Excel created successfully",
		"file":    fileName,
	})

}

// Generate PDF
// POST /api/v1/payslip/generate/pdf
// @Summary Generate PDF
// @Description Generate PDF
// @Tags Payslip
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /payslip/generate/pdf [post]
func (h *PayslipHandler) GeneratePdf(c *gin.Context) {
	fileName, err := GenerateInvoiceOne()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "PDF created successfully",
		"file":    fileName,
	})
}

// income formate
func IncomeFormate(data []domain.ExportIncome) (string, error) {
	f := excelize.NewFile()
	sheet := "Users_Income"

	utils.SetupSheet(f, sheet)

	styles, err := utils.CreateStyles(f, "$#,##0.00")
	if err != nil {
		return "", err
	}

	utils.WriteHeaders(f, sheet, domain.HeadersIncome, styles.Header)

	for i, d := range data {
		row := i + 2
		bg := "#E6F5F8"
		if i%2 != 0 {
			bg = "#FFFFFF"
		}

		// Text
		f.SetCellValue(sheet, "A"+fmt.Sprint(row), d.Q)
		f.SetCellValue(sheet, "B"+fmt.Sprint(row), d.IncomeType)
		f.SetCellValue(sheet, "C"+fmt.Sprint(row), d.LabRecord)
		f.SetCellValue(sheet, "D"+fmt.Sprint(row), d.PaymentDate)
		f.SetCellValue(sheet, "E"+fmt.Sprint(row), d.DentalPractice)
		f.SetCellStyle(sheet, "A"+fmt.Sprint(row), "E"+fmt.Sprint(row), styles.Text[bg])

		// Numbers
		values := []struct {
			col string
			val float64
		}{
			{"F", d.Adjustments}, {"G", d.GrossIncomeG1},
			{"H", d.LabFees}, {"I", d.GrossNetLabFees},
			{"J", d.GSTPayable1A}, {"K", d.GSTFree},
			{"L", d.ManagementFeesG11}, {"M", d.Percentage},
			{"N", d.GSTRefundable1B}, {"O", d.NetPayment},
		}

		for _, v := range values {
			cell := v.col + fmt.Sprint(row)
			f.SetCellValue(sheet, cell, v.val)
			f.SetCellStyle(sheet, cell, cell, styles.Number[bg])
		}
	}

	var incomeColWidths = map[string]float64{
		"A": 6, "B": 15, "C": 18, "D": 15, "E": 30,
		"F": 12, "G": 12, "H": 12, "I": 12, "J": 12,
		"K": 12, "L": 12, "M": 12, "N": 12, "O": 12,
	}

	utils.SetColWidths(f, sheet, incomeColWidths)
	f.AutoFilter(sheet, "A1:O1", nil)
	file := "Income.xlsx"
	return file, f.SaveAs(file)
}

// expenses formate
func ExpensesFormate(data []domain.ExportExpenses) (string, error) {
	f := excelize.NewFile()
	sheet := "Users_Expenses"

	utils.SetupSheet(f, sheet)

	styles, err := utils.CreateStyles(f, "$#,##0.00")
	if err != nil {
		return "", err
	}

	utils.WriteHeaders(f, sheet, domain.HeadersExpenses, styles.Header)

	for i, d := range data {
		row := i + 2
		bg := "#E6F5F8"
		if i%2 != 0 {
			bg = "#FFFFFF"
		}

		// Text
		f.SetCellValue(sheet, "A"+fmt.Sprint(row), d.Q)
		f.SetCellValue(sheet, "B"+fmt.Sprint(row), d.Date)
		f.SetCellValue(sheet, "C"+fmt.Sprint(row), d.Supplier)
		f.SetCellValue(sheet, "D"+fmt.Sprint(row), d.Category)
		f.SetCellValue(sheet, "I"+fmt.Sprint(row), d.Remarks)

		f.SetCellStyle(sheet, "A"+fmt.Sprint(row), "D"+fmt.Sprint(row), styles.Text[bg])
		f.SetCellStyle(sheet, "I"+fmt.Sprint(row), "I"+fmt.Sprint(row), styles.Text[bg])

		// Numbers
		nums := map[string]float64{
			"E": d.Amount,
			"F": d.Bas,
			"G": d.Net,
			"H": d.GST,
		}

		for col, val := range nums {
			cell := col + fmt.Sprint(row)
			f.SetCellValue(sheet, cell, val)
			f.SetCellStyle(sheet, cell, cell, styles.Number[bg])
		}
	}
	var expensesColWidths = map[string]float64{
		"A": 6, "B": 15, "C": 18, "D": 35,
		"E": 14, "F": 14, "G": 14, "H": 14,
		"I": 50,
	}
	utils.SetColWidths(f, sheet, expensesColWidths)
	f.AutoFilter(sheet, "A1:I1", nil)

	file := "Expenses.xlsx"
	return file, f.SaveAs(file)
}

// pdf formatting
func GenerateInvoiceOne() (string, error) {
	// pdf := utils.NewPDF()
	// config := utils.GetDefaultPDFConfig()

	return "", nil
}
