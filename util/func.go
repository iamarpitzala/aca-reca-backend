package util

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func JSONResponse(c *gin.Context, status int, message string, data interface{}, err error) {
	resp := gin.H{
		"success": err == nil,
		"message": message,
	}

	if err != nil {
		resp["error"] = err.Error()
	} else {
		resp["data"] = data
	}

	c.JSON(status, resp)
}

func StructToMap(v interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}
