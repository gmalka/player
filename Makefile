all			:			net server client

net			:
	docker network create MyMusicPlayer

server		:
	docker build -t server ./MyServer
	docker run -d --rm -p 9879:9879 --net=MyMusicPlayer --name serv server
client		:
	docker build -t client ./Client
	docker run -it --name cli --device /dev/snd --net=MyMusicPlayer client

clean		:
	-docker stop serv
	-docker stop cli
	-docker network rm MyMusicPlayer
	-docker rm serv
	-docker rm cli
	-docker image rm server
	-docker image rm client