package logd

import (
	"database/sql"
	"fmt"
	"runtime"
	"time"
)

/* DESIRES
- ability to choose whether to print to console, file, or db

*/
// init in main, pass through
type Logder struct {
	File string
	Logs []Logd
	Prj  string
	DB   *sql.DB
}

type Logd struct {
	Msg      string
	Caller   string
	LogTime  time.Time
	TimeStr  string
	Err      error
	RowCount int64
}

func (lg *Logder) Log(msg string, err error, rc int64) {
	l := Logd{Msg: msg, Err: err, LogTime: time.Now()}
	l.SetCaller()
	l.FormatTime()
	lg.Logs = append(lg.Logs, l)
	l.WriteLog()
	if lg.DB != nil {
		lg.DB.Exec(
			`insert into log.log (prj, msg, ltime, ltstr, caller, err, rc) values
			($1, $2, $3, $4, $5, $6, $7)
		`, lg.Prj, l.Msg, l.LogTime, l.TimeStr, l.Caller, l.Err, l.RowCount)
	}
}

func (l *Logd) WriteLog() {
	if l.Err != nil {
		l.ErrLog()
	}
	fmt.Printf("%s | caller: %s\n%s\n", l.TimeStr, l.Caller, l.Msg)
}

func (l *Logd) ErrLog() {
	fmt.Printf("** %s | ** ERROR | caller: %s\n%s\n** %v\n",
		l.TimeStr, l.Caller, l.Msg, l.Err)
}

func (l *Logd) FormatTime() {
	l.TimeStr = l.LogTime.Format("010206_150405")
}

func (l *Logd) SetCaller() {
	pc, _, _, _ := runtime.Caller(2)
	l.Caller = runtime.FuncForPC(pc).Name()
}
