package main

import "sync"

type Seat struct {
	number   int
	occupied bool
	mutex    sync.Mutex
}

func (s *Seat) Number() int { return s.number }
func (s *Seat) IsOccupied() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.occupied
}
func (s *Seat) AttemptToOccupy() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.occupied {
		s.occupied = true
		return true
	}
	return false
}
func (s *Seat) Vacate() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.occupied = false
}

type SeatManager struct {
	seats []*Seat
}

var seatManager = &SeatManager{
	seats: make([]*Seat, len(PHILOSOPHER_NAMES)),
}

func init() {
	for i := range seatManager.seats {
		seatManager.seats[i] = &Seat{number: i}
	}
}

func (sm *SeatManager) Seats() []*Seat {
	result := make([]*Seat, len(sm.seats))
	copy(result, sm.seats)
	return result
}

func (sm *SeatManager) AvailableSeat() *Seat {
	for _, seat := range sm.seats {
		if !seat.IsOccupied() {
			return seat
		}
	}
	return nil
} 