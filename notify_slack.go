package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var action string
var channel string
var icon string
var title string
var color string
var bot string
var text string
var dateTimeStr string
var help bool

func slackWebHook() error {
	parameters_string := "{"
	parameters_string += "'channel': '" + channel + "', "
	parameters_string += "'icon_emoji': '" + icon + "', "
	parameters_string += "'attachments': [{'title': '" + title + "', "
	parameters_string += "'color': '" + color + "', "
	parameters_string += "'username': '" + bot + "', "
	parameters_string += "'text': '" + text + "'}]}"

	parameters_bytes := []byte(parameters_string)

	parameters_ioreader := bytes.NewReader(parameters_bytes)

	resp, err := http.Post("https://hooks.slack.com/services/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",
		"application/json",
		parameters_ioreader)
	if err != nil {
		log.Fatalf("error: can't call hooks.slack.com")
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	var result string
	if err != nil {
		log.Fatal(err)
		result = "Fatal"
	} else {
		result = string(bodyBytes)
	}
	if result != "ok" {
		return errors.New("Failed")
	}
	fmt.Println(result)
	return nil
}

func cronWebHook() error {

	currentOS := runtime.GOOS
	if currentOS != "linux" && currentOS != "darwin" {
		fmt.Println("This action is available only on Macos and Linux Operating Systems.")
		return nil
	}

	// List of cron-jobs
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.CombinedOutput()
	var currentCronJobs string
	if err != nil {
		currentCronJobs = ""
	} else {
		currentCronJobs = string(output)
	}
	//fmt.Println(currentCronJobs)

	layout := "2006-01-02 15:04:05 MST"
	t, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	/*
	   tMinus7Days := t.AddDate(0, 0, -7)
	   tMinus3Hours := t.Add(-3 * time.Hour)
	   fmt.Println("t ==> ", t)
	   fmt.Println("tMinus7Days ==> ", tMinus7Days)
	   fmt.Println("tMinus3Hours ==> ", tMinus3Hours) */

	// Extract components to construct cron format
	minute := t.Minute()
	hour := t.Hour()
	day := t.Day()
	month := int(t.Month())
	newCronTime := fmt.Sprintf("%d %d %d %d * ", minute, hour, day, month)

	// fmt.Println(newCronTime)

	scPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Println(err)
		return err
	}

	newCronJob := newCronTime
	newCronJob += scPath + " -action sendMessage "
	newCronJob += " -channel '" + channel + "'"
	newCronJob += " -icon '" + icon + "'"
	newCronJob += " -title '" + title + "'"
	newCronJob += " -color " + color
	newCronJob += " -bot " + bot
	newCronJob += " -text '" + text + "'"

	// fmt.Println(newCronJob)
	if strings.Contains(currentCronJobs, newCronJob) {
		fmt.Println("This job is already exist")
	} else {
		updatedCron := string(currentCronJobs) + newCronJob + "\n"
		tmpFile, err := os.CreateTemp("", "crontab")
		if err != nil {
			fmt.Println("Error creating temporary file:", err)
			os.Exit(1)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write([]byte(updatedCron)); err != nil {
			fmt.Println("Error writing to temporary file:", err)
			os.Exit(1)
		}
		if err := tmpFile.Close(); err != nil {
			fmt.Println("Error closing temporary file:", err)
			os.Exit(1)
		}
		cmd = exec.Command("crontab", tmpFile.Name())
		if err := cmd.Run(); err != nil {
			fmt.Println("Error installing new crontab:", err)
			os.Exit(1)
		}
		fmt.Println("New cron job added successfully.")
	}
	return nil
}

func showHelp() {
	str := `
We can run this command in two ways:
  notify_slack [OPTIONS]
  notify_slack -help

[OPTIONS]:
  -help         Shows this text.

  -action       Could be either "sendMessage" or "cronJobMessage". 
                "sendMessage": Will send a message rightaway to a slack public channel or a private group.
                "cronJobMessage": Will schedule a cronjob on the node for sending a messgae later
                default value is: ""

  -dateTimeStr  This option will be used only if we have chosen "cronJobMessage" action. Format of this string will be similar to "2024-06-13 14:30:00 PDT"
                default value is: ""
  
  -channel      This option identifies the public channel or private group. 
                Examples: public channel "#public-chanenel-name" or private group "private-group-name"
                default value is: ""
  
  -icon         The icon or emoji which will be used for message. 
                Example: ":approved-9724:" or ":failed:"
                default value is: ":approved-9724:"
  
  -title        Title of the message.
                default value is: "Test Title"

  -bot          This is the Slack bot which we use for sending message.
                Default value is: "webhookbot"

  -text         Text of the message.
                default value is: "Test Message Text"

Example:
$ ./notify_slack -action sendMessage -title "Test title" -text "This is a test." -channel "channel-name" -icon ":failed:"
`
	fmt.Println(str)
}

func main() {

	flag.StringVar(&action, "action", "", "Action")
	flag.StringVar(&dateTimeStr, "dateTimeStr", "", "Date time string including time-zone like \"2024-06-13 14:30:00 PDT\"")
	flag.StringVar(&channel, "channel", "#all-keivan39s-personal", "Slack channel name")
	flag.StringVar(&icon, "icon", ":approved-9724:", "Message Icon/Emoji")
	flag.StringVar(&title, "title", "Test Title", "Message Title")
	flag.StringVar(&color, "color", "info", "Message Color")
	flag.StringVar(&bot, "bot", "webhookbot", "Slack BOT")
	flag.StringVar(&text, "text", "Test Message Text", "Message Text")
	flag.BoolVar(&help, "help", false, "Showing help")
	flag.Parse()

	switch action {
	case "sendMessage":
		slackWebHook()
	case "cronJobMessage":
		cronWebHook()
	case "":
		if help {
			showHelp()
		} else {
			fmt.Println("Undefined Action")
			showHelp()
		}
	default:
		fmt.Println("Undefined Action")
		showHelp()
	}
}
