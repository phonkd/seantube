import os
import base64
import json
import subprocess
import glob
from urllib.parse import urlparse
# Variables
input_str = "no tomorrow mc orsen"
audio = True
audio_format = "mp3"
video_format = ""

# Commands
video_cmd_url = "yt-dlp --format " + video_format + " input_str"
video_cmd_nourl = "yt-dlp --format " + video_format + " ytsearch:" + '"' + input_str + '"'
audio_cmd_url = "yt-dlp -x --audio-format " + audio_format + " " + input_str
audio_cmd_nourl = "yt-dlp -x --audio-format " + audio_format + " ytsearch:" + '"' + input_str + '"'



# Functions
def is_valid_url(url):
    try:
        result = urlparse(url)
        return all([result.scheme, result.netloc])
    except ValueError:
        return False
is_url = is_valid_url(input_str) # Check if the input string is a url or just some keywords
if (audio == True):
    if (is_url == False):
        os.system(audio_cmd_nourl)
    else:
        os.system(audio_cmd_url)
else:
    if (is_url == False):
        os.system(video_cmd_nourl)
    else:
        os.system(video_cmd_url)
        


