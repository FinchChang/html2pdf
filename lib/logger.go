package logger

import(
	"log"
)

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

func main(){
	SetLogger()
}
