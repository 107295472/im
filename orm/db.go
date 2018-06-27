package orm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"yin/xxd/models"
	"yin/xxd/util"
	"time"
	"github.com/json-iterator/go"
	"strconv"
	"strings"
)
var DbAccounts *gorm.DB
var IMDB *gorm.DB
func init()  {
	accout:="sqlserver://sa:sql2014@localhost:1433?database=RYAccountsDB&connection+timeout=30"
	imdb:="sqlserver://sa:sql2014@localhost:1433?database=IMDB&connection+timeout=30"
	initConn(1,accout)
	initConn(2,imdb)
}
func initConn(DBtype int,dbStr string)  {
	conn,err:=gorm.Open("mssql", dbStr)
	if err != nil {
		panic("failed to connect database")
	}
	conn.LogMode(true)
	conn.DB().SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	conn.DB().SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	conn.DB().SetConnMaxLifetime(time.Hour)
	switch DBtype {
	case 1:DbAccounts=conn
	case 2:IMDB=conn
	}
}
//获取聊天记录
func GetChatLog(userid string) ([]byte,error){
	//user:=Bing{Titlename:title,Imgurl:imgUrl}
	//sql:="SELECT userid,DynamicPass,DynamicPassTime,nickname,Service FROM dbo.AccountsInfo where userid="+strconv.FormatInt(userid,10)
	m:=[]models.IMRecord{}
	//err:=DbAccounts.Raw("SELECT userid,DynamicPass,DynamicPassTime,nickname,Service FROM dbo.AccountsInfo where userid=?", userid).Scan(&model).Error
	rows,err := IMDB.Model(&models.IMRecord{}).Where("sender = ? or recver=?",userid,userid).Select("id,content,sender,recver,time,sendername,recvername,avatar").Rows() // (*sql.Row)
	defer rows.Close()
	for rows.Next() {
		model:=models.IMRecord{}
		IMDB.ScanRows(rows,&model)
		//rows.Scan(&model.Id,&model.Content,&model.Sender,&model.Recver,&model.Time,&model.Sendername,&model.Recvername,&model.Avatar)
		//model.Username=model.Recvername
		m=append(m, model)
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonModel:=models.ChatLog{}
	jsonModel.Total=1
	jsonModel.Records=m
	jsonModel.Size=20
	 strByte,err:= json.Marshal(&jsonModel)
	//jsonStr:=convert(strByte)

	return strByte,err
}
func convert( b []byte ) string {
	s := make([]string,len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s,",")
}
//获取头像数据
func GetAvatar(userid string) ([]byte,error){
	//user:=Bing{Titlename:title,Imgurl:imgUrl}
	//sql:="SELECT userid,DynamicPass,DynamicPassTime,nickname,Service FROM dbo.AccountsInfo where userid="+strconv.FormatInt(userid,10)
	model:=[]byte{}
	//err:=DbAccounts.Raw("SELECT userid,DynamicPass,DynamicPassTime,nickname,Service FROM dbo.AccountsInfo where userid=?", userid).Scan(&model).Error
	row := DbAccounts.Table("AccountsFace").Where("UserID = ?", userid).Select("CustomFace").Row() // (*sql.Row)
	err:=row.Scan(&model)
	return model,err
}
//获取用户
func GetUser(userid int64) (models.User,error){
	//user:=Bing{Titlename:title,Imgurl:imgUrl}
	//sql:="SELECT userid,DynamicPass,DynamicPassTime,nickname,Service FROM dbo.AccountsInfo where userid="+strconv.FormatInt(userid,10)
	model:=models.User{}
	//err:=DbAccounts.Raw("SELECT userid,DynamicPass,DynamicPassTime,nickname,Service FROM dbo.AccountsInfo where userid=?", userid).Scan(&model).Error
	row := DbAccounts.Table("AccountsInfo").Where("UserID = ?", userid).Select("UserID, DynamicPass,DynamicPassTime,NickName,Service").Row() // (*sql.Row)
	err:=row.Scan(&model.UserID, &model.DynamicPass,&model.DynamicPassTime,&model.NickName,&model.Service)
	return model,err
}
//添加聊天记录
func AddRecord(user *models.IMRecord){
	//db:=IMDB.Exec("insert into im_records values(?,?,?,?,default,?,?,?)",user.Id,user.Msg,user.Sender,user.Recver,user.Sendername,user.Recvername,user.Offline)
	IMDB.NewRecord(user) // => 主键为空返回`true`
	IMDB.Create(&user)
	IMDB.NewRecord(user)
	if IMDB.Error!=nil{
		util.LogInfo().Println("AddRecord error",IMDB.Error)
	}
	//defer IMDB.Close()
	//modelArr:=make([]models.User,len(data))
	//model:=models.User{}
	//for index, obj := range data {
	//	model.UserID=string(obj["userid"])
	//	model.DynamicPass=string(obj["DynamicPass"])
	//	t,err:=time.ParseInLocation(util.Time_TIMEMSSQL, string(obj["DynamicPassTime"]), time.Local)//string(obj["DynamicPassTime"])
	//	util.Check(err)
	//	model.Service=string(obj["Service"])
	//	model.UserName=string(obj["nickname"])
	//	model.DynamicPassTime=t
	//	modelArr[index]=model
	//}
	//total := 0
	//for _, value := range a {
	//	total += value
	//}
	//result,err:=engine.Insert(&user)
}
//获取离线消息
func OfflineMessages(userid int64)([]models.OfflineMsg)  {
	//err:=DbAccounts.Raw("SELECT userid,DynamicPass,DynamicPassTime,nickname,Service FROM dbo.AccountsInfo where userid=?", userid).Scan(&model).Error
	rows,err:= IMDB.Table("IMRecord").Where("offline = 0 and recver=?",userid).Select("[content],[sender],[recver],[time],[sendername],[recvername]").Rows() // (*sql.Row)
	util.LogPrint("get im_records error:",err)
	m :=[]models.OfflineMsg{}
	defer rows.Close()
	for rows.Next() {
		model:=models.OfflineMsg{}
		rows.Scan(&model.Content,&model.Sender,&model.Recver,&model.Time,&model.Sendername,&model.Recvername)
		model.Method="offmsg"
		model.Module="chat"
		m=append(m, model)
	}
	IMDB.Exec("update IMRecord set offline=1 where recver=?",userid)
	//defer IMDB.Close()
	return m
}