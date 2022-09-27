package session

import "day2/log"

func (s *Session) Begin() (err error) {
	log.Info("transaction begin")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session) Commit() (err error) {
	log.Info("transaction commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Rollback() (err error) {
	log.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		return err
	}
	return
}

type TxFunc func(*Session) (interface{}, error)

func (s *Session) Transaction(f TxFunc) (result interface{}, err error) {
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			s.Rollback()
		} else {
			err = s.Commit()
			if err != nil {
				s.Rollback()
			}
		}
	}()

	return f(s)
}
