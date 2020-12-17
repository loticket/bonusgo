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
func (tb *TicketBonus) CalculatePrize(ticket Ticket, bonusOpen []string, bonusLevel map[string][]int) TicketCalResult {
	tbm, tam, big, le, tid := tb.lotteryBonus.CalculatePrize(ticket, bonusOpen, bonusLevel)
	return TicketCalResult{
		TaxBeforeMoneyAll: tbm,
		TaxAfterMoneyAll:  tam,
		Big:               big,
		Levels:            le,
		Tid:               tid,
	}
}

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
