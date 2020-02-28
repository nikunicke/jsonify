if [ "$#" -ne 1 ] || ! [ -f "$1" ];
then
	echo "Usage: $0 [FILE]" >&2
	exit 1
fi

FILE=/var/tmp/checksum
TARGET=$1
MD5=$(md5sum $TARGET)
if [ ! -f $FILE ] ;
then
	echo "$MD5" > $FILE
fi
if [ "$MD5" != "$(cat $FILE)" ] ;
then
	echo "$MD5" > $FILE
	echo "[ $(date) ]: '$TARGET' has been modified" >> /var/log/monitor_file.log
	/root/go/go_tools/dpkg_json/dpkg_json
fi
