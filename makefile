dev:
	go run main.go

deploy:
	fly deploy

dump:
	mv mouji.db mouji.db.bak
	mv mouji.db-shm mouji.db-shm.bak
	mv mouji.db-wal mouji.db-wal.bak
	fly ssh sftp get ./data/mouji.db ./mouji.db
	fly ssh sftp get ./data/mouji.db-shm ./mouji.db-shm
	fly ssh sftp get ./data/mouji.db-wal ./mouji.db-wal