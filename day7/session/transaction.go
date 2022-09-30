package session

import (
	"day7/log"
)

func (s *Session) Begin() (err error) {
	log.Info("transaction begin..")
	s.tx, err = s.db.Begin()
	return
}

func (s *Session) Commit() error {
	log.Info("commit ..")
	return s.tx.Commit()
}

func (s *Session) Rollback() error {
	log.Info("rollback ..")
	return s.tx.Rollback()
}

type txFunc func(*Session) error

func (s *Session) Transaction(f txFunc) error {
	err := s.Begin()
	if err != nil {
		log.Error("transaction err")
	}

	defer func() {
		if p := recover(); p != nil {
			log.Info("rollback")
			s.Rollback()
		} else if err != nil {
			s.Rollback()
		} else {
			s.Commit()
		}
	}()
	return f(s)

}
