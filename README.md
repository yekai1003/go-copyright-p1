# go-copyright-p1

## 1 目录介绍

-  static  html页面目录
-  configs 配置文件读取处理
-  routes  路由处理
-  etc   配置文件处理
-  utils 通用处理，错误信息
-  dbs   数据库处理文件
-  eths  以太坊相关处理

## 2 环境安装需要

### 2.1 配置文件读取插件需要安装

```
go get -u github.com/BurntSushi/toml
```

### 2.2 echo框架安装

先安装crypto,labstack使用了该库，需要借助github，go get不能用
```
cd $GOPATH/src
mkdir -p golang.org/x/
cd golang.org/x/
git clone https://github.com/golang/crypto.git
```
安装echo
```
go get -u github.com/labstack/echo
go get -u github.com/labstack/echo-contrib/session
```



### 2.3 开发过程可能需要用到的库安装方法如下


mysql的go语言驱动安装
```
go get -u github.com/go-sql-driver/mysql
```

其他可能涉及的库
```
go get -u github.com/labstack/gommon/
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/go-sql-driver/mysql
```

### 2.4 echo框架学习资料

[echo框架学习](https://echo.labstack.com/guide)

## 3 数据库建库脚本

[建库语句](etc/copyright.sql)

## 4 接口说明

### 4.1 注册

名称 | 说明
---|---
url | /account
method | POST
请求数据 | json:email,identity_id,username
响应数据 | json:errno,errmsg


响应举例
```
{"errno":"0","errmsg":"成功","data":null}
```


### 4.2 会话（注册后会有此请求）

名称 | 说明
---|---
url | /session
method | GET
请求数据 | 无
响应数据 | json:errno,errmsg

响应举例
```
{"errno":"0","errmsg":"成功","data":null}
```

### 4.3 登陆

名称 | 说明
---|---
url | /account
method | POST
请求数据 | json:identity_id,username
响应数据 | json:errno,errmsg


响应举例
```
{"errno":"0","errmsg":"成功","data":null}
```

### 4.4 上传图片

名称 | 说明
---|---
url | /content
method | POST
请求数据 | form中文件二进制数据
响应数据 | json:errno,errmsg

请求举例

```
curl --form "fileupload=@wyz.jpg" http://localhost:8086/content
```


### 4.5 查看用户所有图片

名称 | 说明
---|---
url | /content
method | GET
请求数据 | 无
响应数据 | json:{Errno,Errmsg,Data:{total_page,current_page,contents:[]}}

说明：contents 是一个结构体数组，举例如下：
```
{
	"errno": "0",
	"errmsg": "成功",
	"data": {
		"contents": [{
			"content_hash": "2b48f7f9b1156bd26d13d498d2a2441ca857922d743e5ab7d8fde20dd18e3f21",
			"title": "马夏尔2.jpeg"
		}, {
			"content_hash": "f72d0385722cb58fddaa493557950fe83fe317a0aa41d44ddb386c82f902d9dd",
			"title": "鸟叔2.jpg"
		}, {
			"content_hash": "6bb7e5ec0c5f057cfdb25e54a2f0f09bca5819213761697fa57e6368eed1ec9d",
			"title": "am9.jpg"
		}],
		"current_page": 1,
		"total_page": 1
	}
}
```

### 4.6 查看单个图片

名称 | 说明
---|---
url | /content/==picturename==
method | GET
请求数据 | 无
响应数据 | 图片二进制数据

### 4.7 开始拍卖

用户自行发起拍卖，指定份数，指定价格

名称 | 说明
---|---
url | /aution
method | POST
请求数据 | json:{content_hash,percent,price}
响应数据 | 图片二进制数据

请求数据举例：
```
content_hash: "6bb7e5ec0c5f057cfdb25e54a2f0f09bca5819213761697fa57e6368eed1ec9d",
percent: 30,
price: 20
```

### 4.8 查看当前拍卖

名称 | 说明
---|---
url | /auctions
method | GET
请求数据 | 无
响应数据 | json:{errno,errmsg,data;[{price,title,token_id}]}

响应数据举例

```
{"errno":"0","errmsg":"成功","data":[{"price":"100","title":"马夏尔和女友.jpeg","token_id":"6"}]}
```

### 4.9 拍卖图片

名称 | 说明
---|---
url | /auctions
method | GET
请求数据 | 无
响应数据 | json:{errno,errmsg,data;[{price,title,token_id}]}


## 5 功能函数使用


### 5.1 解析订阅事件内容相关
```
func LogDataUnpack(start, end int, val interface{}, data []byte) (err error) {
	length := len(data)
	fmt.Println("call--- LogDataUnpack begin", reflect.TypeOf(val).String(), length)

	if start >= length || end > length {
		return errors.New("too short datas")
	}
	pdata := data[start:end]

	fmt.Println(string(data), string(pdata))
	if reflect.TypeOf(val).String() == "int64" || reflect.TypeOf(val).String() == "*int64" {
		var tmpval *int64 = val.(*int64)
		*tmpval, err = strconv.ParseInt(string(pdata), 16, 32)
		fmt.Println("call ParseInt", val)
	} else if reflect.TypeOf(val).String() == "string" || reflect.TypeOf(val).String() == "*string" {
		var tmpval *string = val.(*string)
		*tmpval = string(pdata)
		fmt.Println("call ParseInt", val)
	}

	fmt.Println("call--- LogDataUnpack end", val)
	return nil
}
func ParseMintEvent2Db(data []byte) error {
	fmt.Println(string(data))
	var tokenId int64
	err := LogDataUnpack(32*5, 32*6, &tokenId, data)
	if err != nil {
		fmt.Println("faile to get tokenid", err)
		return err
	}
	fmt.Println("tokenid===", tokenId)
	var pixHash string
	err = LogDataUnpack(32*0, 32*2, &pixHash, data)
	if err != nil {
		return err
	}
	fmt.Println("pixHash===", pixHash)
	var pixAddr string
	err = LogDataUnpack(88, 128, &pixAddr, data)
	if err != nil {
		return err
	}
	pixAddr = "0x" + pixAddr
	fmt.Println("pixAddr===", pixAddr)
	//插入到数据库中
	sql := fmt.Sprintf("insert into account_content(content_hash,token_id,address) values('%s',%d,'%s')", pixHash, tokenId, pixAddr)
	_, err = dbs.Create(sql)
	if err != nil {
		fmt.Println("failed to insert into mysql ", sql, err)
		return err
	}
	return nil
}
```

