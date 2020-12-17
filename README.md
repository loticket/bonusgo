# bonusgo
>  数字彩计算奖金，是否中奖，中奖金额,中奖级别


#### 双色球为例子计算奖金

```
package main

import (
	"github.com/loticket/bonusgo"
  
)

func main(){

  tickets := bonusgo.Ticket{
		Tid:      "2018",
		PlayType: 1,
		LotType:  2,
		LotNum:   "03,10,13,17,19,21,22,23,28-13,14",
		Money:    3024,
		BetNum:   1512,
		Multiple: 1,
	}
  
  //单协程计算
  var ss *bonusgo.TicketBonus = bonusgo.NewTicketBonus()
	err, t2 := ss.CalculatePrizeRun("ssq", "03,10,13,22,23,28-13", []int{809250800, 21109100}, tickets)
	fmt.Println(err)
	fmt.Println(t2)
  
  //多协程计算

  ss.StartCalculatePrize("ssq", "03,10,13,22,23,28-15", []int{678818300, 543054600, 24662700, 19730100}, func(res bonusgo.TicketCalResult) {
		fmt.Println(res)
	})
	for i := 0; i < 100; i++ {
		tickets.Tid = strconv.Itoa(i)
		ss.SetTicketChan(tickets)
	}
  for{}

}

```
