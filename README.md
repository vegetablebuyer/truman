# truman
The word truman comes from the movie *The Truman Show*.At the very beginning,this project is designed to solve the problems:
- How to monitor the specified directories, record every file event
- How to synchronize the directories between multiple servers
---
The solution for file synchronization use the [rsync algorithms](https://rsync.samba.org/tech_report/).

The current progress:
- Successfully tell the minimal differences between the source file and destination file
