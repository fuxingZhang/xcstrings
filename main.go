package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

type StringsDict struct {
	SourceLanguage string                       `json:"sourceLanguage"`
	Version        string                       `json:"version"`
	Strings        map[string]LocalizationGroup `json:"strings"`
}

type LocalizationGroup struct {
	Comment         *string                     `json:"comment"`
	ExtractionState *string                     `json:"extractionState"`
	Localizations   map[string]LocalizationUnit `json:"localizations"`
}

type LocalizationUnit struct {
	StringUnit    *StringUnit                  `json:"stringUnit"`
	Variations    *VariationsUnit              `json:"variations"`
	Substitutions map[string]SubstitutionsUnit `json:"substitutions"`
}

type StringUnit struct {
	State string `json:"state"`
	Value string `json:"value"`
}

type VariationsUnit struct {
	Plural *PluralVariation `json:"plural"`
	Device *DeviceVariation `json:"device"`
}

type SubstitutionsUnit struct {
	FormatSpecifier string         `json:"formatSpecifier"`
	Variations      VariationsUnit `json:"variations"`
}

type PluralVariation struct {
	Zero  *StringUnit `json:"zero"`
	One   *StringUnit `json:"one"`
	Two   *StringUnit `json:"two"`
	Few   *StringUnit `json:"few"`
	Many  *StringUnit `json:"many"`
	Other *StringUnit `json:"other"`
}

type DeviceVariation struct {
	Appletv     *StringUnit `json:"appletv"`
	Applevision *StringUnit `json:"applevision"`
	Applewatch  *StringUnit `json:"applewatch"`
	Ipad        *StringUnit `json:"ipad"`
	Iphone      *StringUnit `json:"iphone"`
	Ipod        *StringUnit `json:"ipod"`
	Mac         *StringUnit `json:"mac"`
	Other       *StringUnit `json:"other"`
}

func translateString(stringToTranslate, targetLanguage string) (string, error) {
	subscriptionKey := "YOUR_SUBSCRIPTION_KEY"
	endpoint := "https://api.cognitive.microsofttranslator.com"

	// 构建请求
	url := fmt.Sprintf("%s/translate?api-version=3.0&to=%s", endpoint, targetLanguage)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, strings.NewReader(`[{"Text": "`+stringToTranslate+`"}]`))
	req.Header.Add("Ocp-Apim-Subscription-Key", subscriptionKey)
	req.Header.Add("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(body))

	// var result TranslationResult
	// err = json.Unmarshal(body, &result)
	// if err != nil {
	// 	return "", err
	// }

	// if len(result.Translations) > 0 {
	// 	return result.Translations[0].Text, nil
	// }

	return "", fmt.Errorf("no translation received")
}

func main() {
	jsonPath := "Localizable.xcstrings"

	// 读取文件
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	// 解析JSON
	var xStrings StringsDict
	if err := json.Unmarshal(jsonData, &xStrings); err != nil {
		log.Fatalf("Failed to parse JSON: %s", err)
	}

	keys := make([]string, 0, len(xStrings.Strings))
	// 遍历所有字符串
	for key, _ := range xStrings.Strings {
		//fmt.Println(key)
		keys = append(keys, key)
		// // 如果存在本地化条目
		// if localizationGroup.Localizations != nil {
		// 	for language, localizationUnit := range localizationGroup.Localizations {
		// 		fmt.Printf("%s: %+v\n", language, localizationUnit)
		// 		os.Exit(1)
		// 		// 如果字符串单元存在且没有翻译
		// 		if localizationUnit.StringUnit != nil && localizationUnit.StringUnit.Value == "" && localizationUnit.StringUnit.State != "translated" {
		// 			// 	// 翻译
		// 			// 	translatedValue, err := translateString(localizationUnit.StringUnit.Value, language)
		// 			// 	if err != nil {
		// 			// 		log.Printf("Translation failed for key: %s, error: %v", key, err)
		// 			// 		continue
		// 			// 	}
		// 			// 	localizationUnit.StringUnit.Value = translatedValue
		// 		}
		// 	}
		// }
	}

	// 按照字典顺序对键进行排序
	sort.Strings(keys)

	// 输出排序后的键及其对应的值
	//fmt.Println("Sorted keys:")
	//fmt.Println(strings.Join(keys, "\n"))
	for _, v := range keys {
		fmt.Println("|" + v + "| |")
	}

	// 将翻译后的JSON输出到文件或控制台
	// updatedJson, err := json.MarshalIndent(xStrings, "", "    ")
	// if err != nil {
	// 	log.Fatalf("Failed to marshal updated strings: %s", err)
	// }

	// // 输出翻译后的JSON到控制台
	// fmt.Println(string(updatedJson))

	// // 保存修改后的JSON文件
	// if err := os.WriteFile(jsonPath, updatedJson, 0644); err != nil {
	// 	log.Fatalf("Failed to write JSON back to file: %s", err)
	// }
}
