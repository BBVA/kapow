kapow.gif: tty.gif
	cp tty.gif kapow.gif
tty.gif: ttyrecord
	WINDOWID=$(shell xdotool getwindowfocus) ttygif ttyrecord
ttyrecord:
	stty cols 110
	ttyrec -e bash ttyrecord
