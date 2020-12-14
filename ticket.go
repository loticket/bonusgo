package bonusgo

//数字彩投注内容
type Ticket struct {
	PlayType int    `json:"playtype" binding:"required"` //子玩法
	LotType  int    `json:"lottype" binding:"required"`  //选号方式
	LotNum   string `json:"lotnum" binding:"required"`   //投注号码
	Money    int    `json:"money" binding:"required"`    //投注金额
	BetNum   int    `json:"betnum" binding:"required"`   //注数
	Multiple int    `json:"multiple" binding:"required"` //倍数
}

//开奖数据
type OpenBonus struct {
	LotId    string           `json:"lotid" binding:"lotid"`       //彩种编码
	OpenNum  string           `json:"opennum" binding:"opennum"`   //开奖号码
	OpenInfo []OpenBonusLevel `json:"openinfo" binding:"openinfo"` //开奖号码
}

type OpenBonusLevel struct {
	Money int `json:"money" binding:"money"` //中奖金额
	Types int `json:"types" binding:"types"` //中奖类型  1 普通   2 其他
	Level int `json:"level" binding:"level"` //中奖级别
}
