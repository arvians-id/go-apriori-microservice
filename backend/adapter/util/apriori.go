package util

import (
	"fmt"
	model2 "github.com/arvians-id/go-apriori-microservice/adapter/model"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func FindFirstItemSet(transactionsSet []*model2.Transaction, minimumSupport float32) ([]*model2.GetProductNameTransactionResponse, map[string]float32, []string) {
	// Generate all product
	var transactions []*model2.GetProductNameTransactionResponse
	for _, transaction := range transactionsSet {
		productName := strings.Split(transaction.ProductName, ", ")
		transactions = append(transactions, &model2.GetProductNameTransactionResponse{
			ProductName: productName,
		})
	}

	// Count all every product name
	var productName = make(map[string]float32)
	for _, value := range transactions {
		for _, product := range value.ProductName {
			product = strings.ToLower(product)
			productName[product] = productName[product] + 1
		}
	}

	// Finding one item set
	var propertyProduct []string
	for nameOfProduct, total := range productName {
		support := total / float32(len(transactionsSet)) * 100
		if support >= minimumSupport {
			supportValue := fmt.Sprintf("%.2f", support)
			totalValue := strconv.Itoa(int(total))

			propertyProduct = append(propertyProduct, nameOfProduct+":"+supportValue+"/"+totalValue+"*"+"Eligible")
		} else {
			supportValue := fmt.Sprintf("%.2f", support)
			totalValue := strconv.Itoa(int(total))

			propertyProduct = append(propertyProduct, nameOfProduct+":"+supportValue+"/"+totalValue+"*"+"Not Eligible")

		}
	}

	return transactions, productName, propertyProduct
}
func HandleMapsProblem(propertyProduct []string, minSupport float32) ([]string, []float32, []int32, []string, []string) {
	var oneSet []string
	var support []float32
	var totalTransaction []int32
	var checkEligible []string
	var cleanSet []string

	sort.Strings(propertyProduct)

	for i := 0; i < len(propertyProduct); i++ {
		// Split property
		nameOfProduct := strings.Split(propertyProduct[i], ":")
		transaction := strings.Split(nameOfProduct[1], "/")
		isEligible := strings.Split(transaction[1], "*")

		// Insert product name
		oneSet = append(oneSet, nameOfProduct[0])

		// Convert and insert support
		number, _ := strconv.ParseFloat(transaction[0], 64)
		support = append(support, float32(number))

		if float32(number) >= minSupport {
			cleanSet = append(cleanSet, nameOfProduct[0])
		}

		// Convert and insert total transaction
		transactionNumber, _ := strconv.Atoi(isEligible[0])
		totalTransaction = append(totalTransaction, int32(transactionNumber))

		// Check Is Eligible
		checkEligible = append(checkEligible, isEligible[1])
	}

	return oneSet, support, totalTransaction, checkEligible, cleanSet
}

func FindConfidence(apriori []*model2.GenerateApriori, productName map[string]float32, minSupport float32, minConfidence float32) []*model2.GenerateApriori {
	var confidence []*model2.GenerateApriori
	for _, value := range apriori {
		if value.Iterate == apriori[len(apriori)-1].Iterate {
			if val, ok := productName[value.ItemSet[0]]; ok && value.Support >= float32(minSupport) && float32(value.Transaction)/val*100 >= minConfidence {
				confidence = append(confidence, &model2.GenerateApriori{
					ItemSet:     value.ItemSet,
					Support:     value.Support,
					Iterate:     value.Iterate,
					Transaction: value.Transaction,
					Confidence:  float32(float32(value.Transaction) / val * 100),
				})
			}
		}
	}

	return confidence
}

func IsDuplicate(array []string) bool {
	visited := make(map[string]bool, 0)
	for i := 0; i < len(array); i++ {
		if visited[array[i]] == true {
			return true
		} else {
			visited[array[i]] = true
		}

	}
	return false
}

func FindCandidate(data []string, transactions []*model2.GetProductNameTransactionResponse) int32 {
	var counter int32
	for _, j := range transactions {
		results := make([]string, 0) // slice to store the result

		for i := 0; i < len(data); i++ {
			for k := 0; k < len(j.ProductName); k++ {
				if data[i] != j.ProductName[k] {
					continue
				}
				// append a value in result only if
				// it exists both in names and board
				results = append(results, data[i])
			}
			if len(results) == len(data) {
				counter++
			}
		}
	}
	return counter
}

func FindDiscount(apriori []*model2.GenerateApriori, minDiscount float32, maxDiscount float32) []*model2.GenerateApriori {
	var discounts []*model2.GenerateApriori
	var calculateDiscount = (maxDiscount - minDiscount) / float32(len(apriori))

	// Sorting if the value is greater, then the discount given will be large
	sort.Slice(apriori, func(i, j int) bool {
		return apriori[i].Confidence < apriori[j].Confidence
	})

	for _, value := range apriori {
		minDiscount += calculateDiscount
		discounts = append(discounts, &model2.GenerateApriori{
			ItemSet:     value.ItemSet,
			Support:     value.Support,
			Iterate:     value.Iterate,
			Transaction: value.Transaction,
			Confidence:  value.Confidence,
			Discount:    float32(minDiscount),
		})
	}

	return discounts
}

func FilterCandidateInSlice(dataTemp [][]string) [][]string {
	for i := 0; i < len(dataTemp); i++ {
		for j := i + 1; j < len(dataTemp); j++ {
			sort.Strings(dataTemp[i])
			sort.Strings(dataTemp[j])

			filter := reflect.DeepEqual(dataTemp[i], dataTemp[j])
			if filter {
				dataTemp = append(dataTemp[:j], dataTemp[j+1:]...)
				j--
			}
		}
	}

	return dataTemp
}
