#!/bin/bash

#NAMEDIR=`pwd|awk -F '/' '{print $3}'`
NAMEDIR="check"
ARPGHOME="/home/${NAMEDIR}"

start() {
	echo  "Starting timemachine "
	${ARPGHOME}/timemachine >/dev/null 2>&1 &
}

stop() {
	echo  "Stopping timemachine "
	killall -15 ${ARPGHOME}/timemachine
}

restart() {
	stop
	start
}


case "$1" in
	start)
		start
	;;

	stop)
		stop
	;;

	restart)
		restart
	;;

	*)
		echo "Usage: $0 {start|stop|restart}"
		exit 1
	;;

esac
