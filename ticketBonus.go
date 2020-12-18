package bonusgo

import (
	"fmt"
	"sync"
)

func NewTicketBonus() *TicketBonus {
	return &TicketBonus{
		TicketCal:    make(chan Ticket, 1024),
		TicketCalRes: make(chan TicketCalResult, 1024),
		Stop:         make(chan bool, 0),
		Busy:         false,
	}
}

type TicketBonus struct {
	lotteryBonus LotteryBonusInterface
	TicketCal    chan Ticket
	TicketCalRes chan TicketCalResult
	Stop         chan bool
	Busy         bool
}

/**
 * @name:多协程算奖提高算奖效率
 * @msg:根据传进去的字符串，把相应的算奖结构体付值给 lotteryBonus
 * @param alg string openNum 开奖号码 level 开奖等级金额  callback 回掉函数
 * @return: error 错误提示
 */
func (tb *TicketBonus) StartCalculatePrize(alg string, openNum string, level []int, callback func(TicketCalResult)) error {
	if err := tb.CreateBonus(alg); err != nil {
		return err
	}

	tb.lotteryBonus.SetOpenBonus(openNum)
	tb.lotteryBonus.SetBonusMoney(level)

	formateBonus, errf := tb.lotteryBonus.FormateOpenNum()
	if errf != nil {
		return errf
	}

	var levelFormate map[string][]int = tb.lotteryBonus.BonusLevelFormate()
	go func() {
		for {
			select {
			case ticket := <-tb.TicketCal:
				go func() {
					res := tb.CalculatePrize(ticket, formateBonus, levelFormate)
					tb.TicketCalRes <- res
					callback(res)
				}()
			case <-tb.Stop:
				goto Loop
			}
			tb.wath()
			fmt.Println("-------检察中呢-----")
		}

	Loop:
	}()
	return nil

}

/**
 * @name:观察算奖是否已经停止
 * @msg:观察算奖状态，busy 为 true 正在算奖  false 算奖结束
 * @param nil
 * @return: nil
 */
func (tb *TicketBonus) wath() {
	if (len(tb.TicketCal) > 0 && tb.Busy) || (len(tb.TicketCal) == 0 && !tb.Busy) {
		return
	}

	if len(tb.TicketCal) > 0 && !tb.Busy {
		tb.Busy = true
	} else if len(tb.TicketCal) == 0 && tb.Busy {
		tb.Busy = false
	}
}

/**
 * @name:给指定的票算奖
 * @msg:给指定的票算奖，把票信息放到chan中
 * @param ticket 票信息
 * @return: nil
 */
func (tb *TicketBonus) SetTicketChan(ticket Ticket) {
	tb.Busy = true
	tb.TicketCal <- ticket
}

/**
 * @name:停止算奖
 * @msg:停止算奖
 * @param nil
 * @return: nil
 */

func (tb *TicketBonus) StopCalculatePrize() {
	if tb.Busy {
		return
	}
	tb.Stop <- true
}

/**
 * @name:单协程算奖
 * @msg:根据传进去的字符串，把相应的算奖结构体付值给 lotteryBonus
 * @param alg string openNum 开奖号码 level 开奖等级金额  ticket 需要算奖的票
 * @return: error 错误提示 TicketCalResult算奖统一结果
 */
func (tb *TicketBonus) CalculatePrizeRun(alg string, openNum string, level []int, ticket Ticket) (error, TicketCalResult) {
	if err := tb.CreateBonus(alg); err != nil {
		return err, TicketCalResult{}
	}

	tb.lotteryBonus.SetOpenBonus(openNum)
	tb.lotteryBonus.SetBonusMoney(level)

	formateBonus, errf := tb.lotteryBonus.FormateOpenNum()
	if errf != nil {
		return errf, TicketCalResult{}
	}

	var levelFormate map[string][]int = tb.lotteryBonus.BonusLevelFormate()
	ticketCalResult := tb.CalculatePrize(ticket, formateBonus, levelFormate)
	return nil, ticketCalResult
}

/**
 * @name:获取对应的彩种算奖对象
 * @msg:根据传进去的字符串，把相应的算奖结构体付值给 lotteryBonus
 * @param alg string
 * @return: error 错误提示
 */
func (tb *TicketBonus) CreateBonus(alg string) error {
	var err error = nil
	if methodF, ok := bonusMethods[alg]; ok {
		tb.lotteryBonus = methodF()
	} else {
		err = fmt.Errorf("bonus: unknown struct name %q (forgot to import?)", alg)
	}
	return err
}

/**
 * @name:计算一张票的奖金
 * @msg:计算单注奖金，可以设定完成开奖号码，和开奖级别处理以后，通过设置不同的ticket计算同期号的所有奖金
 * @param bonusOpen 开奖号码数组 []string{红球字符串，篮球字符串}(可调用 FormateOpenNum()函数获取)
 * @param bonusLevel 开奖信息  map[string][]int map[红球个数_篮球个数_开奖方式][]int{奖级，奖金}(可调用 BonusLevelFormate()函数)
 * @return: int64 税前金额  int64 税后金额  int 是否是大额奖金  []int 中奖注数 string 票的唯一标示
 */
func (tb *TicketBonus) CalculatePrize(ticket Ticket, bonusOpen []string, bonusLevel map[string][]int) TicketCalResult {
	//格式化中奖
	var taxBeforeMoneyAll int64 = 0 //税前金额
	var taxAfterMoneyAll int64 = 0  //税后金额
	var maxlevels int = tb.lotteryBonus.GetmaxLervel()
	var levels []int = make([]int, maxlevels) //中奖注数
	var big int = 0
	//获取拆票信息
	var ticketAll []string = tb.lotteryBonus.SplitTicket(ticket)

	for i := 0; i < len(ticketAll); i++ {
		grad, money := tb.lotteryBonus.Calculation(ticket.PlayType, ticketAll[i], bonusOpen[0], bonusOpen[1], bonusLevel)
		if grad == 0 {
			continue
		}
		taxBeforeMoneyAll += int64(money * ticket.Multiple)
		big = 2
		if grad < 3 && money > 1000000 {
			taxAfterMoneyAll += int64((money * 8 / 10) * ticket.Multiple)
			big = 1
		} else {
			taxAfterMoneyAll += int64(money * ticket.Multiple)
		}

		levels[grad-1]++
	}

	return TicketCalResult{
		TaxBeforeMoneyAll: taxBeforeMoneyAll,
		TaxAfterMoneyAll:  taxAfterMoneyAll,
		Big:               big,
		Levels:            levels,
		Tid:               ticket.Tid,
	}

}

/**
 * @name:付值彩种信息
 * @msg:付值彩种信息
 * @param ticket Ticket
 * @return: nil
 */
func (tb *TicketBonus) SetBonusTicketInfo(playtype int, lottype int, lotnum string, money int, betnum int, multiple int) Ticket {
	return Ticket{
		PlayType: playtype,
		LotType:  lottype,
		LotNum:   lotnum,
		Money:    money,
		BetNum:   betnum,
		Multiple: multiple,
	}

}

/**
 * @name:付值开奖信息
 * @msg:付值彩种信息
 * @param openNum 开奖号码  level 开奖等级奖金
 * @return: nil
 */
func (tb *TicketBonus) SetBonusInfo(openNum string, level []int) {
	tb.lotteryBonus.SetBonusMoney(level)
	tb.lotteryBonus.SetOpenBonus(openNum)
}

/**
 * @name:付值开奖信息
 * @msg:付值彩种信息
 * @param openNum 开奖号码  level 开奖等级奖金
 * @return: nil
 */
func (tb *TicketBonus) BonusOpen() ([]string, error) {
	return tb.lotteryBonus.FormateOpenNum()
}

/**
 * @name:付值开奖信息
 * @msg:付值彩种信息
 * @param openNum 开奖号码  level 开奖等级奖金
 * @return: nil
 */
/*func (tb *TicketBonus) CalculatePrize(ticket Ticket, bonusOpen []string, bonusLevel map[string][]int) TicketCalResult {
	tbm, tam, big, le, tid := tb.lotteryBonus.CalculatePrize(ticket, bonusOpen, bonusLevel)
	return TicketCalResult{
		TaxBeforeMoneyAll: tbm,
		TaxAfterMoneyAll:  tam,
		Big:               big,
		Levels:            le,
		Tid:               tid,
	}
}*/

/**
 * @name:付值彩种信息
 * @msg:付值彩种信息
 * @param ticket Ticket
 * @return: nil
 */
func (tb *TicketBonus) SetBonusTicket(ticket Ticket) {
	tb.lotteryBonus.SetTicket(ticket)
}

var bonusMethods = map[string]func() LotteryBonusInterface{}
var signingMethodLock = new(sync.RWMutex)

func RegisterBonusMethod(alg string, f func() LotteryBonusInterface) {
	signingMethodLock.Lock()
	defer signingMethodLock.Unlock()
	if _, ok := bonusMethods[alg]; ok {
		return
	}

	bonusMethods[alg] = f
}

func GetBonusMethod(alg string) (method LotteryBonusInterface, err error) {
	signingMethodLock.RLock()
	defer signingMethodLock.RUnlock()

	if methodF, ok := bonusMethods[alg]; ok {
		method = methodF()
		err = nil
	} else {
		err = fmt.Errorf("bonus: unknown struct name %q (forgot to import?)", alg)
	}
	return
}
