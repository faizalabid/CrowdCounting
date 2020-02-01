package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	//"image"
	//"image/draw"
    "strconv"

	//"image/png"
	//"image/color"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	//"math"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	no         int
	sqlversion string
)

func main() {
	http.HandleFunc("/Show", ImageClassHandler)

	fmt.Printf("Starting server for	 HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

//QueryExecDB for query executeable without data
func QueryExecDB(qry string) {

	conn := GetConn()
	if conn != nil {
		_, err := conn.Exec(qry)
		if err != nil {
			fmt.Println(err)
		}
	}

	defer conn.Close()
}

//GetConn connect to db
func GetConn() *sql.DB {
	condb, errdb := sql.Open("sqlserver", "server=localhost;user id=sa;password=deemes;database=ThesisData")
	if errdb != nil {
		fmt.Println(" Error open db:", errdb.Error())
	}
	return condb
}

//ResetDot reset dot sesuai nomor
func ResetDot (no string){
	Q := "delete GROUNDTRUTH_DATA WHERE type = 'Train' and seq = '" + no + "'"

	QueryExecDB(Q)
}

//SaveRedDot SAVE
func SaveRedDot(x string, y string, filenm string, no string) {
	Q := "INSERT INTO GROUNDTRUTH_DATA(FILENAME,X,Y,SEQ,TYPE) " + "VALUES('" + filenm + "','" + x + "','" + y + "','" + "null" + "','Train')"

	QueryExecDB(Q)
}

//ImageClassHandler Screen
func ImageClassHandler(w http.ResponseWriter, r *http.Request) {
	img := ""

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
    
	k := r.PostFormValue("Mode")
	X := r.PostFormValue("X")
	Y := r.PostFormValue("Y")
	// PostBack := r.Form.Get("iscallback")
	if k == "1" {
		no = no + 1
	} else if k == "2" {
		no = no - 1
	} else if k == "3" {
		no = no + 0
	} else if k == "4"{
		ResetDot(strconv.Itoa(no))
	} else {
		no = 0
	}

	RealImage := Readfile(no)
	filenm := RealImage.Name()

	img = EncodeImage(RealImage)

	if X != "" && Y != "" {
		SaveRedDot(X, Y, filenm, strconv.Itoa(no))
	}

	var data = map[string]string{"image": img}

	t, _ := template.ParseFiles("Screen.html")
	t.Execute(w, data)
}

//RetriveDot Retrive dot from Database
// func RetriveDot(myImage image.Image, FILENAME string) image.Image {
// 	imageRect := myImage.Bounds()
// 	Image1 := image.NewYCbCr(imageRect, image.YCbCrSubsampleRatio420)
// 	conn := GetConn()
// 	rows, err := conn.Query("select X, Y FROM GROUNDTRUTH_DATA WHERE FILENAME='" + FILENAME + "'")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	var (
// 		X float64
// 		Y float64
// 	)

// 	for rows.Next() {
// 		if err := rows.Scan(&X, &Y); err == nil {
// 			xin := int(math.Round(X))
// 			yin := int(math.Round(Y))

// 			var colRED = color.RGBA{255, 0, 0, 255}
// 			Image1.Set(xin, yin, colRED)
// 		}
// 	}

// 	return Image1
// }

//EncodeImage file image to base64
func EncodeImage(images os.FileInfo) string {

	ImgFile, err := os.Open("./PhotoFolder/faild/" + images.Name())
	if err != nil {
		// Handle error
	}

	defer ImgFile.Close()
	myImage, err := jpeg.Decode(ImgFile)
	if err != nil {
		log.Print("File Bukan PNG")
	}

	//myImage1 := RetriveDot(myImage, images.Name())

	var buff bytes.Buffer
	err10 := jpeg.Encode(&buff, myImage, &jpeg.Options{95})
	if err10 != nil {
		log.Print("Eror Encode")
	}

	encodedString := base64.StdEncoding.EncodeToString(buff.Bytes())
	return encodedString
}

//Readfile => FUngsi Baca File from directory
func Readfile(i int) os.FileInfo {
	files, err1 := ioutil.ReadDir("./PhotoFolder/faild/")
	if err1 != nil {
		log.Fatal(err1)
	}

	if i < 0 {
		i = 0
		no = 0
	} else if i > len(files)-1 {
		i = len(files) - 1
		no = len(files) - 1
	}

	FileImage := files[i]
	return FileImage
}

//save data and move file image
func save(i int, opt string) {
	file := Readfile(i)
	var dst string
	dstCrowd := "./PhotoCrowded/" + file.Name()
	dstNCrowd := "./PhotoNonCrowded/" + file.Name()

	source := "./PhotoFolder/" + file.Name()
	if opt == "C" {
		dst = dstCrowd
	} else {
		dst = dstNCrowd
	}

	//delete file bilang sudah ada di folder
	errAll := os.Remove(dstCrowd)
	if errAll != nil {
	}

	errAll = os.Remove(dstNCrowd)
	if errAll != nil {
	}

	err := Copy(source, dst)
	if err != nil {
		log.Fatal(err)
	}
}

//Copy file
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
