package logs

type Outputer interface {
   Write(*LogData)
   Close()
}