package main

//Vehicle represents car, motorcycle, and bus
type Vehicle interface {
	getRegistration() *string
	getColour() *string
	getSlot() *int
	getSlotsNeeded() int
	getType() string
}

type baseVehicle struct {
	name         string //Type of vehicle
	registration string //Registration number of car
	colour       string //Colour of car
	slot         int    //Slot number in which the motorcycle is parked
}

func (basevehicle *baseVehicle) fit() bool {
	return true
}

func (basevehicle *baseVehicle) getRegistration() *string {
	return &basevehicle.registration
}

func (basevehicle *baseVehicle) getColour() *string {
	return &basevehicle.colour
}

func (basevehicle *baseVehicle) getSlot() *int {
	return &basevehicle.slot
}

func (basevehicle *baseVehicle) getType() string {
	return basevehicle.name
}

// Car represents the properties of a car
type Car struct {
	baseVehicle
}

func (car *Car) getSlotsNeeded() int {
	return 2
}

//NewCar is a car constructor function
func NewCar() *Car {
	return &Car{baseVehicle: baseVehicle{name: "Car"}}
}

//Motorcycle represents the properties of a motorcycle
type Motorcycle struct {
	baseVehicle
}

func (motorcycle *Motorcycle) getSlotsNeeded() int {
	return 1
}

//NewMotorcycle is a motorcycle constructor function
func NewMotorcycle() *Motorcycle {
	return &Motorcycle{baseVehicle: baseVehicle{name: "Motorcycle"}}
}

//Bus represents the properties of a bus
type Bus struct {
	baseVehicle
}

func (bus *Bus) getSlotsNeeded() int {
	return 3
}

//NewBus is a bus constructor function
func NewBus() *Bus {
	return &Bus{baseVehicle: baseVehicle{name: "Bus"}}
}
