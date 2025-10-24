package logd

import (
	"database/sql"
	"fmt"
	"runtime"
	"sync"
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
	mu   *sync.Mutex
}

type Logd struct {
	Msg      string
	Caller   string
	LogTime  time.Time
	TimeStr  string
	Err      error
	RowCount int64
}

func (lg *Logder) Log(msg string, err error, rc *int64) {
	var row_count int64
	if rc == nil {
		row_count = 0
	} else {
		row_count = *rc
	}
	l := Logd{Msg: msg, Err: err, LogTime: time.Now(), RowCount: row_count}
	l.SetCaller()
	l.FormatTime()
	lg.Logs = append(lg.Logs, l)
	l.WriteLog()
	if rc != nil {
		if lg.DB != nil {
			lg.LogToDB(&l)
		}
	}
}

func (lg *Logder) CCLog(msg string, err error, rc *int64, mu *sync.Mutex) {
	var row_count int64
	if rc == nil {
		row_count = 0
	} else {
		row_count = *rc
	}
	l := Logd{Msg: msg, Err: err, LogTime: time.Now(), RowCount: row_count}
	l.SetCaller()
	l.FormatTime()
	lg.Logs = append(lg.Logs, l)
	l.WriteLog()
	if rc != nil {
		if lg.DB != nil {
			lg.CCLogToDB(&l, mu)
		}
	}
}

func (lg *Logder) LogToDB(l *Logd) {
	if lg.DB != nil {
		lg.DB.Exec(
			`insert into log.log (prj, msg, ltime, ltstr, caller, err, rc) values
			($1, $2, $3, $4, $5, $6, $7)
		`, lg.Prj, l.Msg, l.LogTime, l.TimeStr, l.Caller, l.Err, l.RowCount)
	}
}

func (lg *Logder) CCLogToDB(l *Logd, mu *sync.Mutex) {
	if lg.DB != nil && lg.mu != nil {
		lg.mu.Lock()
		_, err := lg.DB.Exec(
			`insert into log.log (prj, msg, ltime, ltstr, caller, err, rc) values
			($1, $2, $3, $4, $5, $6, $7)
		`, lg.Prj, l.Msg, l.LogTime, l.TimeStr, l.Caller, l.Err, l.RowCount)
		if err != nil {
			fmt.Println(err)
		}
		lg.mu.Unlock()
	} else {
		return
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
