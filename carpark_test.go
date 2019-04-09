package main

import (
	"container/list"
	"reflect"
	"testing"
)

//variables act as a struct of all parameters used in testing
type variables struct {
	vehicle0     *Motorcycle
	vehicle1     *Motorcycle
	vehicle2     *Motorcycle
	map0         map[int]Vehicle
	map1         map[int]Vehicle
	map2         map[int]Vehicle
	mapAll       map[int]Vehicle
	item1        interface{}
	item2        interface{}
	emptySlot0   *list.List
	emptySlot1   *list.List
	emptySlot2   *list.List
	emptySlotAll *list.List
}

//values() acts a storage of default values and return a 'variables' struct containing default values
func values() variables {
	defaultValues := variables{
		vehicle0:   &Motorcycle{baseVehicle: baseVehicle{registration: "KA-01-HH-2701", colour: "Blue"}},
		vehicle1:   &Motorcycle{baseVehicle: baseVehicle{registration: "KA-01-HH-1234", colour: "White", slot: 1}},
		vehicle2:   &Motorcycle{baseVehicle: baseVehicle{registration: "KA-01-HH-7777", colour: "Red", slot: 2}},
		map0:       make(map[int]Vehicle),
		item1:      1,
		item2:      2,
		emptySlot0: list.New(),
	}
	defaultValues.map1 = map[int]Vehicle{1: defaultValues.vehicle1}
	defaultValues.map2 = map[int]Vehicle{2: defaultValues.vehicle2}
	defaultValues.mapAll = map[int]Vehicle{1: defaultValues.vehicle1, 2: defaultValues.vehicle2}
	defaultValues.emptySlot1 = list.New()
	defaultValues.emptySlot1.PushBack(defaultValues.item1)
	defaultValues.emptySlot2 = list.New()
	defaultValues.emptySlot2.PushBack(defaultValues.item2)
	defaultValues.emptySlotAll = list.New()
	defaultValues.emptySlotAll.PushBack(defaultValues.item1)
	defaultValues.emptySlotAll.PushBack(defaultValues.item2)

	return defaultValues
}

//Compare two 'Carpark' structs
func compareCarpark(t *testing.T, carpark *Carpark, wantCarpark *Carpark) {
	if !reflect.DeepEqual(carpark.Map, wantCarpark.Map) ||
		!reflect.DeepEqual(carpark.emptySlots, wantCarpark.emptySlots) ||
		carpark.highestSlot != wantCarpark.highestSlot ||
		carpark.maxSlot != wantCarpark.maxSlot {
		t.Errorf("gotCarpark = %v, wantCarpark = %v", carpark, wantCarpark)
	}
}

func TestCarpark_init(t *testing.T) {
	type args struct {
		maxSlot int
	}
	tests := []struct {
		name        string
		carpark     *Carpark
		args        args
		wantErr     bool
		wantCarpark *Carpark
	}{
		{name: "Carpark not initialized",
			carpark:     &Carpark{},
			args:        args{maxSlot: 12},
			wantErr:     false,
			wantCarpark: &Carpark{Map: values().map0, emptySlots: values().emptySlot0, highestSlot: 0, maxSlot: 12},
		},
		{name: "Carpark already initialized",
			carpark:     &Carpark{Map: values().map0, emptySlots: values().emptySlot0, highestSlot: 8, maxSlot: 10},
			args:        args{maxSlot: 12},
			wantErr:     true,
			wantCarpark: &Carpark{Map: values().map0, emptySlots: values().emptySlot0, highestSlot: 8, maxSlot: 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.carpark.init(tt.args.maxSlot)
			if (err != nil) != tt.wantErr {
				t.Errorf("Carpark.init() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			compareCarpark(t, tt.carpark, tt.wantCarpark)
		})
	}
}

func TestCarpark_insertCar(t *testing.T) {
	type args struct {
		car *Motorcycle
	}
	tests := []struct {
		name        string
		carpark     *Carpark
		args        args
		want        int
		wantErr     bool
		wantCarpark *Carpark
	}{
		{name: "Carpark not initialized",
			carpark:     &Carpark{},
			args:        args{car: values().vehicle1},
			want:        0,
			wantErr:     true,
			wantCarpark: &Carpark{},
		},
		{name: "Insert car into new slot",
			carpark:     &Carpark{Map: values().map1, emptySlots: values().emptySlot0, highestSlot: 1, maxSlot: 10},
			args:        args{car: values().vehicle2},
			want:        2,
			wantErr:     false,
			wantCarpark: &Carpark{Map: values().mapAll, emptySlots: values().emptySlot0, highestSlot: 2, maxSlot: 10},
		},
		{name: "Insert car into a previously occupied but now free slot",
			carpark:     &Carpark{Map: values().map2, emptySlots: values().emptySlot1, highestSlot: 2, maxSlot: 10},
			args:        args{car: values().vehicle1},
			want:        1,
			wantErr:     false,
			wantCarpark: &Carpark{Map: values().mapAll, emptySlots: values().emptySlot0, highestSlot: 2, maxSlot: 10},
		},
		{name: "Insert car beyond maxSlot",
			carpark:     &Carpark{Map: values().mapAll, emptySlots: values().emptySlot0, highestSlot: 2, maxSlot: 2},
			args:        args{car: values().vehicle0},
			want:        0,
			wantErr:     true,
			wantCarpark: &Carpark{Map: values().mapAll, emptySlots: values().emptySlot0, highestSlot: 2, maxSlot: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.carpark.insertCar(tt.args.car)
			if (err != nil) != tt.wantErr {
				t.Errorf("Carpark.insertCar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Carpark.insertCar() = %v, want %v", got, tt.want)
			}
			compareCarpark(t, tt.carpark, tt.wantCarpark)
		})
	}
}

func TestCarpark_removeCar(t *testing.T) {
	type args struct {
		slotNo int
	}
	tests := []struct {
		name        string
		carpark     *Carpark
		args        args
		wantErr     bool
		wantCarpark *Carpark
	}{
		{name: "Carpark not initialized",
			carpark:     &Carpark{},
			args:        args{slotNo: 1},
			wantErr:     true,
			wantCarpark: &Carpark{},
		},
		{name: "Remove car",
			carpark:     &Carpark{Map: values().mapAll, emptySlots: values().emptySlot0, highestSlot: 2, maxSlot: 10},
			args:        args{slotNo: 1},
			wantErr:     false,
			wantCarpark: &Carpark{Map: values().map2, emptySlots: values().emptySlot1, highestSlot: 2, maxSlot: 10},
		},
		{name: "Remove non-existent car",
			carpark:     &Carpark{Map: values().map1, emptySlots: values().emptySlot2, highestSlot: 2, maxSlot: 10},
			args:        args{slotNo: 2},
			wantErr:     true,
			wantCarpark: &Carpark{Map: values().map1, emptySlots: values().emptySlot2, highestSlot: 2, maxSlot: 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.carpark.removeCar(tt.args.slotNo); (err != nil) != tt.wantErr {
				t.Errorf("Carpark.removeCar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			compareCarpark(t, tt.carpark, tt.wantCarpark)
		})
	}
}

func TestCarpark_getCarsWithColour(t *testing.T) {
	type args struct {
		colour string
	}
	tests := []struct {
		name    string
		carpark *Carpark
		args    args
		want    []int
		want1   []string
		wantErr bool
	}{
		{name: "Carpark with car of requested colour",
			carpark: &Carpark{Map: values().map1, emptySlots: values().emptySlot0, highestSlot: 1, maxSlot: 10},
			args:    args{colour: "White"},
			want:    []int{1},
			want1:   []string{"KA-01-HH-1234"},
			wantErr: false,
		},
		{name: "Carpark without car of requested colour",
			carpark: &Carpark{Map: values().map2, emptySlots: values().emptySlot1, highestSlot: 2, maxSlot: 10},
			args:    args{colour: "White"},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
		{name: "Empty carpark",
			carpark: &Carpark{Map: values().map0, emptySlots: values().emptySlot0, highestSlot: 0, maxSlot: 10},
			args:    args{colour: "White"},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
		{name: "Uninitialized carpark",
			carpark: &Carpark{},
			args:    args{colour: "White"},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.carpark.getCarsWithColour(tt.args.colour)
			if (err != nil) != tt.wantErr {
				t.Errorf("Carpark.getCarsWithColour() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Carpark.getCarsWithColour() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Carpark.getCarsWithColour() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCarpark_getCarWithRegistrationNo(t *testing.T) {
	type args struct {
		registration string
	}
	tests := []struct {
		name    string
		carpark *Carpark
		args    args
		want    int
		wantErr bool
	}{
		{name: "Carpark with car of requested colour",
			carpark: &Carpark{Map: values().map1, emptySlots: values().emptySlot0, highestSlot: 1, maxSlot: 10},
			args:    args{registration: "KA-01-HH-1234"},
			want:    1,
			wantErr: false,
		},
		{name: "Carpark without car of requested colour",
			carpark: &Carpark{Map: values().map2, emptySlots: values().emptySlot1, highestSlot: 2, maxSlot: 10},
			args:    args{registration: "KA-01-HH-1234"},
			want:    0,
			wantErr: true,
		},
		{name: "Empty carpark",
			carpark: &Carpark{Map: values().map0, emptySlots: values().emptySlot0, highestSlot: 0, maxSlot: 10},
			args:    args{registration: "KA-01-HH-1234"},
			want:    0,
			wantErr: true,
		},
		{name: "Uninitialized carpark",
			carpark: &Carpark{},
			args:    args{registration: "KA-01-HH-1234"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.carpark.getCarWithRegistrationNo(tt.args.registration)
			if (err != nil) != tt.wantErr {
				t.Errorf("Carpark.getCarWithRegistrationNo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Carpark.getCarWithRegistrationNo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCarpark_getStatus(t *testing.T) {
	tests := []struct {
		name    string
		carpark *Carpark
		want    []Vehicle
	}{
		{name: "Carpark not initialized",
			carpark: &Carpark{},
			want:    nil,
		},
		{name: "Empty carpark",
			carpark: &Carpark{Map: values().map0, emptySlots: values().emptySlot0, highestSlot: 0, maxSlot: 10},
			want:    nil,
		},
		{name: "Carpark with cars",
			carpark: &Carpark{Map: values().mapAll, emptySlots: values().emptySlot0, highestSlot: 2, maxSlot: 10},
			want:    []Vehicle{values().vehicle1, values().vehicle2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.carpark.getStatus(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Carpark.getStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
