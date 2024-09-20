package parsing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/inflection"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Recipe contains the info for the file and the lines
type Recipe struct {
	FileName    string       `json:"filename"`
	FileContent string       `json:"file_content"`
	Lines       []LineInfo   `json:"lines"`
	Ingredients []Ingredient `json:"ingredients"`
}

// LineInfo has all the information for the parsing of a given line
type LineInfo struct {
	LineOriginal        string
	Line                string         `json:",omitempty"`
	IngredientsInString []WordPosition `json:",omitempty"`
	AmountInString      []WordPosition `json:",omitempty"`
	MeasureInString     []WordPosition `json:",omitempty"`
	Ingredient          Ingredient     `json:",omitempty"`
}

// IngredientList is a list of ingredients
type IngredientList struct {
	Ingredients []Ingredient `json:"ingredients"`
}

// Ingredient is the basic struct for ingredients
type Ingredient struct {
	Name    string  `json:"name,omitempty"`
	Comment string  `json:"comment,omitempty"`
	Measure Measure `json:"measure,omitempty"`
	Line    string  `json:"line,omitempty"`
}

// Measure includes the amount, name and the cups for conversions
type Measure struct {
	Amount float64 `json:"amount"`
	Name   string  `json:"name"`
	Cups   float64 `json:"cups"`
	Weight float64 `json:"weight,omitempty"`
}

// WordPosition shows a word and its position
// Note: the position is memory-dependent as it will
// be the position after the last deleted word
type WordPosition struct {
	Word     string
	Position int
}

// Parse is the main parser for a given recipe.
func (r *Recipe) parseHTML() (rerr error) {
	if r == nil {
		r = &Recipe{}
	}
	if r.FileContent == "" || r.FileName == "" {
		rerr = fmt.Errorf("no file loaded")
		return
	}

	r.Lines, rerr = getIngredientLinesInHTML(r.FileContent)
	return r.parseRecipe()

}

func (r *Recipe) parseRecipe() (rerr error) {
	goodLines := make([]LineInfo, len(r.Lines))
	j := 0
	for _, lineInfo := range r.Lines {
		if len(strings.TrimSpace(lineInfo.Line)) < 3 || len(strings.TrimSpace(lineInfo.Line)) > 150 {
			continue
		}
		if strings.Contains(strings.ToLower(lineInfo.Line), "serving size") {
			continue
		}
		if strings.Contains(strings.ToLower(lineInfo.Line), "yield") {
			continue
		}

		// singularlize
		lineInfo.Ingredient.Measure = Measure{}

		// get amount, continue if there is an error
		err := lineInfo.getTotalAmount()
		if err != nil {
			continue
		}

		// get ingredient, continue if its not found
		err = lineInfo.getIngredient()
		if err != nil {
			continue
		}

		// get measure
		lineInfo.getMeasure()

		// get comment
		if len(lineInfo.MeasureInString) > 0 && len(lineInfo.IngredientsInString) > 0 {
			lineInfo.Ingredient.Comment = getOtherInBetweenPositions(lineInfo.Line, lineInfo.MeasureInString[0], lineInfo.IngredientsInString[0])
		}

		goodLines[j] = lineInfo
		j++
	}
	r.Lines = goodLines[:j]

	// consolidate ingredients
	ingredients := make(map[string]Ingredient)
	ingredientList := []string{}
	for _, line := range r.Lines {
		if _, ok := ingredients[line.Ingredient.Name]; ok {
			if ingredients[line.Ingredient.Name].Measure.Name == line.Ingredient.Measure.Name {
				ingredients[line.Ingredient.Name] = Ingredient{
					Name:    line.Ingredient.Name,
					Comment: ingredients[line.Ingredient.Name].Comment,
					Measure: Measure{
						Name:   ingredients[line.Ingredient.Name].Measure.Name,
						Amount: ingredients[line.Ingredient.Name].Measure.Amount + line.Ingredient.Measure.Amount,
						Cups:   ingredients[line.Ingredient.Name].Measure.Cups + line.Ingredient.Measure.Cups,
					},
				}
			} else {
				ingredients[line.Ingredient.Name] = Ingredient{
					Name:    line.Ingredient.Name,
					Comment: ingredients[line.Ingredient.Name].Comment,
					Measure: Measure{
						Name:   ingredients[line.Ingredient.Name].Measure.Name,
						Amount: ingredients[line.Ingredient.Name].Measure.Amount,
						Cups:   ingredients[line.Ingredient.Name].Measure.Cups + line.Ingredient.Measure.Cups,
					},
				}
			}
		} else {
			ingredientList = append(ingredientList, line.Ingredient.Name)
			ingredients[line.Ingredient.Name] = Ingredient{
				Name:    line.Ingredient.Name,
				Comment: line.Ingredient.Comment,
				Measure: Measure{
					Name:   line.Ingredient.Measure.Name,
					Amount: line.Ingredient.Measure.Amount,
					Cups:   line.Ingredient.Measure.Cups + line.Ingredient.Measure.Cups,
				},
			}
		}
	}
	r.Ingredients = make([]Ingredient, len(ingredients))
	for i, ing := range ingredientList {
		r.Ingredients[i] = ingredients[ing]
	}

	return
}

func (lineInfo *LineInfo) getTotalAmount() (err error) {
	lastPosition := -1
	totalAmount := 0.0
	wps := lineInfo.AmountInString
	for i := range wps {
		wps[i].Word = strings.TrimSpace(wps[i].Word)
		if lastPosition == -1 {
			totalAmount = ConvertStringToNumber(wps[i].Word)
		} else if math.Abs(float64(wps[i].Position-lastPosition)) < 6 {
			totalAmount += ConvertStringToNumber(wps[i].Word)
		}
		lastPosition = wps[i].Position + len(wps[i].Word)
	}
	if totalAmount == 0 && strings.Contains(lineInfo.Line, "whole") {
		totalAmount = 1
	}
	if totalAmount == 0 {
		err = fmt.Errorf("no amount found")
	} else {
		lineInfo.Ingredient.Measure.Amount = totalAmount
	}
	return
}

func (lineInfo *LineInfo) getIngredient() (err error) {
	if len(lineInfo.IngredientsInString) == 0 {
		err = fmt.Errorf("no ingredient found")
		return
	}
	lineInfo.Ingredient.Name = inflection.Singular(lineInfo.IngredientsInString[0].Word)
	return
}

func (lineInfo *LineInfo) getMeasure() (err error) {
	if len(lineInfo.MeasureInString) == 0 {
		lineInfo.Ingredient.Measure.Name = "whole"
		return
	}
	lineInfo.Ingredient.Measure.Name = lineInfo.MeasureInString[0].Word
	return
}

func getIngredientLinesInHTML(htmlS string) (lineInfos []LineInfo, err error) {
	doc, err := html.Parse(bytes.NewReader([]byte(htmlS)))
	if err != nil {
		return
	}
	var f func(n *html.Node, lineInfos *[]LineInfo) (s string, done bool)
	f = func(n *html.Node, lineInfos *[]LineInfo) (s string, done bool) {
		childrenLineInfo := []LineInfo{}
		score := 0
		isScript := n.DataAtom == atom.Script
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if isScript {
				// try to capture JSON and if successful, do a hard exit
				lis, errJSON := extractLinesFromJavascript(c.Data)
				if errJSON == nil && len(lis) > 2 {
					*lineInfos = lis
					done = true
					return
				}
			}
			var childText string
			childText, done = f(c, lineInfos)
			if done {
				return
			}
			if childText != "" {
				scoreOfLine, lineInfo := scoreLine(childText)
				childrenLineInfo = append(childrenLineInfo, lineInfo)
				score += scoreOfLine
			}
		}
		if score > 2 && len(childrenLineInfo) < 25 && len(childrenLineInfo) > 2 {
			*lineInfos = append(*lineInfos, childrenLineInfo...)
			for _, child := range childrenLineInfo {
				log.Println("[%s]", child.LineOriginal)
			}
		}
		if len(childrenLineInfo) > 0 {
			// fmt.Println(childrenLineInfo)
			childrenText := make([]string, len(childrenLineInfo))
			for i := range childrenLineInfo {
				childrenText[i] = childrenLineInfo[i].LineOriginal
			}
			s = strings.Join(childrenText, " ")
		} else if n.DataAtom == 0 && strings.TrimSpace(n.Data) != "" {
			s = strings.TrimSpace(n.Data)
		}
		return
	}
	f(doc, &lineInfos)
	return
}

func extractLinesFromJavascript(jsString string) (lineInfo []LineInfo, err error) {
	var arrayMap = []map[string]interface{}{}
	var regMap = make(map[string]interface{})
	err = json.Unmarshal([]byte(jsString), &regMap)
	if err != nil {
		err = json.Unmarshal([]byte(jsString), &arrayMap)
		if err != nil {
			return
		}
		if len(arrayMap) == 0 {
			err = fmt.Errorf("nothing to parse")
			return
		}
		parseMap(arrayMap[0], &lineInfo)
		err = nil
	} else {
		parseMap(regMap, &lineInfo)
		err = nil
	}

	return
}

func parseMap(aMap map[string]interface{}, lineInfo *[]LineInfo) {
	for _, val := range aMap {
		switch val.(type) {
		case map[string]interface{}:
			parseMap(val.(map[string]interface{}), lineInfo)
		case []interface{}:
			parseArray(val.([]interface{}), lineInfo)
		default:
			// fmt.Println(key, ":", concreteVal)
		}
	}
}

func parseArray(anArray []interface{}, lineInfo *[]LineInfo) {
	concreteLines := []string{}
	for _, val := range anArray {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			parseMap(val.(map[string]interface{}), lineInfo)
		case []interface{}:
			parseArray(val.([]interface{}), lineInfo)
		default:
			switch v := concreteVal.(type) {
			case string:
				concreteLines = append(concreteLines, v)
			}
		}
	}

	score, li := scoreLines(concreteLines)
	if score > 20 {
		*lineInfo = li
	}
}

func scoreLines(lines []string) (score int, lineInfo []LineInfo) {
	if len(lines) < 2 {
		return
	}
	lineInfo = make([]LineInfo, len(lines))
	for i, line := range lines {
		var scored int
		scored, lineInfo[i] = scoreLine(line)
		score += scored
	}
	return
}

func scoreLine(line string) (score int, lineInfo LineInfo) {
	lineInfo = LineInfo{}
	lineInfo.LineOriginal = line
	lineInfo.Line = SanitizeLine(line)
	lineInfo.IngredientsInString = GetIngredientsInString(lineInfo.Line)
	lineInfo.AmountInString = GetNumbersInString(lineInfo.Line)
	lineInfo.MeasureInString = GetMeasuresInString(lineInfo.Line)
	if len(lineInfo.IngredientsInString) == 2 && len(lineInfo.IngredientsInString[1].Word) > len(lineInfo.IngredientsInString[0].Word) {
		lineInfo.IngredientsInString[0] = lineInfo.IngredientsInString[1]
	}

	if len(lineInfo.LineOriginal) > 50 {
		return
	}

	// does it contain an ingredient?
	if len(lineInfo.IngredientsInString) > 0 {
		score++
	}

	// disfavor containing multiple ingredients
	if len(lineInfo.IngredientsInString) > 1 {
		score = score - len(lineInfo.IngredientsInString) + 1
	}

	// does it contain an amount?
	if len(lineInfo.AmountInString) > 0 {
		score++
	}
	// does it contain a measure (cups, tsps)?
	if len(lineInfo.MeasureInString) > 0 {
		score++
	}
	// does the ingredient come after the measure?
	if len(lineInfo.IngredientsInString) > 0 && len(lineInfo.MeasureInString) > 0 && lineInfo.IngredientsInString[0].Position > lineInfo.MeasureInString[0].Position {
		score++
	}
	// does the ingredient come after the amount?
	if len(lineInfo.IngredientsInString) > 0 && len(lineInfo.AmountInString) > 0 && lineInfo.IngredientsInString[0].Position > lineInfo.AmountInString[0].Position {
		score++
	}
	// does the measure come after the amount?
	if len(lineInfo.MeasureInString) > 0 && len(lineInfo.AmountInString) > 0 && lineInfo.MeasureInString[0].Position > lineInfo.AmountInString[0].Position {
		score++
	}

	// disfavor lots of puncuation
	puncuation := []string{".", ",", "!", "?"}
	for _, punc := range puncuation {
		if strings.Count(lineInfo.LineOriginal, punc) > 1 {
			score--
		}
	}

	// disfavor long lines
	if len(lineInfo.Line) > 30 {
		score = score - (len(lineInfo.Line) - 30)
	}
	if len(lineInfo.Line) > 250 {
		score = 0
	}

	// does it start with a list indicator (* or -)?
	fields := strings.Fields(lineInfo.Line)
	if len(fields) > 0 && (fields[0] == "*" || fields[0] == "-") {
		score++
	}
	// if only one thing is right, its wrong
	if score == 1 {
		score = 0.0
	}
	return
}

// GetIngredientsInString returns the word positions of the ingredients
func GetIngredientsInString(s string) (wordPositions []WordPosition) {
	return getWordPositions(s, corpusIngredients)
}

// GetNumbersInString returns the word positions of the numbers in the ingredient string
func GetNumbersInString(s string) (wordPositions []WordPosition) {
	return getWordPositions(s, corpusNumbers)
}

// GetMeasuresInString returns the word positions of the measures in a ingredient string
func GetMeasuresInString(s string) (wordPositions []WordPosition) {
	return getWordPositions(s, corpusMeasures)
}

func getWordPositions(s string, corpus []string) (wordPositions []WordPosition) {
	wordPositions = []WordPosition{}
	for _, ing := range corpus {
		pos := strings.Index(s, ing)
		if pos > -1 {
			s = strings.Replace(s, ing, strings.Repeat(" ", utf8.RuneCountInString(ing)), 1)
			ing = strings.TrimSpace(ing)
			wordPositions = append(wordPositions, WordPosition{ing, pos})
		}
	}
	sort.Slice(wordPositions, func(i, j int) bool {
		return wordPositions[i].Position < wordPositions[j].Position
	})
	return
}

func ConvertStringToNumber(s string) float64 {
	switch s {
	case "½":
		return 0.5
	case "¼":
		return 0.25
	case "¾":
		return 0.75
	case "⅛":
		return 1.0 / 8
	case "⅜":
		return 3.0 / 8
	case "⅝":
		return 5.0 / 8
	case "⅞":
		return 7.0 / 8
	case "⅔":
		return 2.0 / 3
	case "⅓":
		return 1.0 / 3
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// getOtherInBetweenPositions returns the word positions comment string in the ingredients
func getOtherInBetweenPositions(s string, pos1, pos2 WordPosition) (other string) {
	if pos1.Position > pos2.Position {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			errorWriter := gin.DefaultErrorWriter
			errorWriter.Write([]byte(fmt.Errorf(s, pos1, pos2).Error()))
			errorWriter.Write([]byte(fmt.Errorf("%v", r).Error()))
		}
	}()
	other = s[pos1.Position+len(pos1.Word)+1 : pos2.Position]
	other = strings.TrimSpace(other)
	return
}
