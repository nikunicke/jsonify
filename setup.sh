# if [ "$#" -ne 1 ] ;
# then
# 	echo "Usage: $0 [DIR]" >&2
# 	exit 1
# fi
# if [ ! -d "$1" ] ;
# then
# 	if [ -f "$1" ] ;
# 	then
# 		echo "$1: Not a directory"
# 	else
# 		echo "$1: No such directory"
# 	fi
# fi
mkdir data
if [ "$(grep 'monitor_file' /var/spool/cron/crontabs/*)" ] ;
then
	echo "Service already setup on crontab"
else
	echo "Adding service to crontab"
	(crontab -l 2>/dev/null; echo "* * * * * ~/scripts/monitor_file.sh /var/lib/dpkg/status") | crontab -
fi