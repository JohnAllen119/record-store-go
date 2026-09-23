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
func addRecord(records []Record,record Record)[]Record{
	records=append(records,record)
	
	return records
}
func findRecordByID(records []Record,id int)(Record,bool){
	for _,record :=range records{
		if record.ID==id{
			return record,true
		}
	}
	return Record{},false
}
func filterRecords(records []Record,artist string,minPrice float64)[]Record{
	var records2 []Record
	for _, record := range records{
		if record.Artist==artist &&record.Price>=minPrice {
			records2=append(records2,record)
		}
	}
	return records2
}
func main() {
	var records []Record
	record1 := Record{
		ID:     1,
		Title:  "Revolver",
		Artist: "The Beatles",
		Price:  148.0,
	}

	record2 := Record{
		ID:     2,
		Title:  "Purple Rain",
		Artist: "Prince",
		Price:  138.0,
	}
	record3 := Record{
		ID:     3,
		Title:  "Voodoo",
		Artist: "D'Angelo",
		Price:  168.0,
	}
	records=addRecord(records,record1)
	records=addRecord(records,record2)
	records=addRecord(records,record3)

	for _, record := range records {
		fmt.Printf("%+v\n", record)
	}
	fmt.Println("唱片数量等于:", len(records))

	record,found := findRecordByID(records,99)
	if found{
		fmt.Printf("%+v\n",record)
	}else{
		fmt.Println("没有找到唱片")
	}
	records2:=filterRecords(records, "The Beatles", 1000.0)
	if len(records2) == 0 {
	fmt.Println("没有符合条件的唱片")
}
	for _, filteredrecord := range records2{
		fmt.Printf("%+v\n",filteredrecord)
	}
	fmt.Println("筛选结果数量:", len(records2))
	return
}
