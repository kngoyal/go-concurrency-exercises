//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"fmt"
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID                 int
	IsPremium          bool
	TimeUsed           int64 // in seconds
	ConcurrentRequests int
	sync.Mutex
}

func HandleRequestByRequestLimit(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}
	c := make(chan bool)

	go func() {
		process()
		c <- true
	}()

	select {
	case <-c:
		return true
	case <-time.After(10 * time.Second):
		return false
	}
}

func (u *User) GetTimeUsed() int64 {
	u.Lock()
	defer u.Unlock()
	return u.TimeUsed
}

func (u *User) SetTimeUsed(t int64) {
	u.Lock()
	defer u.Unlock()
	u.TimeUsed += t
}

func (u *User) UpdateConcurrentRequests(t int) {
	u.Lock()
	defer u.Unlock()
	u.ConcurrentRequests += t
}

func (u *User) GetConcurrentRequests() int {
	u.Lock()
	defer u.Unlock()
	return u.ConcurrentRequests
	//return int(math.Max(1.0, float64(u.ConcurrentRequests)))
}

func HandleRequestByUserLimit(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}
	c := make(chan time.Time)
	start := time.Now()
	go func() {
		u.UpdateConcurrentRequests(1)
		process()
		u.UpdateConcurrentRequests(-1)
		c <- time.Now()
	}()

	select {
	case t1 := <-c:
		fmt.Println("completed", start, t1)
		u.SetTimeUsed(int64(t1.Sub(start).Seconds()))
		return true
	case t2 := <-time.After(time.Duration(10-u.GetTimeUsed()) / time.Duration(1+u.GetConcurrentRequests()) * time.Second):
		fmt.Println("killed", start, t2)
		u.SetTimeUsed(int64(t2.Sub(start).Seconds()))
		return false
	}
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	return HandleRequestByUserLimit(process, u)
	//return HandleRequestByRequestLimit(process, u)
}

func main() {
	RunMockServer()
}
