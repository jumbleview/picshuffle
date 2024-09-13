# picshuffle

It is the  utility to set Windows desktop background

picshuffle.exe [-s] [-l] path_to_folder_or_jpeg_file

It starts with single argument: path to folder with jpeg files or dierct path to jpeg file.
In case before an argument there is the flag -s program starts in silent mode (without console).

In case argumentis path to file program reads jpeg file, encodes it, crops image so it's width-to-height ratio matches width-to-height ratio of computer monitor, saves the image as bmp file and sets it as a screen background.

In case argumentis path to director with mutiple jpeg files program reads content of directory 
and randomly picks one file from that directory to set as desktop  backgorund. Program tracks selected files in storage picshuffle.db so once selected fle willl not be selected again till all files from this directory are taken.


It is pure  Go program. It should be compiled with command  
* go build -ldflags -H=windowsgui

That allows to assign windows console dynamically or suppress console  by  flag -s.

To track selected files program uses bolt embeded key-value store  (go.etcd.io/bbolt).
If started with flag [-l] program prints list of files used as a backgorund withinlast 365 days.

