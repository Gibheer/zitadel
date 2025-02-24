package initialise

import (
	"database/sql"
	"errors"
	"testing"
)

func Test_verifyEvents(t *testing.T) {
	err := ReadStmts("cockroach") //TODO: check all dialects
	if err != nil {
		t.Errorf("unable to read stmts: %v", err)
		t.FailNow()
	}

	type args struct {
		db db
	}
	tests := []struct {
		name      string
		args      args
		targetErr error
	}{
		{
			name: "unable to begin",
			args: args{
				db: prepareDB(t,
					expectBegin(sql.ErrConnDone),
				),
			},
			targetErr: sql.ErrConnDone,
		},
		{
			name: "create table fails",
			args: args{
				db: prepareDB(t,
					expectBegin(nil),
					expectExec(createEventsStmt, sql.ErrNoRows),
					expectRollback(nil),
				),
			},
			targetErr: sql.ErrNoRows,
		},
		{
			name: "correct",
			args: args{
				db: prepareDB(t,
					expectBegin(nil),
					expectExec(createEventsStmt, nil),
					expectCommit(nil),
				),
			},
			targetErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createEvents(tt.args.db.db); !errors.Is(err, tt.targetErr) {
				t.Errorf("createEvents() error = %v, want: %v", err, tt.targetErr)
			}
			if err := tt.args.db.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_verifyEncryptionKeys(t *testing.T) {
	type args struct {
		db db
	}
	tests := []struct {
		name      string
		args      args
		targetErr error
	}{
		{
			name: "unable to begin",
			args: args{
				db: prepareDB(t,
					expectBegin(sql.ErrConnDone),
				),
			},
			targetErr: sql.ErrConnDone,
		},
		{
			name: "create table fails",
			args: args{
				db: prepareDB(t,
					expectBegin(nil),
					expectExec(createEncryptionKeysStmt, sql.ErrNoRows),
					expectRollback(nil),
				),
			},
			targetErr: sql.ErrNoRows,
		},
		{
			name: "correct",
			args: args{
				db: prepareDB(t,
					expectBegin(nil),
					expectExec(createEncryptionKeysStmt, nil),
					expectCommit(nil),
				),
			},
			targetErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createEncryptionKeys(tt.args.db.db); !errors.Is(err, tt.targetErr) {
				t.Errorf("createEvents() error = %v, want: %v", err, tt.targetErr)
			}
			if err := tt.args.db.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
