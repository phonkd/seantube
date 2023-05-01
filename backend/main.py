import os
import dotenv
import base64
import json
import subprocess
import glob
from urllib.parse import urlparse
from dotenv import load_dotenv
load_dotenv()
# Variables
input_str = os.getenv("URL")
audio = os.getenv("AUDIO")
audio_format = os.getenv("AUDIO_FORMAT")
video_format = os.getenv("VIDEO_FORMAT")


# Commands
video_cmd_url = "yt-dlp --recode-video " + video_format + " --output seantube_download " + input_str
video_cmd_nourl = "yt-dlp --recode-video " + video_format + " --output seantube_download ytsearch:" + '"' + input_str + '"'
audio_cmd_url = "yt-dlp -x --audio-format " + audio_format + " --output seantube_download " + input_str
audio_cmd_nourl = "yt-dlp -x --audio-format " + audio_format + " --output seantube_download ytsearch:" + '"' + input_str + '"'
print(video_cmd_nourl)
change_directory_to_temp = "cd ./static/temp/"
os.system(change_directory_to_temp)

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
        print("Mode: Audio without url")
        os.system(audio_cmd_nourl)
    else:
        print("Mode: Audio with url")
        os.system(audio_cmd_url)
else:
    if (is_url == False):
        print("Mode: Video without url")
        os.system(video_cmd_nourl)
    else:
        print("Mode: Video with url")
        os.system(video_cmd_url)
        


