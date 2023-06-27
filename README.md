<h1 align="center">
  🎶 Mp3 Player
  </h1>
  
  ## 💡 About the project:
  	Simple mp3 music player writen with golang.
  	Writen with client-server architecture.
  	On client side you can use Client Line Interface to manage
  	songs. Also Client can connect server by GRPC to ask a remote
  	song from it.
  	You can also use HTTP to save/delete song on server or get list of it.
   
   ## 🛠 Testing and Usage:
    # Clone the project and access the folder
      git clone https://github.com/gmalka/player.git
    # Perform make to run server in container
      make server
    # Perform make to run client without container
      make clientLocal
    # Perform make to build and run project in containers
      make
     ❗ Note: Music from container may not always be displayed correctly ❗
    # Clean up
      make fclean
    
   ## 🎬 Demonstration(example):
![hippo](https://github.com/gmalka/player/assets/94842625/e70e1406-79ae-4d21-a6e4-615312d6e3b7)
