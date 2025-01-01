// cron_manager.go
//
// Description:
// A CLI tool for managing cron jobs in Linux systems. This tool allows users to:
// - List current cron jobs in a formatted table.
// - Add new cron jobs with example formats and syntax guidance.
// - Remove existing cron jobs using an intuitive interface.
// - Display detailed cron syntax and command examples.
//
// Features:
// - Uses "go-pretty" for colorful table rendering.
// - Handles invalid inputs gracefully and provides user-friendly feedback.
// - Designed for simplicity and ease of use in managing cron tasks.

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// List all cron jobs
func listCronJobs() {
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.Output()

	// Create the table writer
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Current Cron Jobs")
	t.AppendHeader(table.Row{"#", "Cron Job"})

	// Handle errors or empty output
	if err != nil {
		if strings.Contains(err.Error(), "no crontab for") { // No cron jobs set
			t.AppendRow(table.Row{"-", "No current cron jobs"})
			t.Render()
			return
		} else {
			fmt.Println("Error fetching cron jobs:", err)
			return
		}
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		t.AppendRow(table.Row{"-", "No current cron jobs"})
	} else {
		cronJobs := strings.Split(outputStr, "\n")
		count := 0
		for _, job := range cronJobs {
			trimmedJob := strings.TrimSpace(job)
			if trimmedJob != "" { // Skip empty lines
				count++
				t.AppendRow(table.Row{count, trimmedJob})
			}
		}
		if count == 0 {
			t.AppendRow(table.Row{"-", "No current cron jobs"})
		}
	}

	// Render the table
	t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	println("\n") //spacing
	t.Render()
}


func addCronJob() {
	for {
		listCronJobs()

		fmt.Println("\n") // spacing
		printExampleCommands()
		fmt.Println("\nEnter 'help' for detailed cron syntax information.")
		fmt.Println("Enter 'back' to return to the main menu.")
		fmt.Println("\nEnter a new cron job:")

		fmt.Print("\n> ")
		reader := bufio.NewReader(os.Stdin)
		cronJob, _ := reader.ReadString('\n')
		cronJob = strings.TrimSpace(cronJob)

		// Handle special inputs
		switch strings.ToLower(cronJob) {
		case "help":
			printCronSyntaxTable()
			continue // Stay in the addCronJob function after showing help
		case "back":
			return // Exit to the main menu
		}

		// Fetch existing cron jobs
		cmd := exec.Command("crontab", "-l")
		output, err := cmd.Output()
		var currentCronJobs string
		if err == nil { // Handle the case where there are no existing cron jobs
			currentCronJobs = string(output)
		}

		// Add the new cron job to the list
		updatedCronJobs := currentCronJobs + "\n" + cronJob

		// Write the updated cron jobs back
		cmd = exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | crontab -", updatedCronJobs))
		if err := cmd.Run(); err != nil {
			fmt.Println("Error updating cron jobs:", err)
		} else {
			fmt.Println("Cron job added successfully!")
		}
		return // Return to the main menu after adding the job
	}
}


func printExampleCommands() {

	e := table.NewWriter()
	e.SetOutputMirror(os.Stdout)
	e.AppendHeader(table.Row{"Cron Job Example Formats", "Description"})
	e.AppendRow(table.Row{"0 * * * * echo 'Hello'", "Every hour"})
	e.AppendRow(table.Row{"0 0 * * * /path/to/backup.sh", "Every day"})
	e.AppendRow(table.Row{"0 0 1 * * pacman -Sc", "Every month"})
	e.AppendRow(table.Row{"0 0 1 1 * sudo pacman -Syu", "Every year"})
	e.AppendRow(table.Row{"*/5 * * * * ping -c 4 google.com", "Every 5 minutes"})
	e.AppendRow(table.Row{"0 2 * * 0 sudo paccache -r", "Every Sunday at 2am"})
	e.AppendRow(table.Row{"0 10 1 * * df -Th | grep -v fs", "Every 1st of the month at 10am"})
	e.AppendRow(table.Row{"*/2 * * * * command >> /path/to/cron_output.log 2>&1","Log into a file every 2 minutes"})
	//	2>&1: Redirects error output (stderr) to the same file as standard output.

	e.SetStyle(table.StyleColoredBright)
	e.Render()

}

func printEntryFormat() {
	f := table.NewWriter()
	f.SetOutputMirror(os.Stdout)
	// rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	f.SetTitle("Enter the cron job in the format")
	f.AppendRow(table.Row{"min", "hour", "day", "month", "weekday", "command"})
	f.AppendRow(table.Row{"*", "*", "*", "*", "*", "command"})
	f.SetStyle(table.StyleColoredBlackOnGreenWhite)
	f.Render()
}

// Function to print the cron syntax table
func printCronSyntaxTable() {
	// rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	fmt.Println("\n") // spacing
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Cron Job Syntax Table")
	t.AppendHeader(table.Row{"Field", "Allowed Values", "Description"})
	t.AppendRow(table.Row{"min", "0-59", "Minute field"})
	t.AppendRow(table.Row{"hour", "0-23", "Hour field"})
	t.AppendRow(table.Row{"day", "1-31", "Day of the month"})
	t.AppendRow(table.Row{"month", "1-12 or Jan-Dec", "Month field"})
	t.AppendRow(table.Row{"weekday", "0-6 or Sun-Sat", "Day of the week (0 = Sunday)"})
	t.AppendRow(table.Row{"command", "any valid command or script", "The command/script to be executed"})
	
	
	// Set style for table formatting
	t.SetStyle(table.StyleColoredBright)
	t.Render()
	println("\n")

	printEntryFormat()
}

// Remove a cron job
func removeCronJob() {
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.Output()
	println("\n")

	// Create the table writer
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnRedWhite)

	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Remove Cron Jobs")
	t.AppendHeader(table.Row{"#", "Cron Job"})

	// Handle errors or empty output
	if err != nil {
		if strings.Contains(err.Error(), "no crontab for") { // No cron jobs set
			t.AppendRow(table.Row{"-", "No current cron jobs"})
			t.Render()
			return
		} else {
			fmt.Println("Error fetching cron jobs:", err)
			return
		}
	}

	// Process cron jobs and filter empty lines
	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		t.AppendRow(table.Row{"-", "No current cron jobs"})
		t.Render()
		return
	}

	cronJobs := strings.Split(outputStr, "\n")
	validJobs := []string{}
	for _, job := range cronJobs {
		trimmedJob := strings.TrimSpace(job)
		if trimmedJob != "" {
			validJobs = append(validJobs, trimmedJob)
		}
	}

	if len(validJobs) == 0 {
		t.AppendRow(table.Row{"-", "No current cron jobs"})
		t.Render()
		return
	}

	// Display the cron jobs in a table
	for i, job := range validJobs {
		t.AppendRow(table.Row{i + 1, job})
	}
	// t.SetStyle(table.StyleColoredBlackOnRedWhite)
	t.Render()

	// Ask user for the job to remove
	fmt.Println("\nEnter the number (#) of the job to remove.")
	fmt.Println("Enter 'back' to return to the main menu.")

	fmt.Print("\n> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Handle special input "back"
	if strings.ToLower(input) == "back" {
		return // Exit to the main menu
	}

	// Parse the input as an integer
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(validJobs) {
		fmt.Println("Invalid choice!")
		return
	}

	// Remove the selected job
	validJobs = append(validJobs[:choice-1], validJobs[choice:]...)

	// Write the updated cron jobs back
	updatedCronJobs := strings.Join(validJobs, "\n")
	cmd = exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | crontab -", updatedCronJobs))
	if err := cmd.Run(); err != nil {
		fmt.Println("Error updating cron jobs:", err)
	} else {
		fmt.Println("Cron job removed successfully!")
	}
}


func main() {
	for {
		fmt.Println("\nCron Job Manager")
		fmt.Println("\n1. List Cron Jobs")
		fmt.Println("2. Add Cron Job")
		fmt.Println("3. Remove Cron Job")
		fmt.Println("4. Show syntax help")
		fmt.Println("5. Exit")
		fmt.Print("Choose an option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			listCronJobs()
		case 2:
			addCronJob()
		case 3:
			removeCronJob()
		case 4:
			printCronSyntaxTable()
			println("\n")
			printExampleCommands()
		case 5:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
