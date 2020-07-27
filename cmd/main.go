package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"golang.org/x/net/html"
)

var pdfPath, htmlPath, pdopw string

func initlaParamters() {
	pdfPath = "./pdffiles/"
	htmlPath = "./testfiles/"
	pdopw = "example@pw"
}

func init() {
	initlaParamters()
}

// --------------start settin logger--------------

var (
	//Info useed for debug to recored some
	Info *log.Logger
	//Warning useed for record warning
	Warning *log.Logger
	//Error useed for record Error in log file
	Error *log.Logger
)

// SetLogger is used to define the log format
func SetLogger() {
	errFile, err := os.OpenFile("./log/errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("logger setting fail", err)
	}

	Info = log.New(os.Stdout, "[Info]   ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "[Warning]", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, errFile), "[Error]  ", log.Ldate|log.Ltime|log.Lshortfile)

	defer errFile.Close()
}

// --------------end of setting logger--------------

//remove_decentant is used to remove the tag element
// ref:https://golang.hotexamples.com/examples/code.google.com.p.go.net/html/Node/RemoveChild/golang-node-removechild-method-examples.html
func removeDecentant(n *html.Node, tag string) {
	child := n.FirstChild
	for child != nil {
		if child.Type == html.ElementNode && child.Data == tag {
			next := child.NextSibling
			n.RemoveChild(child)
			child = next
		} else {
			removeDecentant(child, tag)
			child = child.NextSibling
		}
	}
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func renderNodeByte(n *html.Node) []byte {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return []byte(buf.String())
}

func main() {

	t1 := time.Now() // get current time

	
	//html2pdf.ConvertHTML("D:\\statment_test\\Data22\\20191126ＮoJS.htm","H123123123")

	// ---- start of products and customers ---
	//---- ref:http://www.hatlonely.com/2018/03/11/golang-%E5%B9%B6%E5%8F%91%E7%BC%96%E7%A8%8B%E4%B9%8B%E7%94%9F%E4%BA%A7%E8%80%85%E6%B6%88%E8%B4%B9%E8%80%85/index.html ---
	/*
	t1 = time.Now() // get current time
	fmt.Println("Start Time = ", time.Now())
	nCPU := runtime.NumCPU()
	var wgp sync.WaitGroup
	var wgc sync.WaitGroup
	stop := false
	products := make(chan string, 10)
	go producer(&wgp, products, &stop)
	wgp.Add(1)
	// 建立nCPU個消費者(也可建立生產者)
	for i := 0; i <= nCPU; i++ {
		go consumer(&wgc, products)
		wgc.Add(1)
	}
	time.Sleep(time.Duration(1) * time.Second)
	stop = true     // 設置生產者 終止flag
	wgp.Wait()      // 等待生產者 退出
	close(products) // 關閉Channel
	wgc.Wait()      // 等待消費者退出
	//---- end of products and customers ---
	pEndTime(t1)
	fmt.Println("End Time = ", time.Now())
	*/
	
	
	// ---------start remmove javascript element from html file----------
	/*
		file, err := osOpe(htmlPath + filename + ".htm") // For read access.
		if err != nil {
			og.Fatal(err)
		}
		defer fie.Close()

		doc, err := htm.Parse(file)
		if err = nil {
			/ ...
		/ }

		newPage = template.ew("templ")
		visit(newPage, doc)

		removeDecentant(oc, "script")

		nwfilename := htmlPath + filename + "2.htm"
		err = ioutil.WrteFile(newfilename, renderNodeByte(doc), 0644)
		if rr != nil {
			og.Fatal(err)
		}

		os.
		clear_dom(doc)
		fmt.Println(renderNode(doc))
	*/
	// ----------end of remove javascript element---------

	elapsed := time.Since(t1)
	fmt.Println("共執行: ", elapsed, "秒")

}

//Contains return true if x is exist in array a
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// start html2pdf

// OptimAndEncryptPDF is used to encrypt and set password to the pdf
func OptimAndEncryptPDF(filename, password string) (err error) {
	err = api.OptimizeFile(filename, "", nil)
	if err != nil {
		Error.Println("compress PDF fail")
		return err
	}
	//change user and owner password
	conf := pdfcpu.NewAESConfiguration(password, pdopw, 256)
	//encrypt file and set password
	err = api.EncryptFile(filename, "", conf)
	if err != nil {
		Error.Println("encrypt PDF fail")
		return err
	}
	return
}

func html2pdf(root string) (err error) {
	err = ConverterFromFile(root)
	if err != nil {
		Error.Fatal("transfer html to pdf fail")
	}
	// ---------end of pdf transfer-----------

	// ----------start to optimize and protect the pdf file----------
	filename := "test"
	err = optAportect(pdfPath+filepath.Base(filename)+".pdf", "test")
	if err != nil {
		Error.Fatal("compress and Encrypt PDF fail")
	}
	// / ----------end of optimize and protect pdf file----------
	return
}

// List all file in root and filger by extension
// extension EX: []string {".htm",".html"}
func listFiles(root string, extension []string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if extension == nil {
			files = append(files, path)
		} else {
			if Contains(extension, filepath.Ext(path)) {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return
}

//ConverterFromFile transfer the html file to pdf file
func ConverterFromFile(FilePath string) (err error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		Error.Println("Create PDF fail")
		return
	}

	pdfg.AddPage(wkhtmltopdf.NewPage(FilePath))
	err = pdfg.Create()
	if err != nil {
		Error.Println("Create PDF fail on add page")
		return
	}
	var extension = filepath.Ext(FilePath)
	var name = FilePath[0 : len(FilePath)-len(extension)]

	err = pdfg.WriteFile(pdfPath + filepath.Base(name) + ".pdf")
	fmt.Println("PDF save at = " + pdfPath + filepath.Base(name) + ".pdf")
	if err != nil {
		Error.Println("save PDF fail")
		return
	}
	return
}

//optAportect is used to optimize and set password the PDF file
func optAportect(filename, password string) error {
	//-------pdfcpu by api-------
	// optimize the PDF file
	err := OptimAndEncryptPDF(filename, password)
	//-------end of pdfcpu by api-------
	if err != nil {
		Error.Fatal(err)
	}
	return err
}

// execPdfCPU execute the windows command
func execPdfCPU(input []string) error {
	excmd := exec.Command("pdfcpu.exe", input...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	excmd.Stdout = &out
	excmd.Stderr = &stderr
	err := excmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	// fmt.Println(out.String())
	return err
}

//end of html2pdf
