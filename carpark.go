package main

import (
	"container/list"
	"errors"
	"sortedlist"
)

//Carpark represents the carpark map, empty slots, and maximum number of slots filled
type Carpark struct {
	Map         map[int]Vehicle //Properties of each vehicle parked in the carpark
	highestSlot int             //Highest number of slots filled throughout carpark operation
	maxSlot     int             //Maximum number of slots available
	emptySlots  *list.List      //List containing sorted empty slots in ascending order
}

//Initialize carpark parameters
func (carpark *Carpark) init(maxSlot int) error {
	if err := carpark.initStatus(); err == nil {
		return errors.New("Carpark already initialized")
	}
	carpark.Map = make(map[int]Vehicle) //Setup a map of the carpark
	carpark.emptySlots = list.New()     //Setup an empty heap of empty parking slots
	carpark.maxSlot = maxSlot           //Set the maximum number of slots
	return nil
}

//Park a vehicle in carpark
func (carpark *Carpark) insertCar(vehicle Vehicle) (int, error) {
	if err := carpark.initStatus(); err != nil {
		return 0, err
	}
	if vehicle == nil {
		return 0, errors.New("Unknown or nil vehicle")
	}

	var slotNo int
	slotsNeeded := vehicle.getSlotsNeeded()
	var emptySlot = sortedlist.FindSeq(carpark.emptySlots, slotsNeeded) //Get nearest empty slot which was previously occupied
	if emptySlot != nil {                                               //Park vehicle at nearest empty slot
		sortedlist.Remove(carpark.emptySlots, emptySlot, slotsNeeded)
		slotNo = emptySlot.Value.(int)
	} else { //Park vehicle at next available highest slot
		if carpark.highestSlot+slotsNeeded > carpark.maxSlot {
			return 0, errors.New("Sorry, parking lot is full")
		}
		slotNo = carpark.highestSlot + 1
		carpark.highestSlot += slotsNeeded
	}

	//Insert the vehicle into the map
	*vehicle.getSlot() = slotNo
	carpark.Map[slotNo] = vehicle
	return slotNo, nil
}

//Remove vehicle from carpark
func (carpark *Carpark) removeCar(slotNo int) error {
	if err := carpark.initStatus(); err != nil {
		return err
	}
	if vehicle, ok := carpark.Map[slotNo]; ok {
		//Remove vehicle from carpark Map
		delete(carpark.Map, slotNo)
		//Add empty slot to the heap
		sortedlist.Insert(carpark.emptySlots, carpark.emptySlots.Back(), slotNo, vehicle.getSlotsNeeded())
		return nil
	}
	return errors.New("Vehicle non-existent in carpark")
}

//Given a vehicle colour, retrieve the vehicle slot and registration numbers
func (carpark *Carpark) getCarsWithColour(colour string) ([]int, []string, error) {
	var slots []int
	var registrations []string
	for i := 1; i <= carpark.highestSlot; i++ {
		vehicle, ok := carpark.Map[i]
		if ok && *vehicle.getColour() == colour {
			slots = append(slots, *vehicle.getSlot())
			registrations = append(registrations, *vehicle.getRegistration())
		}
	}
	if slots == nil {
		return nil, nil, errors.New("Not found")
	}
	return slots, registrations, nil
}

//Given a vehicle registration number, retrieve the vehicle slot number
func (carpark *Carpark) getCarWithRegistrationNo(registration string) (int, error) {
	for _, vehicle := range carpark.Map {
		if *vehicle.getRegistration() == registration {
			return *vehicle.getSlot(), nil
		}
	}
	return 0, errors.New("Not found")
}

//Retrieve ordered sequence of vehicles parked in the carpark
func (carpark *Carpark) getStatus() []Vehicle {
	var vehicles []Vehicle
	for i := 1; i <= carpark.highestSlot; i++ {
		vehicle, ok := carpark.Map[i]
		if ok {
			vehicles = append(vehicles, vehicle)
		}
	}
	return vehicles
}

//Check whether the carpark has been initialized
func (carpark *Carpark) initStatus() error {
	if carpark.Map == nil {
		return errors.New("Carpark not initialized")
	}
	return nil
}
