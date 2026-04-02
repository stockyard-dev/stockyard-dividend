package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Affiliate struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Code string `json:"code"`
	CommissionPct int `json:"commission_pct"`
	TotalEarned int `json:"total_earned"`
	TotalPaid int `json:"total_paid"`
	Referrals int `json:"referrals"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"dividend.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS affiliates(id TEXT PRIMARY KEY,name TEXT NOT NULL,email TEXT DEFAULT '',code TEXT DEFAULT '',commission_pct INTEGER DEFAULT 10,total_earned INTEGER DEFAULT 0,total_paid INTEGER DEFAULT 0,referrals INTEGER DEFAULT 0,status TEXT DEFAULT 'active',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Affiliate)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO affiliates(id,name,email,code,commission_pct,total_earned,total_paid,referrals,status,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Email,e.Code,e.CommissionPct,e.TotalEarned,e.TotalPaid,e.Referrals,e.Status,e.CreatedAt);return err}
func(d *DB)Get(id string)*Affiliate{var e Affiliate;if d.db.QueryRow(`SELECT id,name,email,code,commission_pct,total_earned,total_paid,referrals,status,created_at FROM affiliates WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Email,&e.Code,&e.CommissionPct,&e.TotalEarned,&e.TotalPaid,&e.Referrals,&e.Status,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Affiliate{rows,_:=d.db.Query(`SELECT id,name,email,code,commission_pct,total_earned,total_paid,referrals,status,created_at FROM affiliates ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Affiliate;for rows.Next(){var e Affiliate;rows.Scan(&e.ID,&e.Name,&e.Email,&e.Code,&e.CommissionPct,&e.TotalEarned,&e.TotalPaid,&e.Referrals,&e.Status,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Affiliate)error{_,err:=d.db.Exec(`UPDATE affiliates SET name=?,email=?,code=?,commission_pct=?,total_earned=?,total_paid=?,referrals=?,status=? WHERE id=?`,e.Name,e.Email,e.Code,e.CommissionPct,e.TotalEarned,e.TotalPaid,e.Referrals,e.Status,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM affiliates WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM affiliates`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Affiliate{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ? OR email LIKE ?)"
        args=append(args,"%"+q+"%");args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,email,code,commission_pct,total_earned,total_paid,referrals,status,created_at FROM affiliates WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Affiliate;for rows.Next(){var e Affiliate;rows.Scan(&e.ID,&e.Name,&e.Email,&e.Code,&e.CommissionPct,&e.TotalEarned,&e.TotalPaid,&e.Referrals,&e.Status,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM affiliates GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
