dev:
	go run main.go

deploy:
	fly deploy

backup:
	mv mouji.db mouji.db.bak
	mv mouji.db-shm mouji.db-shm.bak
	mv mouji.db-wal mouji.db-wal.bak

dump:
	fly ssh sftp get ./data/mouji.db ./mouji.db
	fly ssh sftp get ./data/mouji.db-shm ./mouji.db-shm
	fly ssh sftp get ./data/mouji.db-wal ./mouji.db-wal