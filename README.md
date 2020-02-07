# Text to JSON

This program allows you to dynamically generate JSON arrays from structured text files.
The text files need to be of similar structure as the "test" file in this repo.
### KEY: VALUE

Running setup script will make sure your executable directory has a 'data' directory, where it will store the new json file. The script will also schedule a monitoring script on crontab, to make execute updates whenever our source has changed.

TODO:
*	Validate source text structure before converting to JSON
*	Restructure code
*	Prompt 'data' dir from user and config setup.sh accordingly.
*	Add monitoring script to this repo