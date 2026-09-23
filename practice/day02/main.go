package main

import (
	"fmt"
)

type Record struct {
	ID     int
	Title  string
	Artist string
	Price  float64
}

func addRecord(records []Record, record Record) []Record {
	records = append(records, record)
	return records
}
func findRecordByID(records []Record, id int) (Record, bool) {
	for i := range records {
		if records[i].ID == id {
			return records[i],true
		}
	}
	return Record{},false
}
func updateRecordPrice(records []Record, id int, newPrice float64) bool {
	for i := range records {
		if records[i].ID == id {
			records[i].Price = newPrice
			return true
		}
	}
	return false
}
func deleteRecord(records []Record, id int) ([]Record, bool) {
	for i := range records {
		if i+1 == id {
			records = append(records[:i], records[i+1:]...)
			return records, true
		}

	}
	return records, false

}
func countNumbers(nums []int) map[int]int{
	counts := make(map[int]int)
	for i:=0;i<len(nums);i++{
		num:=nums[i]
		counts[num]++
	}
	return counts
}
func main() {
	var records []Record
	record1 := Record{
		ID:     1,
		Title:  "Abbey Road",
		Artist: "The Beatles",
		Price:  169.0,
	}
	record2 := Record{
		ID:     2,
		Title:  "Purple Rain",
		Artist: "Prince",
		Price:  148.0,
	}

	record3 := Record{
		ID:     3,
		Title:  "Voodoo",
		Artist: "D'Angelo",
		Price:  150.0,
	}
	records = addRecord(records, record1)
	records = addRecord(records, record2)
	records = addRecord(records, record3)
	isupdated := updateRecordPrice(records, 99, 100)
	if !isupdated {
		fmt.Println("修改唱片价格失败")
	}
	for i := range records {
		fmt.Printf("%+v\n", records[i])
	}
	records, deleted := deleteRecord(records, 2)
	if !deleted {
		fmt.Println("删除唱片失败")

	} else {
		fmt.Println("删除唱片成功")
	}
	for i := range records {
		fmt.Printf("%+v\n", records[i])
	}

	foundrecord, found := findRecordByID(records, 99)
	if found {
		fmt.Printf("id:%d | title:%s | artist:%s | price:%.2f\n", foundrecord.ID, foundrecord.Title, foundrecord.Artist, foundrecord.Price)
	} else {
		fmt.Println("没找到该唱片")
	}
	nums := []int{1, 2, 1, 3, 2, 1}
	counts := make(map[int]int)
	counts=countNumbers(nums)
	for i:= range counts{
		fmt.Printf("%d:%d次\n",i,counts[i])
	}
	return
}
