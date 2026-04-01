package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-dividend/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){list,_:=s.db.List();if list==nil{list=[]store.Referral{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var ref store.Referral;json.NewDecoder(r.Body).Decode(&ref);if ref.Code==""||ref.Affiliate==""{writeError(w,400,"code and affiliate required");return};ref.Active=true;s.db.Create(&ref);writeJSON(w,201,ref)}
func(s *Server)handleSignup(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var su store.Signup;json.NewDecoder(r.Body).Decode(&su);su.ReferralID=id;if su.Email==""{writeError(w,400,"email required");return};s.db.RecordSignup(&su);writeJSON(w,201,su)}
func(s *Server)handlePay(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.MarkPaid(id);writeJSON(w,200,map[string]string{"status":"paid"})}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
