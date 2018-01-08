export GOPATH=`pwd`
go get github.com/go-sql-driver/mysql
cd src
go build -o radar_ft_collect
mv radar_ft_collect ../bin/

# 运行方式
# nohup ./radar_ft_collect > myout.file 2>&1 &
