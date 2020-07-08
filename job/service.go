package job

type Service interface {
	ServerJobPop(sid, jid string)
}
