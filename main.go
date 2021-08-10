package main

import (
	"fmt"
	pushover2 "github.com/gregdel/pushover"
	"github.com/thedevsaddam/gojsonq"
	"io"
	"log"
	"net/http"
	"time"
)

type StoreShanghai struct {
	Id string
	Name string
}

func getIkeaStock(storeId string, goodId string) float64 {
	ikeaUlr := "https://www.ikea.cn/api-host/store/%s/item/%s/stockAndLocation?itemType=ART"
	ikeaUlr = fmt.Sprintf(ikeaUlr, storeId, goodId)
	body, err := http.Get(ikeaUlr)
	if err != nil {
		fmt.Errorf("get ikea error")
		panic(err)
	}
	defer body.Body.Close()
	content, err := io.ReadAll(body.Body)
	if err != nil{
		panic(err)
	}
	//fmt.Printf("%s\n", content)
	// resp to json
	data := gojsonq.New().FromString(string(content))
	stockNum := data.Find("stock")
	//fmt.Printf("%+v", stockNum)
	//if stockNum.(float64) > 0 {
	//	fmt.Println("In Stock!")
	//
	//}else{
	//	fmt.Println("No stock.")
	//}
	return stockNum.(float64)
}

func pushoverSender(msg string) {
	app := pushover2.New("app")
	recipient := pushover2.NewRecipient("utm")
	message := pushover2.NewMessageWithTitle(msg, "宜家到货提醒")

	response, err := app.SendMessage(message, recipient)
	if err != nil {
		log.Panic(err)
	}
	log.Println(response)
}

func main(){
	storeList := []StoreShanghai{
		{"856", "徐汇"},
		{"247", "宝山"},
		{"585", "杨浦"},
		{"885", "北蔡"},
	}
	goodId := "10476479"
	fmt.Printf("货物ID: %s\n", goodId)
	pushoverSender(fmt.Sprintf("本次开始监控宜家货品ID: %s", goodId))
	for {
		for _, store := range storeList {
			stockResult := getIkeaStock(store.Id, goodId)
			if stockResult > 0 {
				fmt.Printf("%s%s\n", store.Name, "出现库存")
				msg := fmt.Sprintf("货物%s,%s%s\n", goodId, store.Name, "出现库存")
				pushoverSender(msg)
				break
			} else {
				fmt.Printf("%s%s\n", store.Name, "没有库存")
			}
		}
		fmt.Println("等待1小时后继续刷新")
		time.Sleep(time.Duration(1) * time.Hour)
	}

}
