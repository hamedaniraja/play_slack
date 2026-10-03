# play_slack
We can run this script by:
```bash
notify_slack [OPTIONS]
```
[OPTIONS]:  
&emsp; -help&emsp; Shows this text.

&emsp; -action&emsp; Could be either "sendMessage" or "cronJobMessage".   
&emsp; &emsp; "sendMessage": Will send a message rightaway to a slack public channel or a private group.  
&emsp; &emsp; "cronJobMessage": Will schedule a cronjob on the node for sending a messgae later. default value is: ""

&emsp; -dateTimeStr&emsp; This option will be used only if we have chosen "cronJobMessage" action. Format of this string will be similar to "2024-06-13 14:30:00 PDT". default value is: ""
  
&emsp; -channel&emsp; This option identifies the public channel or private group. Examples: public channel "#channle_name" or private group "private_group_name"
                default value is: ""
  
&emsp; -icon&emsp; The icon or emoji which will be used for message. Example: ":approved-9724:" or ":failed:". default value is: ":approved-9724:"
  
&emsp; -title&emsp; Title of the message. default value is: "Test Title"

&emsp; -bot&emsp; This is the Slack bot which we use for sending message Default value is: "webhookbot"

&emsp; -text&emsp; Text of the message. default value is: "Test Message Text"

How to build the binary on Linux and Mac:
```bash
$ cd /path/to/project/dir
$ go build -o notify_slack notify_slack.go 
```

Example:
```bash
$ notify_slack -action "sendMessage" \
-title "Test title" \
-text "This is a test." \
-channel "pe-eloqua-serviceops" \
-icon ":failed:"
```


