OS := $(shell uname)

ifeq ($(OS), Linux)
all			:	net server clientLinux
else ifeq ($(OS), Darwin)
all			:	net server clientMac
else
	$(error Unsupported operating system: $(OS))
endif

net			:
	-docker network create MyMusicPlayer

server		:
	-docker stop serv
	-docker rm serv
	-docker build -t server ./MyServer
	-docker run -d --rm -p 9879:9879 -p 6541:6541 --net=MyMusicPlayer --name serv server

clientLinux		:
	-docker stop cli     
	-docker rm cli
	-docker build -t client ./Client
	-docker run -it --name cli --device /dev/snd --net=MyMusicPlayer client

clientMac		:
	-docker stop cli
	-docker rm cli
	-docker build -t client ./Client
	-docker run -it --name cli --device /dev/null --net=MyMusicPlayer client

test			:
	@echo "	Building Client tests"
	@-docker build -t clienttest -f ./Client/Dockerfile_test ./Client &> /dev/null
	@echo "	Running Client tests"
	@docker run --name clitest clienttest
	@echo "	Building Server tests"
	@-docker build -t servertest -f ./MyServer/Dockerfile_test ./MyServer &> /dev/null
	@echo "	Running Server tests"
	@docker run --name servtest servertest
	@docker rm servtest &> /dev/null
	@docker rm clitest &> /dev/null

clean			:
	-docker stop serv
	-docker stop cli
	-docker network rm MyMusicPlayer
	-docker rm serv
	-docker rm cli

fclean			:	clean
	-docker image rm server
	-docker image rm client
	-docker image rm clienttest
	-docker image rm servertest

remake			:	fclean all