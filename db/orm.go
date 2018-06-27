package db
//
//import (
//	_"github.com/denisenkom/go-mssqldb"
//	"github.com/go-xorm/xorm"
//	"yin/xxd/util"
//	"time"
//	"yin/xxd/models"
//	"fmt"
//	"strconv"
//)
////var DbAccounts *xorm.Engine
////var IMDB *xorm.Engine
////func init()  {
////accout:="odbc:server=localhost;user id=sa;password=sql2014;database=RYAccountsDB;connection timeout=30"
////imdb:="odbc:server=localhost;user id=sa;password=sql2014;database=IMDB;connection timeout=30"
////initConn(1,accout)
////initConn(2,imdb)
////}
////func initConn(DBtype int,dbStr string)  {
////	conn,err:=xorm.NewEngine("mssql", dbStr)
////	if err != nil {
////		panic("failed to connect database")
////	}
////	conn.ShowSQL(true)
////	err=conn.Ping()
////	if err!=nil {
////		fmt.Println(err)
////	}
////	conn.TZLocation=time.Local
////	conn.SetMaxIdleConns(5)
////	conn.SetMaxOpenConns(30)
////	switch DBtype {
////	case 1:DbAccounts=conn
////	case 2:IMDB=conn
////	}
////}
////func GetUser(userid int64) (models.User,error){
////	//user:=Bing{Titlename:title,Imgurl:imgUrl}
////	sql:="SELECT userid,DynamicPass,DynamicPassTime,nickname,Service FROM dbo.AccountsInfo where userid="+strconv.FormatInt(userid,10)
////	DbAccounts.Query(sql)
////	data,err:=DbAccounts.Query()//Where("userid="+string(userid)).Query()
////	if err!=nil{
////		util.LogInfo().Println("查询错误")
////	}
////	modelArr:=make([]models.User,len(data))
////	model:=models.User{}
////	for index, obj := range data {
////	model.UserID=string(obj["userid"])
////	model.DynamicPass=string(obj["DynamicPass"])
////	t,err:=time.ParseInLocation(util.Time_TIMEMSSQL, string(obj["DynamicPassTime"]), time.Local)//string(obj["DynamicPassTime"])
////	util.Check(err)
////	if string(obj["Service"])==""{
////		model.Service="0"
////	}else {
////		model.Service=string(obj["Service"])
////	}
////	model.UserName=string(obj["nickname"])
////	model.DynamicPassTime=t
////	modelArr[index]=model
////	}
////	//total := 0
////	//for _, value := range a {
////	//	total += value
////	//}
////	//result,err:=engine.Insert(&user)
////	if len(modelArr)>0{
////		return modelArr[0],err
////	}
////	return model,err
////}
////func AddRecord(user *models.Record){
////	IMDB.ShowSQL(true)
////	//user:=Bing{Titlename:title,Imgurl:imgUrl}
////	_,err:=IMDB.Exec("insert into values(?,?,?,?,?,default,?,?,?)",user.Id,user.Msg,user.Username,user.Sender,user.Recver,user.Sendername,user.Recvername,user.Offline)
////	if err!=nil{
////		util.LogInfo().Println("添加记录错误")
////	}
/////*	modelArr:=make([]models.User,len(data))
////	model:=models.User{}
////	for index, obj := range data {
////		model.UserID=string(obj["userid"])
////		model.DynamicPass=string(obj["DynamicPass"])
////		t,err:=time.ParseInLocation(util.Time_TIMEMSSQL, string(obj["DynamicPassTime"]), time.Local)//string(obj["DynamicPassTime"])
////		util.Check(err)
////		model.Service=string(obj["Service"])
////		model.UserName=string(obj["nickname"])
////		model.DynamicPassTime=t
////		modelArr[index]=model
////	}*/
////	//total := 0
////	//for _, value := range a {
////	//	total += value
////	//}
//	//result,err:=engine.Insert(&user)
//}