package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var gclient http.Client
var cookies []*http.Cookie

func downloadFile(savefile,url string) {
	res, _ := http.Get(url)
	file, _ := os.Create(savefile)
	io.Copy(file, res.Body)
}


func LoginSetCookie(uname,pwd string) []*http.Cookie {

	encodepwd := base64.StdEncoding.EncodeToString([]byte(pwd))

	song := make(map[string]interface{})
	song["account"] = uname
	song["pwd"] = encodepwd
	bytesData, err := json.Marshal(song)
	if err != nil {
		fmt.Println(err.Error() )
		return nil
	}

	var resp *http.Response
	req, err := http.NewRequest("POST", "http://m.lrts.me/ajax/logon",bytes.NewReader(bytesData))
	if err != nil {
		return nil
	}


	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err1 := gclient.Do(req)
	if err1 != nil {
		fmt.Println("Fatal error ", err.Error())
		return nil
	}
	defer resp.Body.Close()

	//for _, v := range resp.Cookies() {
	//	fmt.Printf("%+v\n", v)
	//	req.AddCookie(v)
	//}

	return resp.Cookies()
}

func GetBookList(bookid int) map[string]interface{} {


	booklisturl := fmt.Sprintf("http://m.lrts.me/ajax/getBookMenu?bookId=%d&pageNum=1&pageSize=5000&sortType=0",bookid)
	reqest, err := http.NewRequest("GET", booklisturl, nil) //建立一个请求
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return nil
	}

	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.95 Safari/537.36")

	response, err := gclient.Do(reqest)
	body, err1 := ioutil.ReadAll(response.Body)
	if err1 != nil {
		fmt.Println("Fatal error ", err.Error())
		return nil
	}

	chapters := make(map[string]interface{})
	err = json.Unmarshal(body, &chapters)

	return chapters

}


// index  从0开始
func GetChapterName(bookinfo map[string]interface{},index int) string{
	chapters, ok:= bookinfo["list"].([]interface{})
	if ok == false {
		fmt.Println("Fatal error ")
		return ""
	}

	onech, ok := chapters[index].(map[string]interface{})
	if ok == false {
		fmt.Println("Fatal error ")
		return ""
	}

	cname,_:= onech["name"].(string)
	return cname
}

// 这里索引要加1
func GetChURL(bookid,cid int) string {
	booklisturl := fmt.Sprintf("http://m.lrts.me/ajax/getPlayPath?entityId=%d&entityType=3&opType=1&sections=[%d]&type=0",bookid,cid+1)
	reqest, err := http.NewRequest("GET", booklisturl, nil) //建立一个请求
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return ""
	}

	reqest.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.95 Safari/537.36")

	for _, v := range cookies {
		reqest.AddCookie(v)
	}

	response, err := gclient.Do(reqest)
	body, err1 := ioutil.ReadAll(response.Body)
	if err1 != nil {
		fmt.Println("Fatal error ", err.Error())
		return ""
	}

	chapters := make(map[string]interface{})
	err = json.Unmarshal(body, &chapters)

	status,_ := chapters["status"].(int)
	if status == 2 {
		print(chapters["msg"])
		return ""
	}

	if status == 0 {
		urls, _:= chapters["list"].([]interface{})
		oneurl, ok := urls[0].(map[string]interface{})
		if ok == false {
			fmt.Println("Fatal error ")
			return ""
		}
		urlname,_:= oneurl["path"].(string)
		return urlname
	}

	return ""

}

func main() {

	gclient = http.Client{}
	cookies = LoginSetCookie("15558630450","aa112233")

	bookinfo := GetBookList(38904)
	chname := GetChapterName(bookinfo,100)
	churl := GetChURL(38904,100)

	filename := fmt.Sprintf("%s.mp3",chname)
	downloadFile(filename,churl)

	print(chname)
	print(churl)

}
