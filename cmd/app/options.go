package main

import (
	"errors"
	"flag"
	"fmt"
	"runtime"
)

type Options struct {
	AutoDetect *bool
	BaudRate *int
	Daemon *bool
	Help *bool
	PlotInterval *int
	PrintGPSCoordsToCLI *bool
	PrintNMEAToCLI *bool
	SerialPort *int
	Silent *bool
	Timeout *int
	Verbose *bool
	WriteCSVFilePath *string
	WriteGPSCoordsFilePath *string
	WriteKMLFilePath *string
	WriteNMEAFilePath *string
}

func parseOptions() *Options {
	o := Options{}

	o.AutoDetect = flag.Bool("autodetect", true, "Auto detect the serial port and baud rate for the connected GPS device. Partially or fully disabled if baud rate and/or port is manually set.")
	o.BaudRate = flag.Int("baudrate", -1, "Set the baud rate for the serial port.")
	o.Daemon = flag.Bool("daemon", false, "Run as a background task.")
	o.Help = flag.Bool ("help", false, "Print help sheet.")
	o.PlotInterval = flag.Int("interval", 30, "Set the plot interval (seconds) for returning a GPS location from device.")
	o.PrintGPSCoordsToCLI = flag.Bool("print-gps", false, "Print the GPS coordinates to standard out.")
	o.PrintNMEAToCLI = flag.Bool("print-nmea", false, "Print NMEA messages to standard out.")
	o.SerialPort = flag.Int("port", -1, "Set the serial port to connect.")
	o.Silent = flag.Bool("silent", false, "No output will be sent to standard out. Cannot be used with flags that write to standard out.")
	o.Timeout = flag.Int("timeout", 60, "Set the timeout (seconds) before disconnecting on error or inactivity.")
	o.Verbose = flag.Bool("verbose", false, "Extra information provided in standard out.")
	o.WriteCSVFilePath = flag.String("write-csv", "", "Write timestamp, GPS coordinates, and NMEA message(s) for location to CSV file at path provided.")
	o.WriteGPSCoordsFilePath = flag.String("write-gps", "", "Write raw GPS coordinates to file at path provided.")
	o.WriteKMLFilePath = flag.String("write-kml", "", "Write Google Maps / Earth KML format as a waypoint workflow to file at path provided.")
	o.WriteNMEAFilePath = flag.String("write-nmea", "", "Write raw NMEA messages to file at path provided.")

	flag.Parse()

	if *o.AutoDetect && (*o.BaudRate > -1 && *o.SerialPort > -1) {
		*o.AutoDetect = false
	}

	return &o
}

func printHelpSheet() {
	fmt.Println("GPS Atlas / gps-usb-serial-reader")
	fmt.Println("Auto-detect, plot, and map with common GPS USB serial devices")
	fmt.Print("\nARGUMENTS:\n\n")
	flag.PrintDefaults()
}

func checkOptionSanity(o *Options) error {
	if runtime.GOOS == "windows" && !*o.AutoDetect && (*o.SerialPort <= 0 || *o.SerialPort > 256) {
		return errors.New("COM serial ports should be between 0-255")
	}

	if !*o.AutoDetect && *o.SerialPort < 0 {
		return errors.New("serial port cannot be negative")
	}

	if !*o.AutoDetect && *o.BaudRate <= 0 {
		return errors.New("baud rate cannot be less than 0")
	}

	if *o.Silent && *o.Verbose {
		return errors.New("Silent and Verbose flags cannot both be set")
	}

	if *o.Silent && (*o.PrintNMEAToCLI || *o.PrintGPSCoordsToCLI) {
		return errors.New("can't be silent and paired with a flag that increases standard out verbosity")
	}

	if *o.Timeout < 0 {
		return errors.New("timeout cannot be negative")
	}

	if *o.PlotInterval <= 0 {
		return errors.New("plot interval cannot be less than 0")
	}

	return nil
}