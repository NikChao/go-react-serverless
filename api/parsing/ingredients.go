package parsing

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// NewFromHTML generates a new parser from a HTML text
func NewFromHTML(name, htmlstring string) (r *Recipe, err error) {
	r = &Recipe{FileName: name}
	r.FileContent = htmlstring
	err = r.parseHTML()
	return
}

func NewFromURL(url string) (r *Recipe, err error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return NewFromHTML(url, string(html))
}

func SanitizeLine(s string) string {
	s = strings.ToLower(s)
	s = strings.Replace(s, "⁄", "/", -1)
	s = strings.Replace(s, " / ", "/", -1)

	// special cases
	s = strings.Replace(s, "butter milk", "buttermilk", -1)
	s = strings.Replace(s, "bicarbonate of soda", "baking soda", -1)
	s = strings.Replace(s, "soda bicarbonate", "baking soda", -1)

	// remove parentheses
	re := regexp.MustCompile(`(?s)\((.*)\)`)
	for _, m := range re.FindAllStringSubmatch(s, -1) {
		s = strings.Replace(s, m[0], " ", 1)
	}

	s = " " + strings.TrimSpace(s) + " "

	// replace unicode fractions with fractions
	for v := range corpusFractionNumberMap {
		s = strings.Replace(s, v, " "+corpusFractionNumberMap[v].fractionString+" ", -1)
	}

	// remove non-alphanumeric
	reg, _ := regexp.Compile("[^a-zA-Z0-9/.]+")
	s = reg.ReplaceAllString(s, " ")

	// replace fractions with unicode fractions
	for v := range corpusFractionNumberMap {
		s = strings.Replace(s, corpusFractionNumberMap[v].fractionString, " "+v+" ", -1)
	}

	s = strings.Replace(s, " one ", " 1 ", -1)

	return s
}

// IngredientList will return a string containing the ingredient list
func (r *Recipe) IngredientList() (ingredientList IngredientList) {
	ingredientList = IngredientList{make([]Ingredient, len(r.Lines))}
	for i, li := range r.Lines {
		ingredientList.Ingredients[i] = li.Ingredient
		ingredientList.Ingredients[i].Line = li.LineOriginal
	}
	return
}
