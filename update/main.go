package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type UpdateAvatarEquipment struct {
	ID              int64
	WebFilePath     string
	WinFilePath     string
	MacFilePath     string
	IosFilePath     string
	AndroidFilePath string
	//CoverFilePath string
}

type UpdateCfgInfo struct {
	id        int64
	resources string
}

//const url = "https://127" // 测试服

const url = "https://127" // 正式服

func main() {
	updateCfgList := []UpdateCfgInfo{
		{
			id:        27,
			resources: "wf_body_01",
		},
		{
			id:        28,
			resources: "wm_body_01",
		},
		{
			id:        30,
			resources: "wm_editor",
		},
		{
			id:        48,
			resources: "wm_lower_008_01_01",
		},
		{
			id:        52,
			resources: "wm_upper_018_01_01",
		},
	}

	base := "./updatesource"
	list := make([]UpdateAvatarEquipment, 0, len(updateCfgList))
	for _, v := range updateCfgList {
		list = append(list, UpdateAvatarEquipment{
			ID:              v.id,
			WebFilePath:     base + "/WebGL/" + v.resources,
			WinFilePath:     base + "/Windows/" + v.resources,
			MacFilePath:     base + "/MacOS/" + v.resources,
			IosFilePath:     base + "/IOS/" + v.resources,
			AndroidFilePath: base + "/Android/" + v.resources,
		})
	}

	for _, v := range list {
		if v.WebFilePath != "" && !fileExists(v.WebFilePath) {
			fmt.Printf("fileExists %s \n", v.WebFilePath)
			return
		}
		if v.WinFilePath != "" && !fileExists(v.WinFilePath) {
			fmt.Printf("fileExists %s \n", v.WinFilePath)
			return
		}
		if v.MacFilePath != "" && !fileExists(v.MacFilePath) {
			fmt.Printf("fileExists %s \n", v.MacFilePath)
			return
		}
		if v.IosFilePath != "" && !fileExists(v.IosFilePath) {
			fmt.Printf("fileExists %s \n", v.IosFilePath)
			return
		}

		if v.AndroidFilePath != "" && !fileExists(v.AndroidFilePath) {
			fmt.Printf("fileExists %s \n", v.AndroidFilePath)
			return
		}
	}

	//上传资源
	for i, v := range list {
		err, res := SendPostFormFile(url, &v)
		if err != nil {
			fmt.Printf("\n%+v\n-----%d----fail------------", res, i)
		} else {
			fmt.Printf("\n %s\n-----%d----ok------------", res, i)
		}
	}
}

func SendPostFormFile(url string, ae *UpdateAvatarEquipment) (error, string) {
	bodBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodBuf)

	// boundary默认会提供一组随机数，也可以自己设置。
	bodyWriter.SetBoundary("Pp7Ye2EeWaFDdAY")
	//boundary :=  body_writer.Boundary()

	// 1. 要上传的数据
	bodyWriter.WriteField("id", fmt.Sprintf("%d", ae.ID))

	if ae.WebFilePath != "" {
		fileField := "web_model_file"
		bodyWriter.WriteField("web_model", fileField)
		_, err := bodyWriter.CreateFormFile(fileField, ae.WebFilePath)
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
	}

	if ae.WinFilePath != "" {
		fileField := "win_model_file"
		bodyWriter.WriteField("win_model", fileField)
		_, err := bodyWriter.CreateFormFile(fileField, ae.WinFilePath)
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

	if ae.MacFilePath != "" {
		fileField := "mac_model_file"
		bodyWriter.WriteField("mac_model", fileField)
		_, err := bodyWriter.CreateFormFile(fileField, ae.MacFilePath)
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
	if ae.IosFilePath != "" {
		fileField := "ios_model_file"
		bodyWriter.WriteField("ios_model", fileField)
		_, err := bodyWriter.CreateFormFile(fileField, ae.IosFilePath)
		if err != nil {
			fmt.Println("CreateFormFile err:", err)
			return err, ""
		}
		fb1, err := ioutil.ReadFile(ae.IosFilePath)
		if err != nil {
			fmt.Println("ReadFile err:", err)
			return err, ""
		}
		bodBuf.Write(fb1)
	}

	if ae.AndroidFilePath != "" {
		fileField := "android_model_file"
		bodyWriter.WriteField("android_model", fileField)
		_, err := bodyWriter.CreateFormFile(fileField, ae.AndroidFilePath)
		if err != nil {
			fmt.Println("CreateFormFile err:", err)
			return err, ""
		}
		fb1, err := ioutil.ReadFile(ae.AndroidFilePath)
		if err != nil {
			fmt.Println("ReadFile err:", err)
			return err, ""
		}
		bodBuf.Write(fb1)
	}

	//if ae.CoverFilePath != "" {
	//	fileField := "cover_file"
	//	bodyWriter.WriteField("cover", fileField)
	//	_, err := bodyWriter.CreateFormFile(fileField, ae.CoverFilePath)
	//	if err != nil {
	//		fmt.Println("CreateFormFile err:", err)
	//		return err, ""
	//	}
	//	fb2, err := ioutil.ReadFile(ae.CoverFilePath)
	//	if err != nil {
	//		fmt.Println("ReadFile err:", err)
	//		return err, ""
	//	}
	//	bodBuf.Write(fb2)
	//}
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
