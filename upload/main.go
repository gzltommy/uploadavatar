package main

import (
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type AvatarEquipment struct {
	LineN int

	WebFilePath string
	WinFilePath string
	MacFilePath string

	CoverFilePath string
	PartID        int32
	Name          string
	RoleType      string
	IsBase        int8
}

//const url = "http://*/avatar/equipment/add"

func main() {
	err, list := LoadSheetFile("./source/ABPublish", "./source/avatar2.xlsx", "./source/Icon")
	if err != nil {
		panic(err)
	}

	//for _, v := range list {
	//	fmt.Println("--------------------------------------------")
	//	fmt.Println(v.LineN, v.WebFilePath, v.MacFilePath, v.WinFilePath, v.CoverFilePath, v.Name, v.RoleType)
	//}

	fmt.Println("===========", len(list))

	//上传资源
	for i, v := range list {
		err, res := SendPostFormFile(url, v)
		if err != nil {
			fmt.Printf("\n%+v\n-----%d----fail------------", res, i)
			return
		} else {
			fmt.Printf("\n %s\n-----%d----ok------------", res, i)
		}
	}
}

func SendPostFormFile(url string, ae *AvatarEquipment) (error, string) {
	bodBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodBuf)

	// boundary默认会提供一组随机数，也可以自己设置。
	bodyWriter.SetBoundary("Pp7Ye2EeWaFDdAY")
	//boundary :=  body_writer.Boundary()

	// 1. 要上传的数据
	bodyWriter.WriteField("part_id", fmt.Sprintf("%d", ae.PartID))
	bodyWriter.WriteField("name", ae.Name)
	bodyWriter.WriteField("role_type", ae.RoleType)
	bodyWriter.WriteField("is_base", fmt.Sprintf("%d", ae.IsBase))

	// 2. 读取文件
	_, err := bodyWriter.CreateFormFile("web_model_file", ae.WebFilePath)
	if err != nil {
		fmt.Println("CreateFormFile err:", err)
		return err, ""
	}
	fb1, err := ioutil.ReadFile(ae.WebFilePath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return err, ""
	}
	bodBuf.Write(fb1)

	{
		_, err := bodyWriter.CreateFormFile("win_model_file", ae.WinFilePath)
		if err != nil {
			fmt.Println("CreateFormFile err:", err)
			return err, ""
		}
		fb1, err := ioutil.ReadFile(ae.WinFilePath)
		if err != nil {
			fmt.Println("ReadFile err:", err)
			return err, ""
		}
		bodBuf.Write(fb1)
	}

	{
		_, err := bodyWriter.CreateFormFile("mac_model_file", ae.MacFilePath)
		if err != nil {
			fmt.Println("CreateFormFile err:", err)
			return err, ""
		}
		fb1, err := ioutil.ReadFile(ae.MacFilePath)
		if err != nil {
			fmt.Println("ReadFile err:", err)
			return err, ""
		}
		bodBuf.Write(fb1)
	}

	if ae.IsBase == 0 {
		_, err = bodyWriter.CreateFormFile("cover_file", ae.CoverFilePath)
		if err != nil {
			fmt.Println("CreateFormFile err:", err)
			return err, ""
		}
		fb2, err := ioutil.ReadFile(ae.CoverFilePath)
		if err != nil {
			fmt.Println("ReadFile err:", err)
			return err, ""
		}
		bodBuf.Write(fb2)
	}
	bodyWriter.Close() // 结束整个消息 body

	//
	reqReader := io.MultiReader(bodBuf)
	req, err := http.NewRequest("POST", url, reqReader)
	if err != nil {
		fmt.Println("NewRequest err:", err)
		return err, ""
	}
	// 添加 Post 头
	req.Header.Set("Connection", "close")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	req.ContentLength = int64(bodBuf.Len())

	// 发送消息
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Do err:", err)
		return err, ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll err:", err)
		return err, ""
	}
	return nil, string(body)
}

/*===================================================================================================================*/
//var modelMap = map[string]int32{
//	"boy":   1, //二次元男
//	"girl":  2, //二次元女
//	"westm": 3, //欧美男（待添加）
//	"westf": 4, //欧美女（待添加）
//}

var platformMap = map[string]int32{
	"WebGL":   1,
	"Windows": 2,
	"MacOS":   3,
}

func LoadSheetFile(modelBasePath, sheetFile, coverDir string, sheetOpt ...string) (error, []*AvatarEquipment) {
	f, err := excelize.OpenFile(sheetFile)
	if err != nil {
		fmt.Println("OpenFile error:", err)
		return err, nil
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Close error:", err)
		}
	}()

	sheet := "Sheet1"
	if len(sheetOpt) == 1 {
		sheet = sheetOpt[0]
	}
	rows, err := f.GetRows(sheet)
	if err != nil {
		fmt.Println("GetRows error:", err)
		return err, nil
	}

	aList := make([]*AvatarEquipment, 0, 100)
	for rn, row := range rows {
		// 去掉前面 1 行
		if rn > 0 {
			var (
				modelFileName string
				coverFileName string
				roleType      string
				partID        int32
				name          string
			)
			for cn, colCell := range row {
				switch cn {
				case 0: // model 资源文件
					modelFileName = colCell
				case 1: // name
					name = colCell
				case 2: // roleType
					roleType = colCell
				case 3: // partID
					if _v, err := strconv.Atoi(colCell); err == nil {
						partID = int32(_v)
					} else {
						fmt.Println("Atoi err:", err, colCell)
						return err, nil
					}
				case 4: // cover 资源文件
					coverFileName = colCell

					//case 5: // cover 资源文件
					//	coverFileName = colCell
				}
			}

			isBase := int8(1)
			coverFilePath := ""
			if coverFileName != "" {
				coverFilePath = coverDir + "/" + coverFileName + ".png"
				if !fileExists(coverFilePath) {
					fmt.Printf("fileExists %s \n", coverFilePath)
					return fmt.Errorf("file(%s) not exists", modelFileName), nil
				}
				isBase = 0
			}
			ae := &AvatarEquipment{
				LineN:         rn + 1,
				CoverFilePath: coverFilePath,
				PartID:        partID,
				Name:          name,
				RoleType:      roleType,
				IsBase:        isBase,
			}
			for platformS, _ := range platformMap {
				modelFilePath := modelBasePath + "/" + platformS + "/" + modelFileName
				if !fileExists(modelFilePath) {
					fmt.Printf("fileExists %s \n", modelFilePath)
					return fmt.Errorf("file(%s) not exists", modelFileName), nil
				}
				if platformS == "WebGL" {
					ae.WebFilePath = modelFilePath
				} else if platformS == "Windows" {
					ae.WinFilePath = modelFilePath
				} else if platformS == "MacOS" {
					ae.MacFilePath = modelFilePath
				}
			}
			aList = append(aList, ae)
		}
	}
	return nil, aList
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
