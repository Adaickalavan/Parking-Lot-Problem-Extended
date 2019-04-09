package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"pretty"
	"runtime"
	"strconv"
	"strings"
	"text/tabwriter"
)

var inputInteractive io.Reader = os.Stdin
var outStream io.Writer = os.Stdout

func main() {

	//Input file or interactive mode
	ii := len(os.Args)
	var scanner *bufio.Scanner
	switch {
	case ii > 2:
		log.Fatal("Unknown command line input")
	case ii == 2:
		inputFile, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer inputFile.Close()
		scanner = bufio.NewScanner(inputFile)
	default:
		scanner = bufio.NewScanner(inputInteractive)
	}

	//Create a carpark
	var carpark = &Carpark{}

	//Operate the carpark
	operateCarpark(carpark, scanner)
}

//operateCarpark reads input queries from console or text file and executes the command
func operateCarpark(carpark *Carpark, scanner *bufio.Scanner) {
	newlineStr := getNewlineStr()
	exit := false
	for !exit && scanner.Scan() {
		input := scanner.Text()
		input = strings.TrimRight(input, newlineStr)
		s := parse(input)

		switch {
		case s[0] == "create_parking_lot" && len(s) == 2: //Initialize carpark
			maxSlot, err := strconv.Atoi(s[1])
			if checkError(err) {
				break
			}
			err = carpark.init(maxSlot)
			if !checkError(err) {
				fmt.Fprintf(outStream, "Created a parking lot with %v slots\n", maxSlot)
			}

		case s[0] == "park" && len(s) == 4: //Park a new vehicle
			var vehicle Vehicle
			switch s[3] {
			case "car":
				car := NewCar()
				car.registration = s[1]
				car.colour = s[2]
				vehicle = car
			case "motorcycle":
				motorcycle := NewMotorcycle()
				motorcycle.registration = s[1]
				motorcycle.colour = s[2]
				vehicle = motorcycle
			case "bus":
				bus := NewBus()
				bus.registration = s[1]
				bus.colour = s[2]
				vehicle = bus
			}
			slotNo, err := carpark.insertCar(vehicle)
			if !checkError(err) {
				fmt.Fprintf(outStream, "Allocated slot number: %v\n", slotNo)
			}

		case s[0] == "leave" && len(s) == 2: //Remove a parked vehicle
			slotNo, err := strconv.Atoi(s[1])
			if checkError(err) {
				break
			}
			err = carpark.removeCar(slotNo)
			if !checkError(err) {
				fmt.Fprintf(outStream, "Slot number %v is free\n", slotNo)
			}

		case s[0] == "registration_numbers_for_cars_with_colour" && len(s) == 2: //Return registration numbers with given vehicle colour
			_, registration, err := carpark.getCarsWithColour(s[1])
			if checkError(err) {
				break
			}
			err = pretty.Printer(registration, outStream)
			if err != nil {
				panic(err.Error())
			}

		case s[0] == "slot_numbers_for_cars_with_colour" && len(s) == 2: //Return slot numbers with given vehicle colour
			slots, _, err := carpark.getCarsWithColour(s[1])
			if checkError(err) {
				break
			}
			err = pretty.Printer(slots, outStream)
			if err != nil {
				panic(err.Error())
			}

		case s[0] == "slot_number_for_registration_number" && len(s) == 2: //Return slot numbers with given vehicle registration number
			slotNo, err := carpark.getCarWithRegistrationNo(s[1])
			if !checkError(err) {
				fmt.Fprintln(outStream, slotNo)
			}

		case s[0] == "status" && len(s) == 1: //Retrieve vehicles parked in carpark
			vehicles := carpark.getStatus()
			var w = tabwriter.NewWriter(outStream, 0, 0, 4, ' ', 0)
			fmt.Fprintln(w, "Slot No.\tRegistration No\tColour\tType")
			for _, vehicle := range vehicles {
				s := fmt.Sprintf("%v\t%s\t%s\t%s", *vehicle.getSlot(), *vehicle.getRegistration(), *vehicle.getColour(), vehicle.getType())
				fmt.Fprintln(w, s)
			}
			w.Flush()

		case s[0] == "exit" && len(s) == 1: //End carpark operation
			exit = true

		default: //Default option
			fmt.Fprintln(outStream, "Unknown input command")
		}
	}
}

//getNewlineStr identifies operating system and returns newline character used
func getNewlineStr() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func parse(input string) []string {
	s := strings.Split(input, " ")
	return s
}

func checkError(err error) bool {
	if err != nil {
		fmt.Fprintln(outStream, err.Error())
		return true
	}
	return false
}
