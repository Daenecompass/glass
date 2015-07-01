package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/timeglass/glass/_vendor/github.com/hashicorp/errwrap"
)

type Keeper struct {
	ledgerPath string
	stop       chan struct{}

	keeperData *keeperData
}

type keeperData struct {
	TickRate time.Duration     `json:"tick_rate"`
	Timers   map[string]*Timer `json:"timers"`
}

func NewKeeper(path string) (*Keeper, error) {
	k := &Keeper{
		stop: make(chan struct{}),
		keeperData: &keeperData{
			Timers:   map[string]*Timer{},
			TickRate: time.Minute,
		},
	}

	//attempt to open json file, if it exsts
	k.ledgerPath = filepath.Join(path, "ledger.json")
	return k, k.Load()
}

func (k *Keeper) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &k.keeperData)
}

func (k *Keeper) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.keeperData)
}

func (k *Keeper) Add(t *Timer) error {
	if tt, ok := k.keeperData.Timers[t.Dir()]; !ok {
		k.keeperData.Timers[t.Dir()] = t
		return t.Start()
	} else {
		tt.Unpause()
		return nil
	}
}

func (k *Keeper) Get(dir string) (*Timer, error) {
	if t, ok := k.keeperData.Timers[dir]; ok {
		return t, nil
	}

	return nil, fmt.Errorf("No known timer for '%s'", dir)
}

func (k *Keeper) Remove(dir string) error {
	if t, ok := k.keeperData.Timers[dir]; ok {
		delete(k.keeperData.Timers, dir)
		return t.Stop()
	}

	return fmt.Errorf("No known timer for '%s'", dir)
}

func (k *Keeper) Stop() {
	k.stop <- struct{}{}
}

func (k *Keeper) Start() {
	log.Printf("Started time keeper on %s", time.Now())
	defer func() {
		log.Printf("Stopped time keeper on %s", time.Now())
	}()

	//@todo, instead of a tickrate, instead only auto-save
	//when certain data was changed
	for {

		//save state
		err := k.Save()
		if err != nil {
			log.Printf("Error while saving to ledger: %s", err)
		}

		//stop or wait for next tick
		select {
		case <-k.stop:
			return
		case <-time.After(k.TickRate()):
		}
	}
}

func (k *Keeper) TickRate() time.Duration {
	return k.keeperData.TickRate
}

func (k *Keeper) Load() error {
	f, err := os.Open(k.ledgerPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return errwrap.Wrapf(fmt.Sprintf("Failed to open '%s': {{err}}", k.ledgerPath), err)
		}
	} else {
		defer f.Close()
		dec := json.NewDecoder(f)
		err := dec.Decode(k)
		if err != nil {
			return errwrap.Wrapf(fmt.Sprintf("Failed to decode JSON in '%s': {{err}}", k.ledgerPath), err)
		}

		//immediately restart if not paused
		for _, t := range k.keeperData.Timers {
			if !t.IsPaused() {
				err := t.Start()
				if err != nil {
					//failure to start a single timer shouldn't result in failure to start the daemon
					//for example, maybe the project was moved or symlinks couldn't be resolved
					log.Printf("Failed to start timer for '%s' after loaded from ledger: %s", t.Dir(), err)
					continue
				}
			}
		}

	}

	return nil
}

func (k *Keeper) Save() error {
	f, err := os.OpenFile(k.ledgerPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(k)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error saving ledger to '%s': {{err}}", k.ledgerPath), err)
	}

	return nil
}
