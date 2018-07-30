package eth_query

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type ethtxrecordst struct {
	Id          int64
	From        string
	To          string
	Amount      string
	Nonce       string
	Fee         string
	TxHash      string
	BlockHeight string
	ConfirmTime string
	Comment     string
	Created     time.Time
	BlockState  int
	TxState     int
	IsToken     int
	Collection  int
}

type ethtokentxrecordst struct {
	ethtxrecordst
	LogIndex string
	Decimal  string
}

// eth-query 相关函数
func GetEthTxrecord(addr string, page uint64, size uint64) (txs *[]ethtxrecordst, r int64) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("t1.id", "t1.from", "t1.to", "t1.amount", "t1.fee", "t1.nonce", "t1.tx_hash",
		"t1.block_height", "t1.confirm_time", "t1.created", "t1.tx_state", "t1.is_token", "t2.comment",
		"0 as collection").From("tx as t1").
		LeftJoin("tx_extra_info as t2").
		On("t1.tx_hash = t2.tx_hash").
		Where("t1.from = ?")
	sql1 := qb.String() + " union all "
	qb.Select("t1.id", "t1.from", "t1.to", "t1.amount", "t1.fee", "t1.nonce", "t1.tx_hash",
		"t1.block_height", "t1.confirm_time", "t1.created", "t1.tx_state", "t1.is_token", "t2.comment",
		"1 as collection").From("tx as t1").
		LeftJoin("tx_extra_info as t2").
		On("t1.tx_hash = t2.tx_hash").
		Where("t1.to = ?").
		OrderBy("created").Desc().
		OrderBy("id").Desc().Limit(int(size)).Offset(int(page))
	sql2 := qb.String()
	sql := sql1 + sql2
	o := orm.NewOrm()
	o.Using(databases)
	r, err := o.Raw(sql, addr, addr).QueryRows(txs)
	if err != nil {
		log.Errorf("GetEthTxrecord error %s", err)
		return
	}
	return
}

func GetEthTokenTxrecord(addr string, contract string, page uint64, size uint64) (txs *[]ethtokentxrecordst, r int64) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("t1.id", "t1.tx_hash", "t1.from", "t1.to", "t1.amount", "t1.decimal", "t1.fee",
		"t1.block_height", "t1.confirm_time", "t1.created", "t1.log_index", "t1.tx_state", "t1.is_token",
		"t2.comment", "0 as collection").
		From("token_tx as t1").
		LeftJoin("tx_extra_info as t2").
		On("t1.tx_hash = t2.tx_hash").
		Where("t1.from = ?").And("t1.contract_addr = ?")
	sql1 := qb.String() + " union all "
	qb.Select("t1.id", "t1.tx_hash", "t1.from", "t1.to", "t1.amount", "t1.decimal", "t1.fee",
		"t1.block_height", "t1.confirm_time", "t1.created", "t1.log_index", "t1.tx_state", "t1.is_token",
		"t2.comment", "1 as collection").
		From("token_tx as t1").
		LeftJoin("tx_extra_info as t2").
		On("t1.tx_hash = t2.tx_hash").
		Where("t1.to = ?").And("t1.contract_addr = ?").
		OrderBy("created").Desc().
		OrderBy("id").Desc().Limit(int(size)).Offset(int(page))

	sql2 := qb.String()
	sql := sql1 + sql2
	o := orm.NewOrm()
	o.Using(databases)
	r, err := o.Raw(sql, addr, contract, addr, contract).QueryRows(txs)
	if err != nil {
		log.Errorf("GetEthTokenTxrecord error %s", err)
		return
	}
	return
}

// 判断eth用户是否存在
func GetEthAddrExist(addr string) bool {
	o := orm.NewOrm()
	o.Using(databases)
	qs := o.QueryTable("address")
	return qs.Filter("addr", addr).Exist()
}
