ver=latest
if [ $# == 1 ];then
	ver=$1
fi

docker build -t "file-service:$ver" .
docker tag file-service:$ver registry.ainirobot.com/arch/file-service:$ver
docker push registry.ainirobot.com/arch/file-service:$ver
