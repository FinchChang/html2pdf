package html2pdf

import(
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

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

func html2pdf(root string) (err error){
	err = ConverterFromFile(htmlPath + filename + ".htm")
	if err != nil {
		Error.Fatal("transfer html to pdf fail")
	}
	// ---------end of pdf transfer-----------

	// ----------start to optimize and protect the pdf file----------
	err = optAportect(pdfPath+filepath.Base(filename)+".pdf", "H123123123")
	if err != nil {
		Error.Fatal("compress and Encrypt PDF fail")
	}
	// / ----------end of optimize and protect pdf file----------

	elapsed := time.Since(t1)
	fmt.Println("共執行: ", elapsed, "秒")
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

