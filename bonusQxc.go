package bonusgo

import (
	"errors"
	"strconv"
	"strings"
)

type BonusQxc struct {
	BaseBonus
}

/**
 * @name:计算单注奖金
 * @msg:计算单注奖金
 * @param playType 投注方式 1:普通  2 追加 lotNum 投注号码  bonusRedBall 开奖红球字符串 bonusBlueBall 篮球字符串 bonusLevel 格式化完成的奖级与开奖金额映射
 * @return: int 中奖金额 int 中奖级别
 */
func (bonus *BonusQxc) Calculation(playType int, lotNum string, bonusRedBall string, bonusBlueBall string, bonusLevel map[string][]int) (int, int) {
	var allBall []string = strings.Split(lotNum, "-")
	var redBall []string = strings.Split(allBall[0], ";")

	var bonusRedBallArr []string = strings.Split(bonusRedBall, ";")

	var redNum int = 0
	var blueNum int = 0

	for k, ball := range bonusRedBallArr {
		if strings.EqualFold(redBall[k], ball) {
			redNum++
		}
	}

	if strings.EqualFold(allBall[1], bonusBlueBall) {
		blueNum++
	}

	if redNum == 0 && blueNum == 0 {
		return 0, 0
	}

	//直接计算出来奖金
	var keys string = strconv.Itoa(redNum) + "_" + strconv.Itoa(blueNum) + "_" + strconv.Itoa(playType)

	levelCal, ok := bonusLevel[keys]

	if !ok {
		return 0, 0
	}

	return levelCal[0], levelCal[1]

}

/**
 * @name:简单验证开奖号码并且格式化开奖号码
 * @msg:在计算开奖之前一定要验证开奖信息
 * @param nil
 * @return: []string{红球字符串,篮球字符串} error 错误提示
 */
func (bonus *BonusQxc) FormateOpenNum() ([]string, error) {
	var OpenBonusArr []string = make([]string, 0)
	if strings.Count(bonus.openNum, "-") != 1 {
		return OpenBonusArr, errors.New("开奖号码错误")
	}

	if strings.Count(bonus.openNum, ";") != 5 {
		return OpenBonusArr, errors.New("开奖号码错误")
	}

	return strings.Split(bonus.openNum, "-"), nil
}

/**
 * @name:单式，复式，胆拖拆票
 * @msg:把全部的票拆成单式，按照单式计算奖金
 * @param nil
 * @return: []string{单式号码字符串,单式号码字符串}
 */
func (bonus *BonusQxc) SplitTicket(ticket Ticket) []string {
	switch ticket.LotType {
	case 1:
		return []string{ticket.LotNum}
	case 2:
		return bonus.spliteDuplicate(ticket)
	default:
		return []string{}
	}
	return []string{}
}

/**
 * @name:复式拆票
 * @msg:把全部的复式票拆成单式，按照单式计算奖金
 * @param nil
 * @return: []string{单式号码字符串,单式号码字符串}
 */
func (bonus *BonusQxc) spliteDuplicate(ticket Ticket) []string {
	var allBall []string = strings.Split(ticket.LotNum, "-")
	var redBall []string = strings.Split(allBall[0], ";")
	var blueBall []string = strings.Split(allBall[1], ",")

	var oneNum []string = strings.Split(redBall[0], ",")
	var twoNum []string = strings.Split(redBall[1], ",")
	var threeNum []string = strings.Split(redBall[2], ",")
	var fourNum []string = strings.Split(redBall[3], ",")
	var fiveNum []string = strings.Split(redBall[4], ",")
	var sixNum []string = strings.Split(redBall[5], ",")

	var allBallArr [][]string = [][]string{twoNum, threeNum, fourNum, fiveNum, sixNum}

	var i int = 0
	var ticketOne []string = oneNum
	for i < 5 {
		var tempTicketOne []string = make([]string, 0)
		for j := 0; j < len(ticketOne); j++ {
			for k := 0; k < len(allBallArr[i]); k++ {
				tempTicketOne = append(tempTicketOne, ticketOne[j]+";"+allBallArr[i][k])
			}
		}
		ticketOne = tempTicketOne
		i++
	}
	var newTicket []string = make([]string, 0)
	for _, blue := range blueBall {
		for _, red := range ticketOne {
			newTicket = append(newTicket, red+"-"+blue)
		}
	}

	return newTicket

}

/**
 * @name:实例化大乐透的算奖结构体
 * @msg:把全部的胆拖票拆成单式，按照单式计算奖金
 * @param nil
 * @return: LotteryBonusInterface 接口(BonusDlt实现了此接口)
 */
func NewBonusQxc() LotteryBonusInterface {
	return &BonusQxc{
		BaseBonus{
			bonus: []Bonus{
				Bonus{Grade: 1, Redball: 6, BlueBall: 1, Monye: 0, Types: 1},
				Bonus{Grade: 2, Redball: 6, BlueBall: 0, Monye: 0, Types: 1},
				Bonus{Grade: 3, Redball: 5, BlueBall: 1, Monye: 300000, Types: 1},
				Bonus{Grade: 4, Redball: 5, BlueBall: 0, Monye: 50000, Types: 1},
				Bonus{Grade: 4, Redball: 4, BlueBall: 1, Monye: 50000, Types: 1},
				Bonus{Grade: 5, Redball: 4, BlueBall: 0, Monye: 3000, Types: 1},
				Bonus{Grade: 5, Redball: 3, BlueBall: 1, Monye: 3000, Types: 1},
				Bonus{Grade: 6, Redball: 3, BlueBall: 0, Monye: 500, Types: 1},
				Bonus{Grade: 6, Redball: 2, BlueBall: 1, Monye: 500, Types: 1},
				Bonus{Grade: 6, Redball: 1, BlueBall: 1, Monye: 500, Types: 1},
				Bonus{Grade: 6, Redball: 0, BlueBall: 1, Monye: 500, Types: 1},
			},
			redBall:   NumBall{Min: 1, Max: 9},
			blueBall:  NumBall{Min: 1, Max: 14},
			maxLervel: 6,
		},
	}
}

/**
 * @name:实例化大乐透算奖结构体
 * @msg:实例化大乐透算奖结构体
 * @param nil
 * @return: func() LotteryBonusInterface 接口(BonusDlt实现了此接口)
 */
func init() {
	RegisterBonusMethod("qxc", NewBonusQxc)
}
