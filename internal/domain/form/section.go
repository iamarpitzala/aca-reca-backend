package form

import "github.com/iamarpitzala/aca-reca-backend/util"

type Section string

const (
	SectionIncome    Section = Section(util.SectionIncome)
	SectionExpense   Section = Section(util.SectionExpense)
	SectionReduction Section = Section(util.SectionReduction)
)
